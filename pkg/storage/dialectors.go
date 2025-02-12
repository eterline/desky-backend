package storage

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
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

type StoragePostgres struct {
	username string
	password string
	host     string
	port     uint16
	dbname   string
	sslmode  bool
	tz       string
}

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

// ====================================================================

type MySQLCharset string

const (
	CharsetUTF8MB4 MySQLCharset = "utf8mb4"
	CharsetUTF8    MySQLCharset = "utf8"
	CharsetLATIN1  MySQLCharset = "latin1"
	CharsetASCII   MySQLCharset = "ascii"
	CharsetUCS2    MySQLCharset = "ucs2"
	CharsetCP1251  MySQLCharset = "cp1251"
	CharsetGBK     MySQLCharset = "gbk"
	CharsetSJIS    MySQLCharset = "sjis"
)

type StorageMySQL struct {
	password, username string
	dbname, host       string

	port uint16

	charset   string
	parseTime bool
}

func NewStorageMySQL(
	username, password, host string,
	port uint16,
	dbname string,
	charset MySQLCharset,
	parseTime bool,
) StorageProvider {

	return &StorageMySQL{
		username:  username,
		password:  password,
		host:      host,
		port:      port,
		dbname:    dbname,
		charset:   string(charset),
		parseTime: parseTime,
	}
}

func (s *StorageMySQL) Socket() gorm.Dialector {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%v)/%s?charset=%s&parseTime=%v",
		s.username, s.password,
		s.host, s.port, s.dbname,
		s.charset, s.parseTime,
	)

	return mysql.Open(dsn)
}

func (s *StorageMySQL) StorageType() string {
	return "mysql"
}

func (s *StorageMySQL) Source() string {
	return s.dbname
}
