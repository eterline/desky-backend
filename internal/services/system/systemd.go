package system

import (
	"encoding/json"
	"fmt"

	"github.com/bitfield/script"
)

type SystemdUnit struct {
	UnitFile string `json:"unit_file"`
	Status   string `json:"state"`
	Preset   string `json:"preset"`
}

func UnitsList() ([]SystemdUnit, error) {
	data := []SystemdUnit{}

	pipe := JqOutputCmd("systemctl list-unit-files --type=service -o json", "")

	out, err := ExecOut(pipe)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, &data); err != nil {
		return nil, ErrExec(err)
	}

	return data, nil
}

func UnitInstance(service string) (unit *SystemdUnit, err error) {

	unit = &SystemdUnit{}

	list, err := UnitsList()
	if err != nil {
		return nil, err
	}

	for _, u := range list {
		if u.UnitFile == service {
			unit := u
			return &unit, nil
		}
	}

	return nil, ErrUnitNotFound
}

func (su *SystemdUnit) Start() error {
	_, err := ExecOut(script.Exec(fmt.Sprintf("systemctl start %s", su.UnitFile)))
	return err
}

func (su *SystemdUnit) Stop() error {
	_, err := ExecOut(script.Exec(fmt.Sprintf("systemctl stop %s", su.UnitFile)))
	return err
}

func (su *SystemdUnit) Restart() error {
	_, err := ExecOut(script.Exec(fmt.Sprintf("systemctl restart %s", su.UnitFile)))
	return err
}
