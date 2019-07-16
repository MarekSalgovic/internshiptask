package main

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

func testCompareLogs(l1 Log, l2 Log) bool {
	if l1.Id == l2.Id && l1.Log == l2.Log && l1.Notification_email == l2.Notification_email &&
		l1.Notification_email_optional == l2.Notification_email_optional {
		return true
	}
	return false
}
