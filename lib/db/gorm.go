package db

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	MysqlType    = 1
	PostgresType = 2

	TxContext = "db tx context"
)

// DbConfig ...
type DbConfig struct {
	DbType    int
	DbName    string
	Host      string
	Port      string
	Username  string
	Pwd       string
	Charset   string
	ParseTime bool
}

type WarpDB struct {
	*gorm.DB
}

func (db WarpDB) WithContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TxContext).(*gorm.DB); ok {
		return tx
	}

	return db.DB

}

// dbEngine global def
var DBEngine WarpDB

// SetupDBEngine init call
func SetupDBEngine(cfg DbConfig) {
	switch cfg.DbType {
	case MysqlType:
		dbEngine, err := NewMysqlEngine(cfg)
		if err != nil {
			log.Panic("NewDBEngine error : ", err)
		}
		DBEngine = WarpDB{dbEngine}
	case PostgresType:
		dbEngine, err := NewPostgresEngine(cfg)
		if err != nil {
			log.Panic("NewDBEngine error : ", err)
		}
		DBEngine = WarpDB{dbEngine}
	}

}

func NewPostgresEngine(dbConfig DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Shanghai",
		dbConfig.Host, dbConfig.Username, dbConfig.Pwd, dbConfig.DbName, dbConfig.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewMysqlEngine(dbConfig DbConfig) (*gorm.DB, error) {
	conn := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	dsn := fmt.Sprintf(conn, dbConfig.Username, dbConfig.Pwd, dbConfig.Host, dbConfig.DbName, dbConfig.Charset, dbConfig.ParseTime)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	return db, nil
}
