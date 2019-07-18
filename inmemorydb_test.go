package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestSuiteInMemory struct {
	suite.Suite
	data []Log
	db   InMemoryDB
}

func (suite *TestSuiteInMemory) SetupSuite() {
	db := createInMemoryDB()
	suite.db = db

}

func (suite *TestSuiteInMemory) SetupTest() {
	suite.data = []Log{
		{
			Id:                          "user1",
			Message:                     "abcde",
			Timestamp:                   time.Now().Unix(),
			Notification_email:          "user1@xxx.cz",
			Notification_email_optional: "",
		}, {
			Id:                          "user1",
			Message:                     "xyz",
			Timestamp:                   time.Now().Unix(),
			Notification_email:          "user1@xxx.cz",
			Notification_email_optional: "",
		}, {
			Id:                          "user2",
			Message:                     "ABCDEF",
			Timestamp:                   time.Now().Unix(),
			Notification_email:          "user2@xxx.sk",
			Notification_email_optional: "",
		}, {
			Id:                          "user1",
			Message:                     "123456",
			Timestamp:                   time.Now().Unix(),
			Notification_email:          "user1@xxx.cz",
			Notification_email_optional: "",
		},
	}
}

func (suite *TestSuiteInMemory) TearDownTest() {
	for k := range suite.db.data {
		delete(suite.db.data, k)
	}
	suite.db.id = 0
}

func (suite *TestSuiteInMemory) TearDownSuite() {
}

func TestSuiteInMemoryFunc(t *testing.T) {
	suite.Run(t, new(TestSuiteInMemory))
}

func (suite *TestSuiteInMemory) TestSqliteDB_Create1() {
	test := CreateLogTableTest{
		Name:   "create log 1",
		Input:  suite.data[0],
		Output: suite.data[0],
	}
	log, err := suite.db.Create(test.Input)
	suite.Nil(err)
	suite.True(testCompareLogs(log, test.Output), "get create1")
}

func (suite *TestSuiteInMemory) TestSqliteDB_Create2() {
	test := CreateLogTableTest{
		Name:   "create log 2",
		Input:  suite.data[1],
		Output: suite.data[1],
	}
	log, err := suite.db.Create(test.Input)
	suite.Nil(err)
	suite.True(testCompareLogs(log, test.Output), "get create2")
}

func (suite *TestSuiteInMemory) TestSqliteDB_Get1() {
	test := GetLogsTableTest{
		Name:   "empty",
		Input:  0,
		Output: 0,
	}
	for i := 0; i < test.Input; i++ {
		c := i % len(suite.data)
		log := suite.data[c]
		_, err := suite.db.Create(log)
		suite.Nil(err)
	}
	logs, err := suite.db.Get()
	suite.Nil(err)
	suite.Equal(test.Output, len(logs), "get test1")
}

func (suite *TestSuiteInMemory) TestSqliteDB_Get2() {
	test := GetLogsTableTest{
		Name:   "empty",
		Input:  3,
		Output: 3,
	}
	for i := 0; i < test.Input; i++ {
		c := i % len(suite.data)
		log := suite.data[c]
		_, err := suite.db.Create(log)
		suite.Nil(err)
	}
	logs, err := suite.db.Get()
	suite.Nil(err)
	suite.Equal(test.Output, len(logs), "get test2")
}

func (suite *TestSuiteInMemory) TestSqliteDB_Get3() {
	test := GetLogsTableTest{
		Name:   "seven records",
		Input:  7,
		Output: 7,
	}
	for i := 0; i < test.Input; i++ {
		c := i % len(suite.data)
		log := suite.data[c]
		_, err := suite.db.Create(log)
		suite.Nil(err)
	}
	logs, err := suite.db.Get()
	suite.Nil(err)
	suite.Equal(test.Output, len(logs), "get test3")
}

func (suite *TestSuiteInMemory) TestInSqliteDB_GetByUser1() {
	test := GetByUserTableTest{
		Name:   "user1 3 logs",
		Input:  "user1",
		Output: 3,
	}
	for i := 0; i < len(suite.data); i++ {
		_, err := suite.db.Create(suite.data[i])
		suite.Nil(err)
	}
	res, err := suite.db.GetByUser(test.Input)
	suite.Nil(err)
	suite.Equal(test.Output, len(res), "get by id test1")
}

func (suite *TestSuiteInMemory) TestInSqliteDB_GetByUser2() {
	test := GetByUserTableTest{
		Name:   "user2 1 log",
		Input:  "user2",
		Output: 1,
	}
	for i := 0; i < len(suite.data); i++ {
		_, err := suite.db.Create(suite.data[i])
		suite.Nil(err)
	}
	res, err := suite.db.GetByUser(test.Input)
	suite.Nil(err)
	suite.Equal(test.Output, len(res), "get by id test2")
}

func (suite *TestSuiteInMemory) TestInSqliteDB_GetByUser3() {
	test := GetByUserTableTest{
		Name:   "no id",
		Input:  "",
		Output: 0,
	}
	for i := 0; i < len(suite.data); i++ {
		_, err := suite.db.Create(suite.data[i])
		suite.Nil(err)
	}
	res, err := suite.db.GetByUser(test.Input)
	suite.Nil(err)
	suite.Equal(test.Output, len(res), "get by id test3")
}

func (suite *TestSuiteInMemory) TestInSqliteDB_Update1() {
	test := UpdateLogsTableTest{
		Name:    "update id 1",
		InputID: 1,
		InputLog: Log{
			Id:                          "user2",
			Message:                     "UPDATE1",
			Notification_email:          "user2@abcd.cz",
			Notification_email_optional: "",
		},
		Output: Log{
			Id:                          "user2",
			Message:                     "UPDATE1",
			Notification_email:          "user2@abcd.cz",
			Notification_email_optional: "",
		},
		OutputError: false,
	}
	for i := 0; i < len(suite.data); i++ {
		_, err := suite.db.Create(suite.data[i])
		suite.Nil(err)
	}
	res, err := suite.db.Update(test.InputID, test.InputLog)
	suite.Equal(test.OutputError, (err != nil))
	suite.Equal(true, testCompareLogs(res, test.Output))
}

func (suite *TestSuiteInMemory) TestInSqliteDB_Update2() {
	test := UpdateLogsTableTest{
		Name:    "update id 3",
		InputID: 3,
		InputLog: Log{
			Id:                          "user1",
			Message:                     "UPDATE_TEST2",
			Notification_email:          "user1@abcd.cz",
			Notification_email_optional: "user1@xxx.sk",
		},
		Output: Log{
			Id:                          "user1",
			Message:                     "UPDATE_TEST2",
			Notification_email:          "user1@abcd.cz",
			Notification_email_optional: "user1@xxx.sk",
		},
		OutputError: false,
	}
	for i := 0; i < len(suite.data); i++ {
		_, err := suite.db.Create(suite.data[i])
		suite.Nil(err)
	}
	res, err := suite.db.Update(test.InputID, test.InputLog)
	suite.Equal(test.OutputError, (err != nil))
	suite.Equal(true, testCompareLogs(res, test.Output))
}

func (suite *TestSuiteInMemory) TestInSqliteDB_Update3() {
	test := UpdateLogsTableTest{
		Name:    "log with id not found",
		InputID: 10,
		InputLog: Log{
			Id:                          "user2",
			Message:                     "UPDATE3",
			Notification_email:          "user2@abcd.cz",
			Notification_email_optional: "",
		},
		Output:      Log{},
		OutputError: true,
	}
	for i := 0; i < len(suite.data); i++ {
		_, err := suite.db.Create(suite.data[i])
		suite.Nil(err)
	}
	res, err := suite.db.Update(test.InputID, test.InputLog)
	suite.Equal(test.OutputError, (err != nil))
	suite.Equal(true, testCompareLogs(res, test.Output))
}
