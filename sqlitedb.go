package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

type SqliteDB struct {
	db *gorm.DB
}

const (
	dialect = "sqlite3"
	dbfile = "internshiptask.db"
)

func createSQLDB() (SqliteDB, error) {
	db, err := gorm.Open(dialect, dbfile)
	if err != nil {
		return SqliteDB{}, err
	}
	return SqliteDB{
		db: db,
	}, nil
}

func (db SqliteDB) Get() ([]Message, error) {
	var messages []Message
	db.db.Find(&messages)
	return messages, nil
}

func (db SqliteDB) Create(message Message) (Message, error) {
	message.Timestamp = time.Now().Unix()
	db.db.Create(&message)
	return message, nil
}

func (Message) TableName() string {
	return "messages"
}
