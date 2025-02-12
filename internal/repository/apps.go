package repository

import (
	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/pkg/storage"
)

type AppsRepository struct {
	DefaultRepository
}

func NewAppsRepository(db *storage.DB) *AppsRepository {
	return &AppsRepository{
		NewDefaultRepository(db),
	}
}

func (r *AppsRepository) ListTopic() ([]models.AppsTopicT, error) {

	topics := make([]models.AppsTopicT, 0)

	if err := r.db.Find(&topics).Error; err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *AppsRepository) CreateTopic(topic *models.AppsTopicT) error {
	return r.db.Create(topic).Error
}

func (r *AppsRepository) Edit(app *models.AppsInstancesT) error {

	if err := r.db.First(new(models.AppsInstancesT), "ID = ?", app.ID).Error; err != nil {
		return err
	}

	return r.db.Model(app).Updates(app).Error
}

func (r *AppsRepository) DeleteTopic(name string) error {
	return r.db.Unscoped().Delete(new(models.AppsTopicT), "Name = ?", name).Error
}

func (r *AppsRepository) Table() ([]models.AppsInstancesT, error) {

	var apps []models.AppsInstancesT

	if err := r.db.Preload("Topic").Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *AppsRepository) CreateApp(app *models.AppsInstancesT) error {

	err := r.db.FirstOrCreate(&app.Topic, "Name = ?", app.Topic.Name).Error
	if err != nil {
		return err
	}

	return r.db.Create(app).Error
}

func (r *AppsRepository) DeleteApp(id uint) error {
	return r.db.Delete(new(models.AppsInstancesT), id).Error
}
