package exporters

import (
	"github.com/eterline/desky-backend/internal/models"
)

type Repository interface {
	Get() ([]models.ExporterInfoT, error)
	GetType(models.ExporterTypeString) ([]models.ExporterInfoT, error)
	Add(obj *models.ExporterInfoT) error
	Edit(obj *models.ExporterInfoT, id uint) error
	Delete(id uint) error
}

type ExporterService struct {
	repo Repository
}

func New(r Repository) *ExporterService {
	return &ExporterService{
		repo: r,
	}
}

func (es *ExporterService) Services(exporter models.ExporterTypeString) ([]models.ExporterInfo, error) {

	infoTableList := make([]models.ExporterInfoT, 0)

	switch exporter {

	case models.ExporterDockerType:
		if err := es.getTypeExporter(&infoTableList, exporter); err != nil {
			return nil, err
		}

	case models.ExporterProxmoxType:
		if err := es.getTypeExporter(&infoTableList, exporter); err != nil {
			return nil, err
		}

	default:
		infoTs, err := es.repo.Get()
		if err != nil {
			return nil, err
		}

		infoTableList = append(infoTableList, infoTs...)
	}

	infoList := make([]models.ExporterInfo, len(infoTableList))

	for idx, info := range infoTableList {

		infoList[idx] = models.ExporterInfo{
			ID:    info.ID,
			Type:  info.ResolveType(),
			API:   info.API,
			Extra: info.ResolveExtra(),
		}

	}

	return infoList, nil
}

func (es *ExporterService) Append(form models.ExporterForm) error {

	info := &models.ExporterInfoT{
		Type:  string(form.ValueType()),
		API:   form.ValueAPI(),
		Extra: form.ValueExtra(),
	}

	if err := es.repo.Add(info); err != nil {
		return err
	}

	return nil
}

func (es *ExporterService) Delete(id int) error {

	if err := es.repo.Delete(uint(id)); err != nil {
		return err
	}

	return nil
}

func (es *ExporterService) getTypeExporter(l *[]models.ExporterInfoT, exporter models.ExporterTypeString) error {
	infoTs, err := es.repo.GetType(exporter)
	if err != nil {
		return err
	}

	*l = infoTs
	return nil
}
