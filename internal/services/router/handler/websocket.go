package handler

import (
	"context"

	"github.com/gorilla/websocket"
)

type WebSocketHandle struct {
	Connect *websocket.Conn
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewSocket(conn *websocket.Conn) *WebSocketHandle {
	return NewSocketWithContext(context.Background(), conn)
}

func NewSocketWithContext(ctx context.Context, conn *websocket.Conn) *WebSocketHandle {
	ctx, cancel := context.WithCancel(ctx)

	return &WebSocketHandle{
		Connect: conn,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (h *WebSocketHandle) AwaitClose(codes ...int) {
	go func() {

		if h.Connect == nil {
			h.cancel()
			return
		}

		for {
			for {
				select {

				case <-h.Done():
					return

				default:
					_, _, err := h.Connect.ReadMessage()

					if err != nil {
						h.Connect.Close()
						h.cancel()
						return
					}
				}
			}
		}
	}()
}

func (h *WebSocketHandle) AwaitMessage(v any) <-chan any {

	channel := make(chan any, 1)

	go func() {

		defer close(channel)

		for {
			select {

			case <-h.Done():
				return

			default:
				err := h.Connect.ReadJSON(v)
				if err != nil {
					<-h.ctx.Done()
					return
				}
				channel <- v
			}
		}
	}()

	return channel
}

func (h *WebSocketHandle) Done() <-chan struct{} {
	return h.ctx.Done()
}

func (h *WebSocketHandle) Exit() {
	h.Connect.Close()
}

func (h *WebSocketHandle) WriteJSON(v any) error {

	defer func() {
		if recover() != nil {
			return
		}
	}()

	return h.Connect.WriteJSON(v)
}

func (h *WebSocketHandle) WriteError(err error) error {

	defer func() {
		if recover() != nil {
			return
		}
	}()

	e := NewWSHandlerError(
		err.Error(),
		false,
	)

	return h.WriteJSON(e)
}

func (h *WebSocketHandle) WriteCloseError(err error) {

	defer func() {
		if recover() != nil {
			return
		}
	}()

	e := NewWSHandlerError(
		err.Error(),
		true,
	)

	h.WriteJSON(e)
	h.Exit()
}
