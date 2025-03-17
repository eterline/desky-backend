package handler

import (
	"context"
	"encoding/base64"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketProtoUpgrader interface {
	Upgrade(http.ResponseWriter, *http.Request, http.Header) (WebSocketConnector, error)
}

type WebSocketConnector interface {
	Close() error

	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	NextReader() (messageType int, r io.Reader, err error)
	NextWriter(messageType int) (io.WriteCloser, error)

	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type WebSocketHandler struct {
	upgrade *websocket.Upgrader
	logger  WebSocketLogger

	ctx context.Context
}

func NewWebSocketHandler(ctx context.Context, u *websocket.Upgrader) *WebSocketHandler {
	return &WebSocketHandler{
		ctx:     ctx,
		upgrade: u,
	}
}

type WebSocketSession struct {
	ID uuid.UUID

	conn WebSocketConnector
	log  WebSocketLogger

	ctx    context.Context
	cancel context.CancelFunc
}

func (h *WebSocketHandler) HandleConnect(w http.ResponseWriter, r *http.Request) (*WebSocketSession, error) {

	uuid := uuid.New()

	head := make(http.Header)
	head.Set("UUID", uuid.String())

	conn, err := h.upgrade.Upgrade(w, r, head)
	if err != nil {
		return nil, err
	}

	if h.logger != nil {
		h.logger.Infof("new socket connection: %s uuid: %s", conn.RemoteAddr().String(), uuid.String())
	}

	context, cancel := context.WithCancel(h.ctx)

	session := &WebSocketSession{
		ID:     uuid,
		conn:   conn,
		log:    h.logger,
		ctx:    context,
		cancel: cancel,
	}

	return session, nil
}

func (h *WebSocketSession) AwaitClose(codes ...int) {
	go func() {

		defer func() {
			if h.log != nil {
				h.log.Infof("socket reader closed uuid: %s", h.ID.String())
			}
			h.conn.Close()
		}()

		if h.conn == nil {
			if h.log != nil {
				h.log.Errorf("socket uuid: %s connect is nil", h.ID.String())
			}
			h.cancel()
			return
		}

		for {
			for {
				select {

				case <-h.ctx.Done():
					return

				default:
					_, _, err := h.conn.ReadMessage()

					if err != nil {
						if h.log != nil && !websocket.IsCloseError(err, codes...) {
							h.log.Errorf("socket uuid: %s error: %v", h.ID.String(), err)
						}
						h.cancel()
						return
					}
				}
			}
		}
	}()
}

func (h *WebSocketSession) AwaitMessage(closeCodes ...int) <-chan SocketMessage {
	channel := make(chan SocketMessage, 1)

	go func() {

		defer func() {
			if h.log != nil {
				h.log.Infof("socket reader closed uuid: %s", h.ID.String())
			}
			close(channel)
		}()

		if h.conn == nil {
			if h.log != nil {
				h.log.Errorf("socket uuid: %s connect is nil", h.ID.String())
			}
			h.cancel()
			return
		}

		for {
			select {

			case <-h.ctx.Done():
				return

			default:
				msgType, bytes, err := h.conn.ReadMessage()

				if err != nil {
					if h.log != nil && !websocket.IsCloseError(err, closeCodes...) {
						h.log.Errorf("socket uuid: %s error: %v", h.ID.String(), err)
					}
					h.conn.Close()
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

func (h *WebSocketSession) SessionDone() <-chan struct{} {
	return h.ctx.Done()
}

func (h *WebSocketSession) WriteBytes(msg []byte) error {
	return h.conn.WriteMessage(websocket.BinaryMessage, msg)
}

func (h *WebSocketSession) WriteBase64(msg []byte) error {

	encodedSize := base64.StdEncoding.EncodedLen(len(msg))
	buf := make([]byte, encodedSize)

	base64.StdEncoding.Encode(buf, msg)

	return h.conn.WriteMessage(websocket.TextMessage, buf)
}

func (h *WebSocketSession) Exit() error {
	h.cancel()
	return h.conn.Close()
}

type WebSocketWriting struct {
	typeMessage int
	conn        WebSocketConnector

	mu sync.Mutex
}

func (h *WebSocketSession) InitWebSocketWriting(bin bool) *WebSocketWriting {

	var typeMessage int

	if bin {
		typeMessage = websocket.BinaryMessage
	} else {
		typeMessage = websocket.TextMessage
	}

	return &WebSocketWriting{
		typeMessage: typeMessage,
		conn:        h.conn,
	}
}

func (wr *WebSocketWriting) Write(p []byte) (n int, err error) {

	wr.mu.Lock()
	defer wr.mu.Unlock()

	if err := wr.conn.WriteMessage(wr.typeMessage, p); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (wr *WebSocketWriting) CloseWriting() error {

	wr.mu.Lock()
	defer wr.mu.Unlock()

	return wr.conn.Close()
}

type WebSocketBase64Writing struct {
	conn WebSocketConnector

	mu sync.Mutex
}

func (h *WebSocketSession) InitWebSocketBase64Writing() *WebSocketBase64Writing {
	return &WebSocketBase64Writing{
		conn: h.conn,
	}
}

func (wrb *WebSocketBase64Writing) Write(p []byte) (n int, err error) {

	wrb.mu.Lock()
	defer wrb.mu.Unlock()

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(p)))

	base64.StdEncoding.Encode(buf, p)

	if err := wrb.conn.WriteMessage(websocket.TextMessage, buf); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (wrb *WebSocketBase64Writing) CloseWriting() error {

	wrb.mu.Lock()
	defer wrb.mu.Unlock()

	return wrb.conn.Close()
}
