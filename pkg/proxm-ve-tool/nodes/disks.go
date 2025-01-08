package nodes

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

const (
	ListPath  = "/disks/list"
	SmartPath = "/disks/smart?disk=%s"
)

func (pn *ProxmoxNode) Disks(ctx context.Context) (dsk *DisksInfo, err error) {
	req := pn.session.MakeRequest(ctx, pn.urlWithName(ListPath))
	defer req.EndTask()
	dsk = &DisksInfo{}

	dsk.Code, err = req.GET()
	if err != nil {
		return nil, err
	}

	if err := req.Resolve(&dsk); err != nil {
		return nil, err
	}

	return dsk, nil
}

func (pn *ProxmoxNode) DiskByDevPath(ctx context.Context, path string) (dsk *DiskUnit, err error) {

	lst, err := pn.Disks(ctx)
	if err != nil {
		return nil, err
	}

	for _, d := range lst.Data {
		if d.Devpath == path {
			return &d, nil
		}
	}

	return nil, ErrDiskPathNotExists(path)
}

func (d *DiskUnit) SMART(ctx context.Context) (out Smart, err error) {
	req := d.pn.session.MakeRequest(ctx, d.pn.urlWithName(fmt.Sprintf(SmartPath, d.Devpath)))
	defer req.EndTask()

	smart := DiskSmartData{}

	smart.Code, err = req.GET()
	if err != nil {
		return Smart{}, err
	}

	if err := req.Resolve(&smart); err != nil {
		return Smart{}, err
	}

	out = Smart{
		Wearout: smart.Data.Wearout,
		Health:  smart.Data.Health,
		Type:    smart.Data.Type,
		Info:    infoTabler(smart.Data.Text),
	}

	return out, nil
}

func (d *DiskUnit) SerialInt() int {
	ser, _ := strconv.Atoi(d.Serial)
	return ser
}

func infoTabler(info string) SmartInfo {

	strs := strings.Split(info, "\n")
	table := make(SmartInfo, len(strs))

	for _, str := range strs {

		param := strings.Split(str, ":")

		if len(param) > 1 {
			table[strings.ReplaceAll(param[0], " ", "")] = strings.ReplaceAll(param[1], " ", "")
		}

	}

	return table
}
