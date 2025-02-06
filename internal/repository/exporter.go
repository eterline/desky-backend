package repository

import (
	"github.com/eterline/desky-backend/internal/models"
	"gorm.io/gorm"
)

type ExporterRepository struct {
	DefaultRepository
}

func NewExporterRepository(db *gorm.DB) *ExporterRepository {
	return &ExporterRepository{
		DefaultRepository{
			db: db,
		},
	}
}

func (r *ExporterRepository) Get() ([]models.ExporterInfoT, error) {
	list := make([]models.ExporterInfoT, 0)

	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *ExporterRepository) GetType(exporterType models.ExporterTypeString) ([]models.ExporterInfoT, error) {

	list := make([]models.ExporterInfoT, 0)

	if err := r.db.Find(&list, "Type = ?", exporterType).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *ExporterRepository) Add(exporter *models.ExporterInfoT) error {
	return r.db.Create(exporter).Error
}

func (r *ExporterRepository) Edit(exporter *models.ExporterInfoT, id uint) error {

	if err := r.db.First(exporter, "ID = ?", id).Error; err != nil {
		return err
	}

	exporter.ID = id

	if err := r.db.Save(exporter).Error; err != nil {
		return err
	}

	return nil
}

func (r *ExporterRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(new(models.ExporterInfoT), "ID = ?", id).Error
}

// func (r *ProxmoxRepository) Create(credentials *models.ProxmoxCredentialsT) error {
// 	return r.db.Create(credentials).Error
// }
