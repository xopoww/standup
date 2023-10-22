package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryInterceptor(l *zap.Logger) grpc.UnaryServerInterceptor {
	s := l.Sugar()
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = WithLogger(ctx, l)

		s.Debugf("GRPC Req %s: %s.", info.FullMethod, marshalJSON(req))

		start := time.Now()
		resp, err = handler(ctx, req)
		delta := time.Since(start)

		if err == nil {
			s.Debugf("GRPC Rsp %s (%s): %s.", info.FullMethod, delta, marshalJSON(resp))
		} else {
			s.Errorf("GRPC Err %s (%s): %s.", info.FullMethod, delta, err)
		}

		return resp, err
	}
}

func UnaryClientInterceptor(l *zap.Logger) grpc.UnaryClientInterceptor {
	s := l.Sugar()
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = WithLogger(ctx, l)

		s.Debugf("GRPC Req %s to %s: %s.", method, cc.Target(), marshalJSON(req))

		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		delta := time.Since(start)

		if err == nil {
			s.Debugf("GRPC Rsp %s (%s): %s.", method, delta, marshalJSON(reply))
		} else {
			s.Errorf("GRPC Err %s (%s): %s.", method, delta, err)
		}

		return err
	}
}

func marshalJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("{!err %s}", err.Error())
	}
	return string(data)
}
