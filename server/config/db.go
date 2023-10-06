package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Host       string
	Username   string
	Password   string
	DBName     string
	Connection *gorm.DB
}

func NewDB() *DB {
	c := NewConfig()
	return newDB(&DB{
		Host:     c.DB.Host,
		Username: c.DB.Username,
		Password: c.DB.Password,
		DBName:   c.DB.DBName,
	})
}

func newDB(d *DB) *DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: false,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)

	sslmode := "disable"
	if os.Getenv("GIN_MODE") == "release" {
		sslmode = "require"
	}

	db, err := gorm.Open(
		postgres.Open("postgresql://"+d.Username+":"+d.Password+"@"+d.Host+":5432/"+d.DBName+"?sslmode="+sslmode),
		&gorm.Config{Logger: newLogger},
	)
	db.Logger = db.Logger.LogMode(logger.Info)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("DB接続成功")
	d.Connection = db
	return d
}

// start transaction
func (db *DB) Begin() *gorm.DB {
	return db.Connection.Begin()
}

func (db *DB) Commit() *gorm.DB {
	return db.Connection.Commit()
}

func (db *DB) Rollback() *gorm.DB {
	return db.Connection.Rollback()
}

func (db *DB) Connect() *gorm.DB {
	return db.Connection
}
