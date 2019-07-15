package main

import "time"

type InMemoryDB struct {
	data map[int]Message
	id   int
}

func createInMemoryDB() InMemoryDB {
	return InMemoryDB{
		data: make(map[int]Message),
		id:   0,
	}
}

func (db *InMemoryDB) Get() ([]Message, error) {
	values := make([]Message, 0, len(db.data))
	for _, v := range db.data {
		values = append(values, v)
	}
	return values, nil
}

func (db *InMemoryDB) Create(message Message) (Message, error) {
	message.Timestamp = time.Now().Unix()
	db.data[db.id] = message
	db.id++
	return message, nil
}
