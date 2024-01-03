package tgmock

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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
	var req struct {
		ChatID           int64  `schema:"chat_id" validate:"required"`
		ReplyToMessageID *int64 `schema:"reply_to_message_id"`
		Text             string `schema:"text"`
	}
	if err := decode(r.Form, &req); err != nil {
		return nil, err
	}

	tm.mx.Lock()
	defer tm.mx.Unlock()

	msg := &control.Message{
		From: tm.me(),
		Date: time.Now().Unix(),
		Chat: &control.Chat{
			Id: req.ChatID,
		},
		Text: req.Text,
	}
	if req.ReplyToMessageID != nil {
		replyToMessage, err := tm.getMessage(req.ChatID, *req.ReplyToMessageID)
		if err != nil {
			return nil, err
		}
		msg.ReplyToMessage = replyToMessage
	}
	tm.addMessage(r.Context(), msg)

	return msg, nil
}

func (tm *TGMock) handleGetUpdates(r *http.Request) (any, error) {
	var req struct {
		Offset  int64 `schema:"offset" validate:"gte=0"`
		Limit   int64 `schema:"limit"`
		Timeout int64 `schema:"timeout"`
	}
	if err := decode(r.Form, &req); err != nil {
		return nil, err
	}

	const defaultLimit = 100
	if req.Limit == 0 {
		req.Limit = defaultLimit
	}

	rsp := make([]*control.Update, 0)
	start := time.Now()
	for {
		tm.mx.Lock()
		if req.Offset < tm.lastOffset {
			req.Offset = tm.lastOffset
		} else {
			tm.lastOffset = req.Offset
		}

		for uid := req.Offset; uid < int64(len(tm.updates)) && len(rsp) < int(req.Limit); uid++ {
			rsp = append(rsp, tm.updates[uid])
		}
		tm.mx.Unlock()

		if len(rsp) > 0 || time.Since(start) > time.Second*time.Duration(req.Limit) {
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

func decode(v url.Values, dst any) error {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	if err := d.Decode(dst, v); err != nil {
		return err
	}
	return validator.New(validator.WithRequiredStructEnabled()).Struct(dst)
}
