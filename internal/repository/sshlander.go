package repository

import (
	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/pkg/storage"
)

type SSHRepository struct {
	DefaultRepository
}

func NewSSHLanderRepository(db *storage.DB) *SSHRepository {
	return &SSHRepository{
		NewDefaultRepository(db),
	}
}

func (r *SSHRepository) QueryAll() ([]models.SSHCredentialsT, error) {

	credentialsList := make([]models.SSHCredentialsT, 0)

	if err := r.db.Preload("OperationSystem").Preload("Security").Find(&credentialsList).Error; err != nil {
		return nil, err
	}
	return credentialsList, nil
}

func (r *SSHRepository) AddHost(
	username string,
	host string,
	port uint16,

	osType string,
	privateKeyUsage bool,
	password, key string,
) error {

	credentialsData := &models.SSHCredentialsT{
		OperationSystem: models.MakeSSHSystemTypesT(osType),

		Username: username,
		Host:     host,
		Port:     port,

		Security: models.MakeSSHSecureT(password, privateKeyUsage, key),
	}

	if err := r.db.Create(credentialsData).Error; err != nil {
		return err
	}
	return nil
}

func (r *SSHRepository) Delete(id int) error {

	if err := r.db.Unscoped().Delete(new(models.SSHCredentialsT), "ID = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *SSHRepository) QueryById(id int) (*models.SSHCredentialsT, error) {

	sshCredentials := new(models.SSHCredentialsT)

	if err := r.db.Preload("OperationSystem").Preload("Security").First(sshCredentials, "ID = ?", id).Error; err != nil {
		return nil, err
	}
	return sshCredentials, nil
}
