package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jiuzhou-zhao/go-templates/config"
	"github.com/jiuzhou-zhao/go-templates/grpc/gens/utpb"
	"github.com/sgostarter/libservicetoolset/servicetoolset"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type utServiceImpl struct {
	utpb.UnimplementedUTServiceServer
}

func (impl *utServiceImpl) Hello(ctx context.Context, request *utpb.HelloRequest) (*utpb.HelloResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "no request body")
	}

	return &utpb.HelloResponse{
		Reply: fmt.Sprintf("Hi %s, I accepted your message: %s", request.Caller, request.Message),
	}, nil
}

func main() {
	cfg := config.GetConfig()

	logger := cfg.Logger

	tlsConfig, err := servicetoolset.GRPCTlsConfigMap(cfg.GRPCTLSConfig4Server)
	if err != nil {
		logger.Fatal(err)
	}

	grpcCfg := &servicetoolset.GRPCServerConfig{
		Address:           cfg.GRPCListen4Server,
		TLSConfig:         tlsConfig,
		KeepAliveDuration: time.Minute * 10,
	}

	s, err := servicetoolset.NewGRPCServer(nil, grpcCfg,
		[]grpc.ServerOption{grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             time.Second * 10,
			PermitWithoutStream: true,
		})}, nil, logger)
	if err != nil {
		logger.Fatal(err)

		return
	}

	err = s.Start(func(s *grpc.Server) error {
		utpb.RegisterUTServiceServer(s, &utServiceImpl{})

		return nil
	})
	if err != nil {
		logger.Fatal(err)

		return
	}

	logger.Info("grpc server listen on: ", cfg.GRPCListen4Server)
	s.Wait()
}
