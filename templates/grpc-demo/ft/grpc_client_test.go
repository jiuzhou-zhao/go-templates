package ft

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jiuzhou-zhao/go-templates/config"
	"github.com/jiuzhou-zhao/go-templates/grpc/gens/utpb"
	"github.com/sgostarter/libservicetoolset/clienttoolset"
	"github.com/sgostarter/libservicetoolset/servicetoolset"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func TestGrpc(t *testing.T) {
	cfg := config.GetConfig()

	tlsConfig, err := servicetoolset.GRPCTlsConfigMap(cfg.GRPCTLSConfig4Client)
	assert.Nil(t, err)

	grpcDial := &clienttoolset.GRPCClientConfig{
		Target:    cfg.GRPCServerAddress4Client,
		TLSConfig: tlsConfig,
	}

	conn, err := clienttoolset.DialGRPC(grpcDial, []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second, // send pings every x seconds if there is no activity
			Timeout:             1 * time.Second,  // wait x second for ping ack before considering the connection dead
			PermitWithoutStream: true,             // send pings even without active streams
		}),
		grpc.WithBlock(),
	})
	assert.Nil(t, err)

	defer conn.Close()

	client := utpb.NewUTServiceClient(conn)

	for idx := 0; idx < 10; idx++ {
		resp, err := client.Hello(context.TODO(), &utpb.HelloRequest{
			Message: fmt.Sprintf("msg-%d", idx),
			Caller:  "bitter",
		})
		assert.Nil(t, err)
		t.Log(resp.GetReply())
		time.Sleep(time.Second)
	}
}
