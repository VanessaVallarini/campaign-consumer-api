package service

import (
	"context"
	"fmt"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type CampaignDao interface {
	Fetch(context.Context, uuid.UUID) (model.Campaign, error)
	Create(context.Context, model.Campaign) error
	Update(context.Context, model.Campaign) error
}

type CampaignHistoryDao interface {
	Create(context.Context, model.CampaignHistory) error
}

type SpentFetcher interface {
	FetchByCampaignIdAndBucket(ctx context.Context, id uuid.UUID, bucket string) (model.Spent, error)
}

type BucketFetcher interface {
	CurrentBucket() model.Bucket
}

type CampaignService struct {
	campaignDao        CampaignDao
	campaignHistoryDao CampaignHistoryDao
	spentFetcher       SpentFetcher
	bucketFetcher      BucketFetcher
}

func NewCampaignService(campaignDao CampaignDao, campaignHistoryDao CampaignHistoryDao, spentFetcher SpentFetcher, bucketFetcher BucketFetcher) CampaignService {
	return CampaignService{
		campaignDao:        campaignDao,
		campaignHistoryDao: campaignHistoryDao,
		spentFetcher:       spentFetcher,
		bucketFetcher:      bucketFetcher,
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
		err := cs.campaignDao.Create(ctx, campaign)
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

		err = cs.campaignDao.Update(ctx, campaign)
		if err != nil {
			easyzap.Errorf("fail to update campaignDb %v to campaign %v: %v", campaignDb, campaign, err)

			return err
		}
	}

	err = cs.registryHistory(ctx, campaign, &campaignDb)
	if err != nil {
		easyzap.Errorf("fail to registry history campaignDb %v to campaign %v: %v", campaignDb, campaign, err)
		campaign.Status = string(model.Cancelled)
		errRollback := cs.campaignDao.Update(ctx, campaignDb)
		if errRollback != nil {
			easyzap.Errorf("[INCONSISTENT] fail to rollback campaign %v to campaignDb %v: %v", campaign, campaignDb, err)

			return err
		}

		return err
	}

	return nil
}

func (cs CampaignService) shouldUpdateAndActivateCampaign(ctx context.Context, campaign model.Campaign, campaignDb model.Campaign) (bool, bool) {
	//shouldUpdate
	//shouldActivateCampaign

	spent, err := cs.spentFetcher.FetchByCampaignIdAndBucket(ctx, campaign.Id, cs.bucketFetcher.CurrentBucket().Key)
	if err != nil {
		if err != model.ErrNotFound {

			return true, false
		}

		easyzap.Errorf("fail to fetch spent by campaign id %s: %v", campaign.Id.String(), err)

		return false, false
	}

	//não pode atualizar o status de uma campanha de pendente para ativo se
	//o budget alterar pra baixo ou não sofrer alteração e for menor/igual ao spent
	if campaignDb.Status == string(model.Suspended) &&
		campaign.Status == string(model.Active) &&
		campaign.Budget <= campaignDb.Budget &&
		campaign.Budget <= spent.TotalSpent {
		return false, false

	}

	//uma campanha deve ser ativa se estiver suspensa e o budget for maior que o spent
	if campaignDb.Status == string(model.Suspended) &&
		campaign.Budget > campaignDb.Budget &&
		campaign.Budget > spent.TotalSpent {
		return true, true

	}

	return true, false
}

func (cs CampaignService) registryHistory(ctx context.Context, campaign model.Campaign, campaignDb *model.Campaign) error {
	if campaignDb.Id == uuid.Nil {
		err := cs.campaignHistoryDao.Create(ctx, model.CampaignHistory{
			Id:          uuid.New(),
			CampaignId:  campaign.Id,
			Status:      campaign.Status,
			Description: model.CampaignCreatedAndActive,
			CreatedBy:   campaign.UpdatedBy,
			CreatedAt:   campaign.UpdatedAt,
		})
		if err != nil {
			easyzap.Errorf("fail to registry history campaign create %v: %v", campaign, err)

			return err
		}
	} else {
		if campaignDb.Status != campaign.Status {
			err := cs.campaignHistoryDao.Create(ctx, model.CampaignHistory{
				Id:          uuid.New(),
				CampaignId:  campaign.Id,
				Status:      campaign.Status,
				Description: fmt.Sprintf(model.CampaignUpdateStatus, campaignDb.Status, campaign.Status),
				CreatedBy:   campaign.UpdatedBy,
				CreatedAt:   campaign.UpdatedAt,
			})
			if err != nil {
				easyzap.Errorf("fail to registry history campaign status from %s to %s for campaignId %s: %v", campaignDb.Status, campaign.Status, campaign.Id.String(), err)

				return err
			}
		}

		if campaignDb.Budget != campaign.Budget {
			err := cs.campaignHistoryDao.Create(ctx, model.CampaignHistory{
				Id:          uuid.New(),
				CampaignId:  campaign.Id,
				Status:      campaign.Status,
				Description: fmt.Sprintf(model.CampaignUpdateBudget, campaignDb.Budget, campaign.Budget),
				CreatedBy:   campaign.UpdatedBy,
				CreatedAt:   campaign.UpdatedAt,
			})
			if err != nil {
				easyzap.Errorf("fail to registry history campaign budget from %.2f to %.2f for campaignId %s: %v", campaignDb.Budget, campaign.Budget, campaign.Id.String(), err)

				return err
			}
		}
	}

	return nil
}
