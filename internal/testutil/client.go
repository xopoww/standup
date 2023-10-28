package testutil

import (
	"context"
	"errors"
	"strings"

	"github.com/stretchr/testify/mock"
	"github.com/xopoww/standup/pkg/api/standup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type MockStandupClient struct {
	mock.Mock
}

func (c *MockStandupClient) called(method string, args ...any) (result mock.Arguments, err error) {
	func() {
		defer func() {
			if i := recover(); i != nil {
				s, ok := i.(string)
				if !ok || !strings.Contains(s, "I don't know what to return because the method call was unexpected") {
					panic(i)
				}
				err = errors.New("unexpected mock call")
			}
		}()
		result = c.MethodCalled(method, args...)
	}()
	return
}

func (c *MockStandupClient) CreateMessage(ctx context.Context, in *standup.CreateMessageRequest, opts ...grpc.CallOption) (*standup.CreateMessageResponse, error) {
	result, err := c.called("CreateMessage", ctx, in)
	if err != nil {
		return nil, err
	}
	return result.Get(0).(*standup.CreateMessageResponse), result.Error(1)
}

func (c *MockStandupClient) GetMessage(ctx context.Context, in *standup.GetMessageRequest, opts ...grpc.CallOption) (*standup.GetMessageResponse, error) {
	result, err := c.called("GetMessage", ctx, in)
	if err != nil {
		return nil, err
	}
	return result.Get(0).(*standup.GetMessageResponse), result.Error(1)
}

func (c *MockStandupClient) ListMessages(ctx context.Context, in *standup.ListMessagesRequest, opts ...grpc.CallOption) (*standup.ListMessagesResponse, error) {
	result, err := c.called("ListMessages", ctx, in)
	if err != nil {
		return nil, err
	}
	return result.Get(0).(*standup.ListMessagesResponse), result.Error(1)
}

// OutgoingMetadata in fact returns mock.argumentMatcher and should only be used
// when preparing mock.Mock for a method call
func OutgoingMetadata(f func(md metadata.MD) bool) any {
	return mock.MatchedBy(func(ctx context.Context) bool {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			return false
		}
		return f(md)
	})
}
