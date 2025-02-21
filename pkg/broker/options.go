package broker

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type ClientOptions struct {
	*mqtt.ClientOptions
	QoS byte
}

type OptionFunc func(*ClientOptions) // wrapper for options setup

func OptionCredentials(user, pwd string) OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.Username = user
		so.ClientOptions.Password = pwd
	}
}

type SenderProto string

const (
	ProtoTCP SenderProto = "tcp"
	ProtoWS  SenderProto = "ws"
	ProtoWSS SenderProto = "wss"
	ProtoSSL SenderProto = "ssl"
)

func OptionServer(proto SenderProto, host string, port uint16) OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.AddBroker(fmt.Sprintf("%s://%s:%d", proto, host, port))
	}
}

func OptionInjectCAFile(file string) OptionFunc {
	cert, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return OptionInjectCA(cert)
}

func OptionInjectCA(ca []byte) OptionFunc {

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	rootCAs.AppendCertsFromPEM(ca)

	return func(so *ClientOptions) {
		tslConf := &tls.Config{
			RootCAs: rootCAs,
		}

		so.ClientOptions.SetTLSConfig(tslConf)
	}
}

func OptionInsecureCerts() OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
}

func OptionSettingsWS(wrBuf, rdBuf int) OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.WebsocketOptions = &mqtt.WebsocketOptions{
			ReadBufferSize:  rdBuf,
			WriteBufferSize: wrBuf,
		}
	}
}

func OptionClientIDString(id string) OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.SetClientID(id)
	}
}

func OptionClientIDfromUUID(uuid uuid.UUID) OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.SetClientID(uuid.String())
	}
}

func OptionClientIDrandomly(size int) OptionFunc {
	return func(so *ClientOptions) {

		buf := make([]byte, size)
		if _, err := rand.Read(buf); err != nil {
			return
		}

		so.ClientOptions.SetClientID(base64.StdEncoding.EncodeToString(buf))
	}
}

func OptionSetupHeaders(h http.Header) OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.SetHTTPHeaders(h)
	}
}

type QoSValue byte

const (
	_ QoSValue = iota
	LowQoS
	BaseQoS
	HightQoS
)

func OptionDefaultQoS(QoS QoSValue) OptionFunc {
	return func(so *ClientOptions) {
		so.QoS = byte(QoS)
	}
}

func OptionReconnecting() OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.SetConnectRetry(true)
	}
}

func OptionOnLost(f func(mqtt.Client, error)) OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.OnConnectionLost = f
	}
}

func OptionOnReconn(f func(mqtt.Client, *mqtt.ClientOptions)) OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.OnReconnecting = f
	}
}

func OptionConnTimeout(t time.Duration) OptionFunc {
	return func(so *ClientOptions) {
		so.ClientOptions.ConnectTimeout = t
	}
}
