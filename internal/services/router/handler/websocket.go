package handler

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsLogger interface {
	Info(args ...interface{})
	Error(args ...interface{})

	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type WsHandlerWrap struct {
	uuid uuid.UUID

	ctx    context.Context
	cancel context.CancelFunc

	logger WsLogger

	Connect *websocket.Conn
}

func NewSocket(conn *websocket.Conn) *WsHandlerWrap {
	return NewSocketWithContext(context.Background(), conn, nil)
}

func NewSocketWithContext(ctx context.Context, conn *websocket.Conn, logger WsLogger) *WsHandlerWrap {
	ctx, cancel := context.WithCancel(ctx)
	uuid := uuid.New()

	if logger != nil {
		logger.Infof("new socket connection: %s uuid: %s", conn.RemoteAddr().String(), uuid.String())
	}

	return &WsHandlerWrap{
		uuid:    uuid,
		ctx:     ctx,
		cancel:  cancel,
		Connect: conn,
	}
}

func (h *WsHandlerWrap) UUID() uuid.UUID {
	return h.uuid
}

func (h *WsHandlerWrap) AwaitClose(codes ...int) {
	go func() {

		defer func() {
			if h.logger != nil {
				h.logger.Infof("socket reader closed uuid: %s", h.uuid.String())
			}
		}()

		if h.Connect == nil {
			if h.logger != nil {
				h.logger.Errorf("socket uuid: %s connect is nil", h.uuid.String())
			}
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
						if h.logger != nil && !websocket.IsCloseError(err, codes...) {
							h.logger.Errorf("socket uuid: %s error: %v", h.uuid.String(), err)
						}
						h.Connect.Close()
						h.cancel()
						return
					}
				}
			}
		}
	}()
}

type SocketMessage struct {
	Type int
	Body []byte
}

func (m *SocketMessage) String() string {
	return string(m.Body)
}

func (h *WsHandlerWrap) AwaitMessage(closeCodes ...int) <-chan SocketMessage {
	channel := make(chan SocketMessage, 10)

	go func() {

		defer func() {
			if h.logger != nil {
				h.logger.Infof("socket reader closed uuid: %s", h.uuid.String())
			}
			close(channel)
		}()

		if h.Connect == nil {
			if h.logger != nil {
				h.logger.Errorf("socket uuid: %s connect is nil", h.uuid.String())
			}
			h.cancel()
			return
		}

		for {
			select {

			case <-h.Done():
				return

			default:
				msgType, bytes, err := h.Connect.ReadMessage()

				if err != nil {
					if h.logger != nil && !websocket.IsCloseError(err, closeCodes...) {
						h.logger.Errorf("socket uuid: %s error: %v", h.uuid.String(), err)
					}
					h.Connect.Close()
					h.cancel()
					return
				}

				channel <- SocketMessage{
					Type: msgType,
					Body: bytes,
				}
			}
		}
	}()

	return channel
}

func (h *WsHandlerWrap) Done() <-chan struct{} {
	return h.ctx.Done()
}

func (h *WsHandlerWrap) Exit() {
	h.Connect.Close()
}

func (h *WsHandlerWrap) WriteBytes(msg []byte) error {
	return h.Connect.WriteMessage(websocket.BinaryMessage, msg)
}

func (h *WsHandlerWrap) WriteBase64(msg []byte) error {

	encodedSize := base64.StdEncoding.EncodedLen(len(msg))
	buf := make([]byte, encodedSize)

	base64.StdEncoding.Encode(buf, msg)

	return h.Connect.WriteMessage(websocket.TextMessage, buf)
}

func (h *WsHandlerWrap) WriteJSON(v any) error {
	return h.Connect.WriteJSON(v)
}

func (h *WsHandlerWrap) WriteError(err error) error {

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

func (h *WsHandlerWrap) WriteCloseError(err error) {

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
