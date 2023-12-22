package tgmock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/xopoww/standup/internal/common/logging"
	"github.com/xopoww/standup/pkg/tgmock/control"
)

func (tm *TGMock) initService() {
	handler := mux.NewRouter()

	router := handler.Headers("Content-Type", "application/x-www-form-urlencoded").Subrouter()

	router.Path("/{token}/getMe").HandlerFunc(tm.handler(tm.handleGetMe))
	router.Path("/{token}/sendMessage").HandlerFunc(tm.handler(tm.handleSendMessage))
	router.Path("/{token}/getUpdates").HandlerFunc(tm.handler(tm.handleGetUpdates))

	tm.service = &http.Server{
		Addr:              tm.cfg.Service,
		Handler:           handler,
		ReadHeaderTimeout: time.Second,
	}
}

func (tm *TGMock) handleGetMe(_ *http.Request) (any, error) {
	return tm.me(), nil
}

func (tm *TGMock) handleSendMessage(r *http.Request) (any, error) {
	chatID, err := mustGetInt(r.Form, "chat_id")
	if err != nil {
		return nil, err
	}

	tm.mx.Lock()
	defer tm.mx.Unlock()

	msg := &control.Message{
		From: tm.me(),
		Date: time.Now().Unix(),
		Chat: &control.Chat{
			Id: chatID,
		},
		ReplyToMessage: &control.Message{},
		Text:           r.Form.Get("text"),
	}
	if replyToMessageID, err := mustGetInt(r.Form, "reply_to_message_id"); err == nil {
		msg.ReplyToMessage = &control.Message{MessageId: replyToMessageID}
	} else if !errors.Is(err, ErrNoField) {
		return nil, err
	}
	tm.addMessage(r.Context(), msg)

	return msg, nil
}

func (tm *TGMock) handleGetUpdates(r *http.Request) (any, error) {
	offset, err := mustGetInt(r.Form, "offset")
	if err != nil {
		if errors.Is(err, ErrNoField) {
			offset = 0
		} else {
			return nil, err
		}
	}
	if offset < 0 {
		return nil, fmt.Errorf("negative offset is not supported")
	}

	limit, err := mustGetInt(r.Form, "limit")
	if err != nil {
		if errors.Is(err, ErrNoField) {
			limit = 100
		} else {
			return nil, err
		}
	}

	timeout, err := mustGetInt(r.Form, "timeout")
	if err != nil {
		if errors.Is(err, ErrNoField) {
			timeout = 0
		} else {
			return nil, err
		}
	}

	rsp := make([]*control.Update, 0)
	start := time.Now()
	for {
		tm.mx.Lock()
		for uid := offset; uid < int64(len(tm.updates)) && len(rsp) < int(limit); uid++ {
			rsp = append(rsp, tm.updates[uid])
		}
		tm.mx.Unlock()

		if len(rsp) > 0 || time.Since(start) > time.Second*time.Duration(timeout) {
			break
		}
		time.Sleep(time.Second)
	}

	return rsp, nil
}

func (tm *TGMock) handler(f func(r *http.Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = logging.WithLogger(ctx, tm.logger)

		if err := func() error {
			if err := r.ParseForm(); err != nil {
				return err
			}

			logging.L(ctx).Sugar().Debugf("%s %q %s", r.Method, r.URL.Path, logging.MarshalJSON(r.Form))

			rsp, err := f(r.WithContext(ctx))
			if err != nil {
				return err
			}
			writeJSONResponse(ctx, w, map[string]any{"ok": true, "result": rsp})
			return nil
		}(); err != nil {
			logging.L(ctx).Sugar().Errorf("Handle: %s.", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func writeJSONResponse(ctx context.Context, w http.ResponseWriter, body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		logging.L(ctx).Sugar().Errorf("Failed to marshal JSON response: %s.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		logging.L(ctx).Sugar().Errorf("Failed to write response body: %s.", err)
	}
}

var ErrNoField = errors.New("missing required field")

func mustGet(v url.Values, key string) (string, error) {
	if !v.Has(key) {
		return "", fmt.Errorf("%w %q", ErrNoField, key)
	}
	return v.Get(key), nil
}

func mustGetInt(v url.Values, key string) (int64, error) {
	s, err := mustGet(v, key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}
