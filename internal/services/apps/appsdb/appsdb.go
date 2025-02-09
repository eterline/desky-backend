package appsdb

import (
	"sync"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/repository"
	"github.com/eterline/desky-backend/internal/utils"
)

type AppsService struct {
	repository *repository.AppsRepository
	sync.Mutex
}

func New(repository *repository.AppsRepository) *AppsService {
	return &AppsService{
		repository: repository,
	}
}

func (sc *AppsService) Append(topic string, app models.AppDetails) error {

	sc.Lock()
	defer sc.Unlock()

	instance := &models.AppsInstancesT{
		Name:        app.Name,
		Icon:        app.Icon,
		Description: app.Description,
		Link:        app.Link,

		Topic: models.AppsTopicT{
			Name: topic,
		},
	}

	if err := sc.repository.CreateApp(instance); err != nil {
		return err
	}

	return nil
}

func (sc *AppsService) Delete(topic string, topicQuery int) error {

	sc.Lock()
	defer sc.Unlock()

	apps, err := sc.getTable()
	if err != nil {
		return err
	}

	t, ok := apps[topic]
	if !ok || !(len(t) > 0) {
		return ErrAppNotFound
	}

	var id *uint = nil

	for q, app := range t {
		if q == topicQuery {
			id = &app.ID
			break
		}
	}

	if id == nil {
		return ErrAppNotFound
	}

	if err := sc.repository.DeleteApp(*id); err != nil {
		return err
	}

	apps[topic] = utils.RemoveSliceIndex(apps[topic], topicQuery)

	if len(apps[topic]) > 0 {
		return nil
	}

	if err := sc.repository.DeleteTopic(topic); err != nil {
		return err
	}

	return nil
}

func (sc *AppsService) Edit(topic string, topicQuery int, app models.AppDetails) error {

	sc.Lock()
	defer sc.Unlock()

	return nil
}

func (sc *AppsService) Table() (models.AppsTable, error) {

	sc.Lock()
	defer sc.Unlock()

	return sc.getTable()
}

func (sc *AppsService) getTable() (models.AppsTable, error) {

	apps, err := sc.repository.Table()
	if err != nil {
		return nil, err
	}

	table := models.AppsTable{}

	if len(apps) > 0 {

		for _, app := range apps {

			table[app.Topic.Name] = append(
				table[app.Topic.Name],
				models.AppDetails{
					ID:          app.ID,
					Name:        app.Name,
					Description: app.Description,
					Icon:        app.Icon,
					Link:        app.Link,
				},
			)
		}

	}

	return table, nil
}
