package handler

import (
	"context"

	"github.com/gorilla/websocket"
)

type WebSocketHandle struct {
	connect *websocket.Conn
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewSocket(conn *websocket.Conn) *WebSocketHandle {
	return NewSocketWithContext(context.Background(), conn)
}

func NewSocketWithContext(ctx context.Context, conn *websocket.Conn) *WebSocketHandle {
	ctx, cancel := context.WithCancel(ctx)

	return &WebSocketHandle{
		connect: conn,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (h *WebSocketHandle) AwaitClose(codes ...int) {
	go func() {
		for {
			_, _, err := h.connect.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, codes...) {
					h.cancel()
					return
				}
			}
		}
	}()
}

func (h *WebSocketHandle) Done() <-chan struct{} {
	return h.ctx.Done()
}

func (h *WebSocketHandle) Exit() {
	h.connect.Close()
}

func (h *WebSocketHandle) WriteJSON(v any) error {

	defer func() {
		if recover() != nil {
			return
		}
	}()

	return h.connect.WriteJSON(v)
}
