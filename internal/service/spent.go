package service

import (
	"context"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
	"github.com/google/uuid"
)

type SpentDao interface {
	FetchByMerchantIdAndBucket(context.Context, uuid.UUID, string) (model.Spent, error)
	Upsert(context.Context, transaction.Transaction, model.Spent) error
}

type LedgerManager interface {
	Create(context.Context, transaction.Transaction, model.Ledger) error
}

type SpentService struct {
	spentDao      SpentDao
	tm            TransactionManager
	ledgerManager LedgerManager
}

func NewSpentService(
	spentDao SpentDao,
	tm TransactionManager,
	ledgerManager LedgerManager) SpentService {
	return SpentService{
		spentDao:      spentDao,
		tm:            tm,
		ledgerManager: ledgerManager,
	}
}

func (ss SpentService) UpsertAndRegsterLedger(ctx context.Context, spent model.Spent, spentEvent model.SpentEvent, regionName string) error {
	funcWithTransaction := func(ctx context.Context, tx transaction.Transaction) error {
		err := ss.spentDao.Upsert(ctx, tx, spent)
		if err != nil {

			return err
		}
		ledger := ss.buildLedger(spent, spentEvent, regionName)

		err = ss.ledgerManager.Create(ctx, tx, ledger)
		if err != nil {

			return err
		}

		return err
	}

	return ss.tm.Execute(ctx, funcWithTransaction)
}

func (ss SpentService) buildLedger(spent model.Spent, spentEvent model.SpentEvent, regionName string) model.Ledger {
	return model.Ledger{
		Id:         uuid.New(),
		SpentId:    spent.Id,
		CampaignId: spentEvent.CampaignId,
		MerchantId: spentEvent.MerchantId,
		SlugName:   spentEvent.SlugName,
		RegionName: regionName,
		UserId:     spentEvent.UserId,
		SessionId:  spentEvent.SessionId,
		EventType:  model.EventTypeFromString(spentEvent.EventType),
		Cost:       spent.TotalSpent,
		Ip:         spentEvent.IP,
		Lat:        0,
		Long:       0,
		CreatedAt:  time.Now(),
		EventTime:  spentEvent.EventTime,
	}
}

func (ss SpentService) FetchByMerchantIdAndBucket(ctx context.Context, id uuid.UUID, bucket string) (model.Spent, error) {
	return ss.spentDao.FetchByMerchantIdAndBucket(ctx, id, bucket)
}
