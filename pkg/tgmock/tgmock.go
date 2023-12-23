package tgmock

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/pkg/tgmock/control"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type TGMock struct {
	control.UnimplementedTGMockControlServer

	control *grpc.Server
	service *http.Server

	cfg    Config
	logger *zap.Logger

	mx      sync.Locker
	chats   map[int64][]*control.Message
	updates []*control.Update
}

func New(cfg Config, logger *zap.Logger) *TGMock {
	tm := &TGMock{
		cfg:    cfg,
		logger: logger,

		mx:    &sync.Mutex{},
		chats: make(map[int64][]*control.Message),
	}
	tm.initService()
	tm.initControl()
	return tm
}

func (tm *TGMock) Start() error {
	lis, err := net.Listen("tcp", tm.cfg.Control)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	ech := make(chan error)
	go func() {
		tm.logger.Sugar().Debugf("Starting control server on %q...", tm.cfg.Control)
		ech <- tm.control.Serve(lis)
	}()
	go func() {
		tm.logger.Sugar().Debugf("Starting service server on %q...", tm.cfg.Service)
		ech <- tm.service.ListenAndServe()
	}()
	return <-ech
}

func (tm *TGMock) Stop() error {
	tm.control.GracefulStop()

	const shutdownTimeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return tm.service.Shutdown(ctx)
}

func (tm *TGMock) me() *control.User {
	return &control.User{
		Id:       0,
		Username: "Test Bot",
	}
}

// addMessage must be called only from under lock
func (tm *TGMock) addMessage(ctx context.Context, msg *control.Message) {
	chatID := msg.GetChat().GetId()
	msg.MessageId = int64(len(tm.chats[chatID]) + 1)
	logging.L(ctx).Sugar().Debugf("Add message %d to chat %d.", msg.GetMessageId(), chatID)
	tm.chats[chatID] = append(tm.chats[chatID], msg)
}

// getMessage must be called only from under lock
func (tm *TGMock) getMessage(chatID, messageID int64) (*control.Message, error) {
	i := messageID - 1
	if i < 0 || i >= int64(len(tm.chats[chatID])) {
		return nil, fmt.Errorf("no message %d in chat %d", messageID, chatID)
	}
	return tm.chats[chatID][i], nil
}
