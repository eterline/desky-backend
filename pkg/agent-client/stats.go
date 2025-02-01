package agentclient

import (
	"net/url"
	"strings"

	"github.com/eterline/desky-backend/pkg/agent-client/requester"
)

// Reg - register new agent session
func Reg(api, key string) (agent *DeskyAgent, err error) {
	url, err := url.Parse(api)
	if err != nil {
		return
	}

	agent = &DeskyAgent{
		Url:   url,
		Token: keyBerearer(key),

		valid: false,
	}

	return
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

	if !a.valid {
		return nil, ErrInvalidAgent
	}

	data, path := a.modelFabric(exporter)
	if data == nil {
		return nil, ErrExporterNotExists
	}

	req, err := requester.Make(a.apiAppend(path), &requester.RequestOptions{
		SSLVerify: false,
		Headers: map[string]string{
			"Authorization": a.Token.Berearer(),
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

func (a *DeskyAgent) Info() (data *HostInfo, ok bool) {

	data = new(HostInfo)

	defer func() {
		if r := recover(); r == nil {
			a.valid = true
			return
		}
	}()

	req, err := requester.Make(
		strings.ReplaceAll(a.Url.String(), "api", "info"),
		&requester.RequestOptions{
			SSLVerify: false,
			Headers: map[string]string{
				"Authorization": a.Token.Berearer(),
			},
		},
	)

	if err != nil {
		a.valid = false
		return nil, a.valid
	}

	if _, err := req.GET(); err != nil {
		a.valid = false
		return nil, a.valid
	}

	if err := req.Resolve(data); err != nil {
		a.valid = false
		return nil, a.valid
	}

	return data, true
}

func (a *DeskyAgent) IsValid() bool {
	return a.valid
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
		data = new(Network)

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
