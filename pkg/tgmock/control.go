//nolint:lll // grpc
package tgmock

import (
	"context"

	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/pkg/tgmock/control"
	"google.golang.org/grpc"
)

func (tm *TGMock) initControl() {
	tm.control = grpc.NewServer(
		grpc.UnaryInterceptor(logging.UnaryInterceptor(tm.logger)),
	)
	control.RegisterTGMockControlServer(tm.control, tm)
}

func (tm *TGMock) ListMessages(_ context.Context, req *control.ListMessagesRequest) (*control.ListMessagesResponse, error) {
	tm.mx.Lock()
	defer tm.mx.Unlock()
	rsp := &control.ListMessagesResponse{}
	for _, msg := range tm.chats[req.GetChatId()] {
		if msg.GetFrom().GetId() == tm.me().GetId() {
			rsp.Messages = append(rsp.GetMessages(), msg)
		}
	}
	return rsp, nil
}

func (tm *TGMock) CreateUpdate(ctx context.Context, req *control.CreateUpdateRequest) (*control.CreateUpdateResponse, error) {
	tm.mx.Lock()
	defer tm.mx.Unlock()
	if msg := req.GetUpdate().GetMessage(); msg != nil {
		tm.addMessage(ctx, msg)
	}
	upd := req.GetUpdate()
	upd.UpdateId = int64(len(tm.updates))
	tm.updates = append(tm.updates, upd)
	return &control.CreateUpdateResponse{}, nil
}
