package service

import (
	"context"
	"fmt"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type CampaignDao interface {
	Fetch(context.Context, uuid.UUID) (model.Campaign, error)
	Create(context.Context, transaction.Transaction, model.Campaign) error
	Update(context.Context, transaction.Transaction, model.Campaign) error
}

type CampaignHistoryDao interface {
	Create(context.Context, transaction.Transaction, model.CampaignHistory) error
}

type SpentFetcher interface {
	FetchByMerchantIdAndBucket(ctx context.Context, id uuid.UUID, bucket string) (model.Spent, error)
}

type BucketFetcher interface {
	CurrentBucket() model.Bucket
}

type CampaignService struct {
	campaignDao        CampaignDao
	campaignHistoryDao CampaignHistoryDao
	spentFetcher       SpentFetcher
	bucketFetcher      BucketFetcher
	tm                 TransactionManager
}

func NewCampaignService(campaignDao CampaignDao, campaignHistoryDao CampaignHistoryDao, spentFetcher SpentFetcher, bucketFetcher BucketFetcher, tm TransactionManager) CampaignService {
	return CampaignService{
		campaignDao:        campaignDao,
		campaignHistoryDao: campaignHistoryDao,
		spentFetcher:       spentFetcher,
		bucketFetcher:      bucketFetcher,
		tm:                 tm,
	}
}

func (cs CampaignService) Upsert(ctx context.Context, campaign model.Campaign) error {
	err := campaign.ValidateCampaign()
	if err != nil {
		easyzap.Errorf("invalid campaign %v: %v", campaign, err)

		return model.ErrInvalid
	}

	campaignDb, err := cs.campaignDao.Fetch(ctx, campaign.Id)
	if err != nil && err != model.ErrNotFound {
		easyzap.Errorf("fail to fetch campaign by campaignId %s: %v", campaign.Id.String(), err)

		return err
	}

	if err != nil && err == model.ErrNotFound {
		err := cs.createAndRegistryHistory(ctx, campaign, &campaignDb)
		if err != nil {
			easyzap.Errorf("fail to create campaign %v: %v", campaign, err)

			return err
		}
	} else {
		shouldUpdate, shouldActivateCampaign := cs.shouldUpdateAndActivateCampaign(ctx, campaign, campaignDb)
		if !shouldUpdate {

			return model.ErrUnprocessableEntity
		}

		if shouldActivateCampaign {
			campaign.Status = string(model.Active)
		}

		err := cs.updateAndRegistryHistory(ctx, campaign, &campaignDb)
		if err != nil {
			easyzap.Errorf("fail to update campaign from %s to %s for campaignId %s: %v", campaignDb.Status, campaign.Status, campaign.Id.String(), err)

			return err
		}
	}

	return nil
}

func (cs CampaignService) createAndRegistryHistory(ctx context.Context, campaign model.Campaign, campaignDb *model.Campaign) error {
	funcWithTransaction := func(ctx context.Context, tx transaction.Transaction) error {
		err := cs.campaignDao.Create(ctx, tx, campaign)
		if err != nil {

			return err
		}

		history := cs.buildHistory(campaign, campaignDb)

		err = cs.campaignHistoryDao.Create(ctx, tx, history)
		if err != nil {

			return err
		}

		return err
	}

	return cs.tm.Execute(ctx, funcWithTransaction)
}

func (cs CampaignService) updateAndRegistryHistory(ctx context.Context, campaign model.Campaign, campaignDb *model.Campaign) error {
	funcWithTransaction := func(ctx context.Context, tx transaction.Transaction) error {
		err := cs.campaignDao.Update(ctx, tx, campaign)
		if err != nil {

			return err
		}

		history := cs.buildHistory(campaign, campaignDb)

		err = cs.campaignHistoryDao.Create(ctx, tx, history)
		if err != nil {

			return err
		}

		return err
	}

	return cs.tm.Execute(ctx, funcWithTransaction)
}

func (cs CampaignService) buildHistory(campaign model.Campaign, campaignDb *model.Campaign) model.CampaignHistory {
	history := model.CampaignHistory{
		Id:          uuid.New(),
		CampaignId:  campaign.Id,
		Status:      campaign.Status,
		Description: model.CampaignCreatedAndActive,
		CreatedBy:   campaign.UpdatedBy,
		CreatedAt:   campaign.UpdatedAt,
	}

	if campaignDb.Id == uuid.Nil {
		history.Description = model.CampaignCreatedAndActive
	} else {
		if campaignDb.Status != campaign.Status {
			history.Description = fmt.Sprintf(model.CampaignUpdateStatus, campaignDb.Status, campaign.Status)
		}
		if campaignDb.Budget != campaign.Budget {
			history.Description = fmt.Sprintf(model.CampaignUpdateBudget, campaignDb.Budget, campaign.Budget)
		}
		if campaignDb.Budget != campaign.Budget && campaignDb.Status != campaign.Status {
			history.Description = fmt.Sprintf("%v E %v",
				fmt.Sprintf(model.CampaignUpdateBudget, campaignDb.Budget, campaign.Budget),
				fmt.Sprintf(model.CampaignUpdateStatus, campaignDb.Status, campaign.Status))
		}
	}

	return history
}

func (cs CampaignService) shouldUpdateAndActivateCampaign(ctx context.Context, campaign model.Campaign, campaignDb model.Campaign) (bool, bool) {
	spent, err := cs.spentFetcher.FetchByMerchantIdAndBucket(ctx, campaign.MerchantId, cs.bucketFetcher.CurrentBucket().Key)
	if err != nil {
		if err != model.ErrNotFound {
			easyzap.Errorf("fail to fetch spent by campaign id %s: %v", campaign.Id.String(), err)

			return true, false
		}

		return false, false
	}

	if campaignDb.Status == string(model.Suspended) &&
		campaign.Status == string(model.Active) &&
		campaign.Budget <= campaignDb.Budget &&
		campaign.Budget <= spent.TotalSpent {

		return false, false

	}

	if campaignDb.Status == string(model.Suspended) &&
		campaign.Budget > campaignDb.Budget &&
		campaign.Budget > spent.TotalSpent {

		return true, true

	}

	return true, false
}

func (cs CampaignService) Fetch(ctx context.Context, campaignId uuid.UUID) (model.Campaign, error) {
	return cs.campaignDao.Fetch(ctx, campaignId)
}
