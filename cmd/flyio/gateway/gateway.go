package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/wenchy/grpcio/internal/atom"
	"github.com/wenchy/grpcio/internal/corepb"
	"google.golang.org/grpc"
)

// Run runs the gRPC-Gateway, dialling the provided address.
func Run(gwServer *http.Server, grpcAddr string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// mux := runtime.NewServeMux()
	mux, ok := gwServer.Handler.(*runtime.ServeMux)
	if !ok {
		return fmt.Errorf("handler is not *runtime.ServeMux")
	}
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := corepb.RegisterClientTestServiceGWFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	atom.Log.Infof("gateway server start: %s", gwServer.Addr)
	return gwServer.ListenAndServe()
}
