package redis

import (
	"context"
	"fmt"
	"time"

	easyzap "github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
	redistracer "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"

	client "github.com/redis/go-redis/v9"
)

type RedisDataClient struct {
	client *client.ClusterClient
}

func NewRedisDataClient(rwAddress string, roAddress string) RedisDataClient {
	clusterSlots := func(ctx context.Context) ([]client.ClusterSlot, error) {
		nodes := []client.ClusterNode{{Addr: rwAddress}, {Addr: roAddress}}

		const maxRedisHashSlot = 16383
		slots := []client.ClusterSlot{
			{
				Start: 0,
				End:   maxRedisHashSlot,
				Nodes: nodes,
			},
		}

		easyzap.Debug(ctx, fmt.Sprintf("solved %v", slots))
		return slots, nil
	}

	redisClient := client.NewClusterClient(&client.ClusterOptions{
		ClusterSlots:  clusterSlots,
		ReadOnly:      true,
		RouteRandomly: true,
	})

	redistracer.WrapClient(redisClient)

	return RedisDataClient{client: redisClient}
}

func (rdc RedisDataClient) Ping(ctx context.Context) {
	err := rdc.client.Ping(ctx).Err()
	if err != nil {
		easyzap.Fatal(err, "failed to reach redis")
	}
}

func (rdc RedisDataClient) Get(ctx context.Context, key string) (string, error) {
	out, err := rdc.client.Get(ctx, key).Result()
	if err != nil {
		if err == client.Nil {
			return "", nil
		}

		return "", errors.Wrapf(err, "failed to GET item of key %s on redis", key)
	}

	return out, nil
}

func (rdc RedisDataClient) SetIfNotExists(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error) {
	created, err := rdc.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, errors.Wrapf(err, "failed to set item of key %s on redis", key)
	}

	return created, nil
}
