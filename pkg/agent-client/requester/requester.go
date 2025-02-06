package requester

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type RequestOptions struct {
	SSLVerify bool
	Headers   map[string]string
	Timeout   time.Duration
}

type RequestProvide struct {
	request  *http.Request
	client   http.Client
	response *http.Response
}

func Make(url string, opts *RequestOptions) (*RequestProvide, error) {

	ctx, _ := context.WithTimeout(
		context.Background(),
		func() time.Duration {
			if opts != nil && opts.Timeout != 0 {
				return opts.Timeout
			}
			return 5 * time.Second
		}(),
	)

	r, err := http.NewRequestWithContext(ctx, "", url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	if opts != nil && len(opts.Headers) > 0 {
		for key, header := range opts.Headers {
			r.Header.Add(key, header)
		}
	}

	return &RequestProvide{
		client: http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !opts.SSLVerify,
			},
		}},
		request:  r,
		response: nil,
	}, nil
}

func (rp *RequestProvide) GET() (code int, err error) {
	rp.request.Method = "GET"

	response, err := rp.fetch(rp.request)
	if err != nil {
		return 0, err
	}

	rp.response = response

	return response.StatusCode, ValidateOKCodes(response.StatusCode)
}

func (rp *RequestProvide) POST(body io.ReadCloser) (code int, err error) {
	rp.request.Method = "POST"

	rp.request.Body = body
	defer body.Close()

	response, err := rp.fetch(rp.request)
	if err != nil {
		return 0, err
	}

	rp.response = response

	return response.StatusCode, ValidateOKCodes(response.StatusCode)
}

func (rp *RequestProvide) fetch(r *http.Request) (*http.Response, error) {
	response, err := rp.client.Do(r)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, ErrResponseNotImplemented
	}

	return response, nil
}

func (rp *RequestProvide) Resolve(v any) error {
	defer rp.response.Body.Close()
	return json.NewDecoder(rp.response.Body).Decode(v)
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
