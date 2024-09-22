package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
)

type TransactionManager interface {
	Execute(context.Context, func(context.Context, transaction.Transaction) error) error
}
