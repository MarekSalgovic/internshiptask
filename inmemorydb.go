package main

import "time"

type InMemoryDB struct {
	data map[int]Log
	id   int
}

func createInMemoryDB() InMemoryDB {
	return InMemoryDB{
		data: make(map[int]Log),
		id:   0,
	}
}

func (db *InMemoryDB) Get() ([]Log, error) {
	values := make([]Log, 0, len(db.data))
	for _, v := range db.data {
		values = append(values, v)
	}
	return values, nil
}

func (db *InMemoryDB) GetByUser(id string) ([]Log, error) {
	values := make([]Log, 0, len(db.data))
	for _, v := range db.data {
		if v.Id == id {
			values = append(values, v)
		}
	}
	return values, nil
}

func (db *InMemoryDB) Create(log Log) (Log, error) {
	log.Timestamp = time.Now().Unix()
	db.data[db.id] = log
	db.id++
	unique := generateRandomString(10)
	for i := 0; i < len(db.data); i++ {
		if db.data[i].Unique_phrase == unique {
			unique = generateRandomString(10)
			i = -1
		}
	}
	return log, nil
}

func (db *InMemoryDB) Update(id int, log Log) (Log, error) {
	db.data[id] = log
	return log, nil
}
