package sshlander

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/router/handler"
)

type SSHRepository interface {
	AddHost(username string, host string, port uint16, osType string, privateKeyUsage bool, password string, key string) error
	Delete(id int) error
	QueryAll() ([]models.SSHCredentialsT, error)
}

type SSHLanderControllers struct {
	ctx     context.Context
	repoSSH SSHRepository
}

func Init(
	ctx context.Context,
	repo SSHRepository,
) *SSHLanderControllers {
	return &SSHLanderControllers{
		ctx:     ctx,
		repoSSH: repo,
	}
}

func (mc *SSHLanderControllers) ListHosts(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sshlander.list-hosts"

	sshList, err := mc.repoSSH.QueryAll()
	if err != nil {
		return op, err
	}

	if sshList == nil || len(sshList) < 1 {
		w.WriteHeader(http.StatusNoContent)
		return op, nil
	}

	resultList := make([]models.SSHInstanceObject, len(sshList))

	for idx, sshInst := range sshList {

		hostString := fmt.Sprintf(
			"%s@%s:%v",
			sshInst.Username,
			sshInst.Host,
			sshInst.Port,
		)

		resultList[idx] = models.SSHInstanceObject{
			ID:            int(sshInst.ID),
			HostString:    hostString,
			PrivateKeyUse: sshInst.Security.PrivateKeyUse,
		}
	}

	return op, handler.WriteJSON(w, http.StatusOK, resultList)
}

func (mc *SSHLanderControllers) AppendHost(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sshlander.append-host"

	form := new(models.RequestFormSSH)
	if err := handler.DecodeRequest(r, form); err != nil {
		return op, err
	}

	if err := handler.Validate(form); err != nil {
		return op, err
	}

	if err := mc.repoSSH.AddHost(
		form.User, form.Host, form.Port,
		form.System,
		form.PrivateKeyUse,
		form.Password, form.PrivateKey,
	); err != nil {
		return op, err
	}

	response := models.ResponseCreateSSH{
		PrivateKeyUse: form.PrivateKeyUse,
		Target:        fmt.Sprintf("%s@%s:%v", form.User, form.Host, form.Port),
	}

	return op, handler.WriteJSON(w, http.StatusCreated, response)
}

func (mc *SSHLanderControllers) DeleteHost(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sshlander.delete-host"

	q, err := handler.ParseURLParameters(r, handler.NumOpts("id"))
	if err != nil {
		return op, err
	}

	if err := mc.repoSSH.Delete(q.GetInt("id")); err != nil {
		return op, err
	}

	return op, handler.StatusOK(w, "host deleted")
}
