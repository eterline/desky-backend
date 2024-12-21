package ve

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/luthermonson/go-proxmox"
)

type ClientAPI struct {
	TokenID string
	Secret  string
	HostURL string
}

type ClientBase struct {
	User     string
	Password string
	HostURL  string
}

func InitAPI(token, secret, host string, port int) ClientAPI {
	return ClientAPI{
		TokenID: token,
		Secret:  secret,
		HostURL: host,
	}
}

func InitBase(user, pass, host string, port int) ClientBase {
	return ClientBase{
		User:     user,
		Password: pass,
		HostURL:  apiURL(host, port),
	}
}

func apiURL(host string, port int) string {
	return fmt.Sprintf("https://%s:%v/api2/json", host, port)
}

func userCreds(user, pass string) proxmox.Credentials {
	return proxmox.Credentials{
		Username: user,
		Password: pass,
	}
}

func InsecHttp(v bool) http.Client {
	return http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: v,
			},
		},
	}
}

func (c *ClientBase) LogIn() *proxmox.Client {
	HTTPClient := InsecHttp(true)
	credentials := userCreds(c.User, c.Password)

	return proxmox.NewClient(
		c.HostURL,
		proxmox.WithHTTPClient(&HTTPClient),
		proxmox.WithCredentials(&credentials),
	)
}

func (c *ClientAPI) LogIn() *proxmox.Client {
	HTTPClient := InsecHttp(true)

	return proxmox.NewClient(
		c.HostURL,
		proxmox.WithHTTPClient(&HTTPClient),
		proxmox.WithAPIToken(c.TokenID, c.Secret),
	)
}

func Authenticate(v Auth) *proxmox.Client {
	return v.LogIn()
}

type Auth interface {
	LogIn() *proxmox.Client
}
