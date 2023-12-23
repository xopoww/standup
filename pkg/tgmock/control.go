//nolint:lll // grpc
package tgmock

import (
	"context"
	"fmt"

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
		if req.GetAll() || msg.GetFrom().GetId() == tm.me().GetId() {
			rsp.Messages = append(rsp.GetMessages(), msg)
		}
	}
	return rsp, nil
}

func (tm *TGMock) CreateUpdate(ctx context.Context, req *control.CreateUpdateRequest) (*control.CreateUpdateResponse, error) {
	tm.mx.Lock()
	defer tm.mx.Unlock()
	rsp := &control.CreateUpdateResponse{}

	if msg := req.GetUpdate().GetMessage(); msg != nil {
		entities, err := parseMessageEntities(msg.GetText())
		if err != nil {
			return nil, fmt.Errorf("invalid message text: %w", err)
		}
		msg.Entities = entities
		tm.addMessage(ctx, msg)
		rsp.MessageId = msg.GetMessageId()
	}

	upd := req.GetUpdate()
	upd.UpdateId = int64(len(tm.updates))
	tm.updates = append(tm.updates, upd)
	return rsp, nil
}
