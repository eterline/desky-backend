package ve

import "fmt"

func uptimeStr(time uint64) string {
	days := time / 86400
	hours := time % 86400 / 3600
	mins := time % 86400 % 3600 / 60
	sec := time % 86400 % 3600 % 60

	return fmt.Sprintf("%vd|%vh|%vm|%vs", days, hours, mins, sec)
}

func sizeStrMB(data uint64) string {
	mbytes := data / 1024 / 1024

	return fmt.Sprintf("%vMB", mbytes)
}
