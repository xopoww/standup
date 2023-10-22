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

func marshalJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("{!err %s}", err.Error())
	}
	return string(data)
}
