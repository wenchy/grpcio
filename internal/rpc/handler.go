package rpc

import (
	"runtime/debug"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/wenchy/grpcio/internal/atom"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	customRecoveryHandler grpc_recovery.RecoveryHandlerFunc
	RecoveryOpts          []grpc_recovery.Option

	customCodeToLevelHandler grpc_zap.CodeToLevel
	ZapOpts                  []grpc_zap.Option
)

func init() {
	// Define customfunc to handle panic
	customRecoveryHandler = func(p interface{}) (err error) {
		atom.Log.Error("stacktrace from panic: \n", string(debug.Stack()))
		return status.Errorf(codes.Internal, "panic triggered: %v", p)
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	RecoveryOpts = []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customRecoveryHandler),
	}

	// CodeToLevel function defines the mapping between gRPC return codes and interceptor log level.
	customCodeToLevelHandler = func(code codes.Code) zapcore.Level{
		switch code {
		case codes.OK:
			return zapcore.DebugLevel
		case codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.Unauthenticated:
			return zapcore.InfoLevel
		case codes.DeadlineExceeded, codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange, codes.Unavailable:
			return zapcore.WarnLevel
		case codes.Unknown, codes.Unimplemented, codes.Internal, codes.DataLoss:
			return zapcore.ErrorLevel
		default:
			return zapcore.ErrorLevel
		}
	}
	// Shared options for the grpc_zap, with a custom gRPC code to log level function.
	ZapOpts = []grpc_zap.Option{
		grpc_zap.WithLevels(customCodeToLevelHandler),
	}
}
