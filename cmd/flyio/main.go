package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/wenchy/grpcio/cmd/flyio/conf"
	"github.com/wenchy/grpcio/cmd/flyio/controller"
	"github.com/wenchy/grpcio/cmd/flyio/gateway"
	"github.com/wenchy/grpcio/cmd/flyio/services"
	"github.com/wenchy/grpcio/internal/atom"

	// _ "github.com/wenchy/grpcio/internal/confpb"
	"github.com/wenchy/grpcio/internal/corepb"
	"github.com/wenchy/grpcio/internal/daemon"
	"github.com/wenchy/grpcio/internal/rpc"
	// "github.com/wenchy/grpcio/internal/rpc/resolver"
)

func main() {
	daemon.Start()
	defer daemon.Fini()

	conf.InitConf(daemon.Args.Conf) // server config
	var a = &atom.Atom{}
	a.InitZap(conf.Conf.Log.Level, conf.Conf.Log.Dir) // log
	defer a.Log.Sync()
	// a.InitRedis(env.EnvConf.Redis.Addrs, env.EnvConf.Redis.Password)
	// a.InitMysqlDB(
	// 	env.EnvConf.Mysql.Address,
	// 	env.EnvConf.Mysql.Database,
	// 	env.EnvConf.Mysql.Username,
	// 	env.EnvConf.Mysql.Password)

	a.InitGin()
	// init exported variables in package `atom`.
	atom.InitFrom(a)

	controller.InitRouter()

	// resolver.Register()

	// errgroup
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	// web server
	httpServer := &http.Server{
		Addr:    conf.Conf.Server.HTTPAddress,
		Handler: a.GinEngine,
	}
	g.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Log.Errorf("http server listen failed: %s\n", err)
			return err
		}
		return nil
	})

	// grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(atom.GetZapLogger(), rpc.ZapOpts...),
			grpc_recovery.UnaryServerInterceptor(rpc.RecoveryOpts...),
		)),
	)

	g.Go(func() error {
		// init gRPC
		lis, err := net.Listen("tcp", conf.Conf.Server.GRPCAddress)
		if err != nil {
			a.Log.Errorf("failed to listen: %v", err)
			return err
		}
		// register services
		corepb.RegisterClientTestServiceServer(grpcServer, services.NewClientTestService())

		if err := grpcServer.Serve(lis); err != nil {
			a.Log.Warnf("failed to serve: %v", err)
			return err
		}
		return nil
	})

	gwServer := &http.Server{
		Addr:    conf.Conf.Server.GatewayAddress,
		Handler: runtime.NewServeMux(),
	}
	g.Go(func() error {
		// TODO(wenchyzhu): gracefully stop gateway
		grpcAddr := conf.Conf.Server.GRPCAddress
		if err := gateway.Run(gwServer, grpcAddr); err != nil && err != http.ErrServerClosed {
			a.Log.Errorf("gateway server listen failed: %s\n", err)
			return err
		}
		return nil
	})

	daemon.ProcessSignal(
		ctx,
		func(sig os.Signal) error {
			// The context is used to inform the server it has 3 seconds to finish
			// the request it is currently handling
			// http gracefully stop
			httpCtx, httpCancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer httpCancel()
			if err := httpServer.Shutdown(httpCtx); err != nil {
				a.Log.Errorf("http server forced to shutdown:", err)
			}

			// The context is used to inform the server it has 3 seconds to finish
			// the request it is currently handling
			// http gracefully stop
			gwCtx, gwCancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer gwCancel()
			if err := gwServer.Shutdown(gwCtx); err != nil {
				a.Log.Errorf("gatewat server forced to shutdown:", err)
			}

			// grpc gracefully stop
			grpcServer.GracefulStop()

			err := g.Wait()
			if err != nil {
				a.Log.Errorf("error: %v", err)
				return err
			}
			a.Log.Info("stop success")
			return nil
		},
		func(sig os.Signal) error {
			a.Log.Info("TODO: conf reload...")
			return nil
		},
	)
}
