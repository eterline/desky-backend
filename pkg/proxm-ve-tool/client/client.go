package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	AccessTicketPostfix = "/access/ticket"
)

type (
	Session struct {
		Config       SessionConfig
		cookie, csrf string
	}

	SessionConfig struct {
		APIurl      *url.URL
		Credentials UserCredentials
		ssl         bool
	}

	UserCredentials struct {
		Username string
		Password string
	}
)

func InitSession(login, password, api, postfixUrl string, SSLauthority bool) *SessionConfig {
	if postfixUrl == "" {
		postfixUrl = "/api2/json"
	}

	user := UserCredentials{
		Username: fmt.Sprintf("%s@pam", login),
		Password: password,
	}

	u, _ := url.Parse(api + postfixUrl)

	return &SessionConfig{
		APIurl:      u,
		Credentials: user,
		ssl:         SSLauthority,
	}
}

func Connect(config *SessionConfig) (*Session, error) {
	cl := http.Client{}
	form := url.Values{}
	data := ProxmoxAuthResponse{}

	cl.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !config.ssl,
		},
	}

	form.Set("username", config.Credentials.Username)
	form.Set("password", config.Credentials.Password)

	response, err := cl.PostForm(config.APIurl.String()+AccessTicketPostfix, form)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	d := &Session{
		Config: *config,
		cookie: data.Data.Ticket,
		csrf:   data.Data.CSRFPreventionToken,
	}

	return d, nil
}

type RequestProvide struct {
	Req      *http.Request
	Client   http.Client
	response *http.Response
	Base     string
}

func (s *Session) MakeRequest(ctx context.Context, path string) *RequestProvide {
	req, _ := http.NewRequestWithContext(ctx, "", s.Config.APIurl.String()+path, nil)

	req.Header.Add("Cookie", fmt.Sprintf("PVEAuthCookie=%s", s.cookie))
	req.Header.Add("CSRFPreventionToken", s.csrf)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !s.Config.ssl,
		},
	}

	return &RequestProvide{
		Req:    req,
		Client: http.Client{Transport: transport},
		Base:   s.Config.APIurl.String(),
	}
}

func (s *Session) ConnGuard(loopTime string, ctx context.Context) error {
	delay, err := time.ParseDuration(loopTime)
	if err != nil {
		return err
	}

	timer := time.NewTicker(delay)

	for {
		select {

		case <-ctx.Done():
			return nil

		case <-timer.C:
			conn, err := Connect(&s.Config)
			if err != nil {
				continue
			}

			s = conn
		}
	}
}

func (rp *RequestProvide) GET() (code int, err error) {
	rp.Req.Method = "GET"

	response, err := rp.Client.Do(rp.Req)
	if err == nil {
		rp.response = response
	}

	return response.StatusCode, ValidateOKCodes(response.StatusCode)
}

func (rp *RequestProvide) POST() (code int, err error) {
	rp.Req.Method = "POST"

	response, err := rp.Client.Do(rp.Req)
	if err == nil {
		rp.response = response
	}

	return response.StatusCode, ValidateOKCodes(response.StatusCode)
}

func (rp *RequestProvide) Resolve(v any) error {
	defer rp.response.Body.Close()
	return json.NewDecoder(rp.response.Body).Decode(v)
}

func (rp *RequestProvide) EndTask() error {
	if rp != nil {
		rp.Client.CloseIdleConnections()
		return rp.response.Body.Close()
	}
	return ErrNilConnection
}

func (rp *RequestProvide) BodyString() string {
	data, _ := io.ReadAll(rp.response.Body)
	return string(data)
}

func ValidateOKCodes(code int) error {
	if 200 > code || code > 299 {
		return ErrBadStatusCode(code)
	}
	return nil
}
