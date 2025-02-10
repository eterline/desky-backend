package storage

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type StorageProvider interface {
	Socket() gorm.Dialector
	StorageType() string
	Source() string
}

// ====================================================================

func NewStorageSQLite(file string) StorageProvider {
	return &StorageSQLite{
		filename: file,
	}
}

type StorageSQLite struct {
	filename string
}

func (s *StorageSQLite) Socket() gorm.Dialector {
	return sqlite.Open(s.filename)
}

func (s *StorageSQLite) StorageType() string {
	return "sqlite"
}

func (s *StorageSQLite) Source() string {
	return s.filename
}

// ====================================================================

func NewStoragePostgres(
	username, password, host string,
	port uint16,
	dbname string,
	sslmode bool,
	tzRegion, tzCity string,
) StorageProvider {

	timeZone := fmt.Sprintf("%s/%s", tzRegion, tzCity)

	return &StoragePostgres{
		username: username,
		password: password,
		host:     host,
		port:     port,
		dbname:   dbname,
		sslmode:  sslmode,
		tz:       timeZone,
	}
}

type StoragePostgres struct {
	username string
	password string
	host     string
	port     uint16
	dbname   string
	sslmode  bool
	tz       string
}

func (s *StoragePostgres) Socket() gorm.Dialector {
	var dsnList [7]string

	dsnList[0] = fmt.Sprintf("host=%s", s.host)
	dsnList[1] = fmt.Sprintf("user=%s", s.username)
	dsnList[2] = fmt.Sprintf("password=%s", s.password)

	dsnList[3] = fmt.Sprintf("dbname=%s", s.dbname)
	dsnList[4] = fmt.Sprintf("port=%v", s.port)
	dsnList[5] = fmt.Sprintf("TimeZone=%s", s.tz)

	if s.sslmode {
		dsnList[6] = fmt.Sprintf("sslmode=enable")
	} else {
		dsnList[6] = fmt.Sprintf("sslmode=disable")
	}

	return postgres.Open(strings.Join(dsnList[:], " "))
}

func (s *StoragePostgres) StorageType() string {
	return "postgres"
}

func (s *StoragePostgres) Source() string {
	return s.dbname
}
