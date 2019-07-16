package main

import (
	"time"
)

type CreateLogTableTest struct {
	Name   string
	Input  Log
	Output Log
}

type GetLogsTableTest struct {
	Name   string
	Input  int
	Output int
}

type GetByUserTableTest struct {
	Name   string
	Input  string
	Output int
}

type UpdateLogsTableTest struct {
	Name        string
	InputID     int
	InputLog    Log
	Output      Log
	OutputError bool
}

var testdata = []Log{
	{
		Id:                          "user1",
		Log:                         "abcde",
		Timestamp:                   time.Now().Unix(),
		Notification_email:          "user1@xxx.cz",
		Notification_email_optional: "",
	}, {
		Id:                          "user1",
		Log:                         "xyz",
		Timestamp:                   time.Now().Unix(),
		Notification_email:          "user1@xxx.cz",
		Notification_email_optional: "",
	}, {
		Id:                          "user2",
		Log:                         "ABCDEF",
		Timestamp:                   time.Now().Unix(),
		Notification_email:          "user2@xxx.sk",
		Notification_email_optional: "",
	}, {
		Id:                          "user1",
		Log:                         "123456",
		Timestamp:                   time.Now().Unix(),
		Notification_email:          "user1@xxx.cz",
		Notification_email_optional: "",
	},
}

func testCompareLogs(l1 Log, l2 Log) bool {
	if l1.Id == l2.Id && l1.Log == l2.Log && l1.Notification_email == l2.Notification_email &&
		l1.Notification_email_optional == l2.Notification_email_optional {
		return true
	}
	return false
}
