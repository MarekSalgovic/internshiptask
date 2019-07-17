package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

type SqliteDB struct {
	db *gorm.DB
}

const (
	DIALECT = "sqlite3"
	DBFILE  = "internshiptask.db"
)

func createSQLDB(dialect string, dbfile string) (SqliteDB, error) {
	db, err := gorm.Open(dialect, dbfile)
	if err != nil {
		return SqliteDB{}, err
	}
	return SqliteDB{
		db: db,
	}, nil
}

func (db *SqliteDB) Get() ([]Log, error) {
	var logs []Log
	db.db.Find(&logs)
	return logs, nil
}

func (db *SqliteDB) GetByUser(id string) ([]Log, error) {
	var logs []Log
	db.db.Where("usrid = ?", id).Find(&logs)
	return logs, nil
}

func (db *SqliteDB) Create(log Log) (Log, error) {
	log.Timestamp = time.Now().Unix()
	unique := generateRandomString(10)
	for db.db.Where("Unique_phrase = ?", unique).Error == gorm.ErrRecordNotFound {
		unique = generateRandomString(10)
	}
	log.Unique_phrase = unique
	db.db.Create(&log)
	return log, nil
}

func (db *SqliteDB) Update(id int, m Log) (Log, error) {
	var log Log
	err := db.db.First(&log, id).Error
	if err != nil {
		return Log{}, err
	}
	unique := log.Unique_phrase
	log.Id = m.Id
	log.Log = m.Log
	log.Notification_email = m.Notification_email
	log.Notification_email_optional = m.Notification_email_optional
	log.Timestamp = time.Now().Unix()
	log.Unique_phrase = unique
	db.db.Save(&log)
	return log, nil
}

func (Log) TableName() string {
	return "logs"
}
