package agentclient

import (
	"net/url"
	"strings"

	"github.com/eterline/desky-backend/pkg/agent-client/requester"
)

// Reg - register new agent session
func Reg(api, key string) (*DeskyAgent, error) {
	url, err := url.Parse(api)
	if err != nil {
		return nil, err
	}

	info, ok := info(url, keyBerearer(key))
	if !ok {
		return nil, ErrInvalidAgent
	}

	return &DeskyAgent{
		Url:   url,
		Info:  info,
		token: keyBerearer(key),
	}, nil
}

// Parameter - get parameter by value name:
// "host" |
// "cpu" |
// "ram" |
// "load" |
// "temperature" |
// "ports" |
// "partitions"
func (a *DeskyAgent) Parameter(exporter string) (any, error) {

	data, path := a.modelFabric(exporter)
	if data == nil {
		return nil, ErrExporterNotExists
	}

	req, err := requester.Make(a.apiAppend(path), &requester.RequestOptions{
		SSLVerify: false,
		Headers: map[string]string{
			"Authorization": a.token.Berearer(),
		},
	})

	if err != nil {
		return nil, err
	}

	if _, err := req.GET(); err != nil {
		return nil, err
	}

	if err := req.Resolve(data); err != nil {
		return nil, err
	}

	return data, nil
}

func info(api *url.URL, token keyBerearer) (data *HostInfo, ok bool) {

	defer func() {
		if r := recover(); r == nil {
			return
		}
	}()

	req, err := requester.Make(
		strings.ReplaceAll(api.String(), "api", "info"),
		&requester.RequestOptions{
			SSLVerify: false,
			Headers: map[string]string{
				"Authorization": token.Berearer(),
			},
		},
	)

	if err != nil {
		return
	}

	if _, err := req.GET(); err != nil {
		return
	}

	data = new(HostInfo)

	if err := req.Resolve(data); err != nil {
		return nil, false
	}

	ok = true
	return
}

func (a *DeskyAgent) modelFabric(exporter string) (data any, url string) {

	lower := strings.ToLower(exporter)
	url = "/stats/" + lower

	switch lower {

	case "host":
		data = new(Host)

	case "cpu":
		data = new(CPU)

	case "ram":
		data = new(RAM)

	case "load":
		data = new(Load)

	case "temperature":
		data = new(SensorList)

	case "ports":
		data = new(Ports)

	case "partitions":
		data = new(PartitionList)

	default:
		data = nil
	}

	return
}

func (a *DeskyAgent) apiAppend(path string) string {
	return a.Url.String() + path
}
