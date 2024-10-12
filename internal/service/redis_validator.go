package service

import (
	"time"

	"golang.org/x/net/context"
)

type CacheClient interface {
	SetIfNotExists(context.Context, string, []byte, time.Duration) (bool, error)
	Get(context.Context, string) (string, error)
}

type Validator struct {
	cacheClient CacheClient
	ttl         time.Duration
}

func NewRedisValidator(client CacheClient, ttl time.Duration) *Validator {
	return &Validator{
		cacheClient: client,
		ttl:         ttl,
	}
}

func (v Validator) SetIfNotExists(ctx context.Context, uniqueKey string, value []byte) (bool, error) {
	createOperationWasPerformed, err := v.cacheClient.SetIfNotExists(ctx, uniqueKey, []byte{}, v.ttl)
	if err != nil {

		return false, err
	}

	return createOperationWasPerformed, nil
}

func (v Validator) Get(ctx context.Context, key string) (string, error) {
	value, err := v.cacheClient.Get(ctx, key)
	if err != nil {

		return value, err
	}

	return value, nil
}
