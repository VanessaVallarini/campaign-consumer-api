package client

import (
	"context"

	addressv1 "github.com/VanessaVallarini/address-api/pkg/api/proto/v1"
	v1 "github.com/VanessaVallarini/address-api/pkg/api/proto/v1"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	client "github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/base_grpc"
	easyzap "github.com/lockp111/go-easyzap"
)

func NewAddressClient(config config.ClientConfig) addressv1.AddressClient {
	conn, err := client.GetGRPCConnection(config)
	if err != nil {
		easyzap.Error(context.Background(), err, "error loading address-api grpc client")

		return nil
	}
	return v1.NewAddressClient(conn)
}
