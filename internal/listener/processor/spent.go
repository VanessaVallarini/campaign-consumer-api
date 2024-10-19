package processor

import (
	"context"

	addressv1 "github.com/VanessaVallarini/address-api/pkg/api/proto/v1"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/geocode"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type SpentDao interface {
	FetchByMerchantIdAndBucket(context.Context, uuid.UUID, string) (model.Spent, error)
	Insert(context.Context, transaction.Transaction, model.Spent) error
	Upsert(context.Context, transaction.Transaction, model.Spent) error
}

type CampaignManager interface {
	Fetch(context.Context, uuid.UUID) (model.Campaign, error)
	Upsert(context.Context, model.Campaign) error
}

type MerchantRetriever interface {
	Fetch(context.Context, uuid.UUID) (model.Merchant, error)
}

type SlugRetriever interface {
	Fetch(context.Context, string) (model.Slug, error)
}

type RegionRetriever interface {
	Fetch(context.Context, uuid.UUID) (model.Region, error)
}

type LedgerManager interface {
	Create(context.Context, transaction.Transaction, model.Ledger) error
}

type BucketFetcher interface {
	CurrentBucket() model.Bucket
}

type SpentManager interface {
	FetchByMerchantIdAndBucket(context.Context, uuid.UUID, string) (model.Spent, error)
	UpsertAndRegsterLedger(context.Context, model.Spent, model.SpentEvent, string) error
}

type SpentProcessor struct {
	spentManager      SpentManager
	campaignManager   CampaignManager
	merchantRetriever MerchantRetriever
	slugRetriever     SlugRetriever
	regionRetriever   RegionRetriever
	ledgerManager     LedgerManager
	bucket            BucketFetcher
	addressApi        addressv1.AddressClient
}

func NewSpentProcessor(
	spentManager SpentManager,
	campaignManager CampaignManager,
	merchantRetriever MerchantRetriever,
	slugRetriever SlugRetriever,
	regionRetriever RegionRetriever,
	ledgerManager LedgerManager,
	bucket BucketFetcher,
	addressApi addressv1.AddressClient) SpentProcessor {
	return SpentProcessor{
		spentManager:      spentManager,
		campaignManager:   campaignManager,
		merchantRetriever: merchantRetriever,
		slugRetriever:     slugRetriever,
		regionRetriever:   regionRetriever,
		ledgerManager:     ledgerManager,
		bucket:            bucket,
		addressApi:        addressApi,
	}
}

func (sp SpentProcessor) ProcessSpentEvent(ctx context.Context, spentEvent model.SpentEvent) error {
	err := spentEvent.ValidateSpentEvent()
	if err != nil {
		easyzap.Errorf("invalid spentEvent %v: %v", spentEvent, err)

		return model.ErrInvalid
	}

	campaign, err := sp.fetchCampaign(ctx, spentEvent.CampaignId)
	if err != nil {

		return err
	}

	merchant, err := sp.fetchMerchant(ctx, spentEvent.MerchantId)
	if err != nil {

		return err
	}

	address, err := sp.addressApi.GetAddressByIp(ctx, &addressv1.IpRequest{Ip: spentEvent.IP})
	if err != nil {

		return err
	}

	region, err := sp.fetchRegion(ctx, merchant.RegionId)
	if err != nil {

		return err
	}

	slug, err := sp.fetchSlug(ctx, spentEvent.SlugName)
	if err != nil {

		return err
	}

	spent, err := sp.spentManager.FetchByMerchantIdAndBucket(ctx, spentEvent.MerchantId, sp.bucket.CurrentBucket().Key)
	if err != nil && err != model.ErrNotFound {

		return err
	} else {
		spent = model.Spent{
			Id:               uuid.New(),
			CampaignId:       campaign.Id,
			MerchantId:       merchant.Id,
			Bucket:           sp.bucket.CurrentBucket().Key,
			TotalSpent:       0,
			TotalClicks:      0,
			TotalImpressions: 0,
		}
	}

	shouldPerformSpent := sp.shouldPerformSpent(
		campaign.Status,
		merchant.Status,
		slug.Status,
		region.Status,
		campaign.Budget,
		spent.TotalSpent,
		slug.Cost,
		region.Cost,
		address.Lat,
		address.Lon,
		region.Lat,
		region.Long,
	)
	if !shouldPerformSpent {

		return nil
	}

	howMuchShouldCharge, shouldRunOut := sp.howMuchShouldChargeAndShouldRunOut(campaign.Budget, spent.TotalSpent, slug.Cost, region.Cost)

	spent.TotalSpent += howMuchShouldCharge
	if spentEvent.EventType == string(model.Click) {
		spent.TotalClicks++
	} else {
		spent.TotalImpressions++
	}

	err = sp.spentManager.UpsertAndRegsterLedger(ctx, spent, spentEvent, region.Name)
	if err != nil {

		return err
	}

	if shouldRunOut {
		campaign.Status = string(model.Suspended)
		err := sp.campaignManager.Upsert(ctx, campaign)
		if err != nil {

			return err
		}
	}

	return nil
}

func (sp SpentProcessor) fetchCampaign(ctx context.Context, campaignId uuid.UUID) (model.Campaign, error) {
	campaign, err := sp.campaignManager.Fetch(ctx, campaignId)
	if err != nil {

		return model.Campaign{}, err
	}

	return campaign, nil
}

func (sp SpentProcessor) fetchMerchant(ctx context.Context, merchantId uuid.UUID) (model.Merchant, error) {
	merchant, err := sp.merchantRetriever.Fetch(ctx, merchantId)
	if err != nil {

		return model.Merchant{}, err
	}

	return merchant, nil
}

func (sp SpentProcessor) fetchSlug(ctx context.Context, slugName string) (model.Slug, error) {
	slug, err := sp.slugRetriever.Fetch(ctx, slugName)
	if err != nil {

		return model.Slug{}, err
	}

	return slug, nil
}

func (sp SpentProcessor) fetchRegion(ctx context.Context, regionId uuid.UUID) (model.Region, error) {
	region, err := sp.regionRetriever.Fetch(ctx, regionId)
	if err != nil {

		return model.Region{}, err
	}

	return region, nil
}

func (sp SpentProcessor) shouldPerformSpent(
	campaignStatus string,
	merchantStatus string,
	slugStatus string,
	regionStatus string,
	budget,
	totalSpent,
	slugCost,
	regionCost float64,
	addressLat float64,
	addressLon float64,
	regionLat float64,
	regionLong float64,
) bool {
	status := string(model.Active)
	if campaignStatus != status || merchantStatus != status || slugStatus != status ||
		regionStatus != status {
		return false
	}

	if addressLat == 0 || addressLon == 0 || regionLat == 0 || regionLong == 0 {
		return false
	}

	if geocode.IsWithinRegion(addressLat, addressLon, regionLat, regionLong) {
		return false
	}

	if (totalSpent + slugCost + regionCost) < budget {
		return true
	}

	return false
}

func (sp SpentProcessor) howMuchShouldChargeAndShouldRunOut(budget, totalSpent, slugCost, regionCost float64) (float64, bool) {
	partialClickValue := budget - (totalSpent + slugCost + regionCost)
	if partialClickValue < (slugCost+regionCost) && partialClickValue > 0 {

		return partialClickValue, true
	}

	return slugCost + regionCost, false
}
