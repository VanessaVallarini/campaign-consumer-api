package base_grpc

import (
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func GetGRPCConnection(config config.ClientConfig) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		config.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Duration(config.KeepAlive.TimeMs) * time.Millisecond,
			Timeout:             time.Duration(config.KeepAlive.TimeoutMs) * time.Millisecond,
			PermitWithoutStream: false,
		}),
	)
	return conn, err
}
