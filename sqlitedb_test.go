package main

import (
	"testing"
)

func TestSqliteDB_Create(t *testing.T) {
	tests := []CreateLogTableTest{
		{
			Name:   "create log 1",
			Input:  testdata[0],
			Output: testdata[0],
		},
		{
			Name:   "create log 2",
			Input:  testdata[1],
			Output: testdata[1],
		},
	}
	for _, test := range tests {
		db, err := createSQLDB(DIALECT, "testdb.db")
		if err != nil {
			t.Fatal(err)
		}
		defer db.db.Close()
		var logs []Log
		db.db.DropTable(&logs)
		db.db.AutoMigrate(&Log{})
		t.Run(test.Name, func(t *testing.T) {
			log, err := db.Create(test.Input)
			if err != nil {
				t.Fatal(err)
			}
			if !(log.Id == test.Output.Id && log.Log == test.Output.Log &&
				log.Notification_email == test.Output.Notification_email &&
				log.Notification_email_optional == test.Output.Notification_email_optional) {
				t.Fail()
			}
		})

	}
}

func TestSqliteDB_Get(t *testing.T) {
	tests := []GetLogsTableTest{
		{
			Name:   "empty",
			Input:  0,
			Output: 0,
		},
		{
			Name:   "one record",
			Input:  1,
			Output: 1,
		},
		{
			Name:   "seven records",
			Input:  7,
			Output: 7,
		},
	}

	for _, test := range tests {
		db, err := createSQLDB(DIALECT, "testdb.db")
		if err != nil {
			t.Fatal(err)
		}
		defer db.db.Close()
		var logs []Log
		db.db.DropTable(&logs)
		db.db.AutoMigrate(&Log{})

		db.db.AutoMigrate(&Log{})
		t.Run(test.Name, func(t *testing.T) {
			for i := 0; i < test.Input; i++ {
				c := i % len(tests)
				log := testdata[c]
				_, err := db.Create(log)
				if err != nil {
					t.Fatal(err)
				}
			}
			logs, err := db.Get()
			if err != nil {
				t.Fatal(err)
			}
			if len(logs) != test.Output {
				t.Fail()
			}
		})

	}
}

func TestInSqliteDB_GetByUser(t *testing.T) {
	tests := []GetByUserTableTest{
		{
			Name:   "user1 3 logs",
			Input:  "user1",
			Output: 3,
		},
		{
			Name:   "user2 1 log",
			Input:  "user2",
			Output: 1,
		},
		{
			Name:   "no id",
			Input:  "",
			Output: 0,
		},
	}

	for _, test := range tests {
		db, err := createSQLDB(DIALECT, "testdb.db")
		if err != nil {
			t.Fatal(err)
		}
		defer db.db.Close()
		var logs []Log
		db.db.DropTable(&logs)
		db.db.AutoMigrate(&Log{})

		db.db.AutoMigrate(&Log{})
		for i := 0; i < len(testdata); i++ {
			_, err := db.Create(testdata[i])
			if err != nil {
				t.Fatal(err)
			}
		}
		t.Run(test.Name, func(t *testing.T) {
			res, err := db.GetByUser(test.Input)
			if err != nil {
				t.Fatal(err)
			}
			if len(res) != test.Output {
				t.Fail()
			}
		})

	}
}

func TestInSqliteDB_Update(t *testing.T) {
	tests := []UpdateLogsTableTest{
		{
			Name:    "update id 1",
			InputID: 1,
			InputLog: Log{
				Id:                          "user2",
				Log:                         "UPDATE1",
				Notification_email:          "user2@abcd.cz",
				Notification_email_optional: "",
			},
			Output: Log{
				Id:                          "user2",
				Log:                         "UPDATE1",
				Notification_email:          "user2@abcd.cz",
				Notification_email_optional: "",
			},
			OutputError: false,
		}, {
			Name:    "update id 3",
			InputID: 3,
			InputLog: Log{
				Id:                          "user1",
				Log:                         "UPDATE_TEST2",
				Notification_email:          "user1@abcd.cz",
				Notification_email_optional: "user1@xxx.sk",
			},
			Output: Log{
				Id:                          "user1",
				Log:                         "UPDATE_TEST2",
				Notification_email:          "user1@abcd.cz",
				Notification_email_optional: "user1@xxx.sk",
			},
			OutputError: false,
		}, {
			Name:    "log with id not found",
			InputID: 10,
			InputLog: Log{
				Id:                          "user2",
				Log:                         "UPDATE3",
				Notification_email:          "user2@abcd.cz",
				Notification_email_optional: "",
			},
			Output:      Log{},
			OutputError: true,
		},
	}
	for _, test := range tests {
		db, err := createSQLDB(DIALECT, "testdb.db")
		if err != nil {
			t.Fatal(err)
		}
		defer db.db.Close()
		var logs []Log
		db.db.DropTable(&logs)
		db.db.AutoMigrate(&Log{})

		db.db.AutoMigrate(&Log{})
		for i := 0; i < len(testdata); i++ {
			_, err := db.Create(testdata[i])
			if err != nil {
				t.Fatal(err)
			}
		}
		t.Run(test.Name, func(t *testing.T) {
			res, err := db.Update(test.InputID, test.InputLog)
			if !(testCompareLogs(res, test.Output) && (!test.OutputError == (err == nil))) {
				t.Fail()
			}
		})

	}
}
