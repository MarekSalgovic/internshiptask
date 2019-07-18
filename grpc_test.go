package main

import (
	"context"
	"fmt"
	pb "github.com/MarekSalgovic/internshiptask/grpc"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestSuiteGrpc struct {
	suite.Suite
	server server
	data   []Log
	db     *SqliteDB
}

func TestSuiteGrpcFunc(t *testing.T) {
	suite.Run(t, new(TestSuiteGrpc))
}

func (suite *TestSuiteGrpc) SetupSuite() {
	db, err := createSQLDB(DIALECT, "testdb.db")
	if err != nil {
		panic(err)
	}
	suite.server.svc = NewService(&db)
	suite.db = &db
}

func (suite *TestSuiteGrpc) SetupTest() {
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
	suite.db.db.AutoMigrate(&Log{})

}

func (suite *TestSuiteGrpc) TearDownTest() {
	var logs []Log
	suite.db.db.DropTable(&logs)
}

func (suite *TestSuiteGrpc) TearDownSuite() {
	suite.db.db.Close()
}

func (suite *TestSuiteGrpc) TestServer_GetLogsEmpty() {
	test := pb.GetRequest{}
	res, err := suite.server.GetLogs(context.Background(), &test)
	suite.Nil(err)
	suite.Equal(len(res.Logs), 0)
}

func (suite *TestSuiteGrpc) TestServer_PostLog() {
	test := pb.PostRequest{
		Log: &pb.NewLog{
			Id:                        suite.data[0].Id,
			Message:                   suite.data[0].Message,
			NotificationEmail:         suite.data[0].Notification_email,
			NotificationEmailOptional: suite.data[0].Notification_email_optional,
		},
	}

	_, err := suite.server.PostLog(context.Background(), &test)
	suite.Nil(err)
}

func (suite *TestSuiteGrpc) TestServer_GetLogsAll() {
	var test []pb.PostRequest
	for i := 0; i < len(suite.data); i++ {
		t := pb.PostRequest{
			Log: &pb.NewLog{
				Id:                        suite.data[i].Id,
				Message:                   suite.data[i].Message,
				NotificationEmail:         suite.data[i].Notification_email,
				NotificationEmailOptional: suite.data[i].Notification_email_optional,
			},
		}
		_, err := suite.server.PostLog(context.Background(), &t)
		suite.Nil(err)
		test = append(test, t)
	}
	testget := pb.GetRequest{}
	res, err := suite.server.GetLogs(context.Background(), &testget)
	suite.Nil(err)
	suite.Equal(len(res.Logs), len(suite.data))
}

func (suite *TestSuiteGrpc) TestServer_GetByUser1() {
	var test []pb.PostRequest
	for i := 0; i < len(suite.data); i++ {
		t := pb.PostRequest{
			Log: &pb.NewLog{
				Id:                        suite.data[i].Id,
				Message:                   suite.data[i].Message,
				NotificationEmail:         suite.data[i].Notification_email,
				NotificationEmailOptional: suite.data[i].Notification_email_optional,
			},
		}
		_, err := suite.server.PostLog(context.Background(), &t)
		suite.Nil(err)
		test = append(test, t)
	}
	testgetbyuser := pb.GetByUserRequest{
		Usrid: "user1",
	}
	res, err := suite.server.GetLogsByUser(context.Background(), &testgetbyuser)
	suite.Nil(err)
	suite.Equal(len(res.Logs), 3)
}

func (suite *TestSuiteGrpc) TestServer_GetByUserEmpty() {
	var test []pb.PostRequest
	for i := 0; i < len(suite.data); i++ {
		t := pb.PostRequest{
			Log: &pb.NewLog{
				Id:                        suite.data[i].Id,
				Message:                   suite.data[i].Message,
				NotificationEmail:         suite.data[i].Notification_email,
				NotificationEmailOptional: suite.data[i].Notification_email_optional,
			},
		}
		_, err := suite.server.PostLog(context.Background(), &t)
		suite.Nil(err)
		test = append(test, t)
	}
	testgetbyuser := pb.GetByUserRequest{
		Usrid: "empty",
	}
	res, err := suite.server.GetLogsByUser(context.Background(), &testgetbyuser)
	suite.Nil(err)
	suite.Equal(len(res.Logs), 0)
}

func (suite *TestSuiteGrpc) TestServer_UpdateLog() {
	var test []pb.PostRequest
	for i := 0; i < len(suite.data); i++ {
		t := pb.PostRequest{
			Log: &pb.NewLog{
				Id:                        suite.data[i].Id,
				Message:                   suite.data[i].Message,
				NotificationEmail:         suite.data[i].Notification_email,
				NotificationEmailOptional: suite.data[i].Notification_email_optional,
			},
		}
		_, err := suite.server.PostLog(context.Background(), &t)
		suite.Nil(err)
		test = append(test, t)
	}
	testupdatelog := pb.PutRequest{
		LogId: 1,
		Log: &pb.Log{
			Usrid:                     "newuser",
			Message:                   "aaaaaaaa",
			Timestamp:                 0,
			UniquePhrase:              "",
			NotificationEmail:         "abcd@gmail.com",
			NotificationEmailOptional: "",
		},
	}
	_, err := suite.server.UpdateLog(context.Background(), &testupdatelog)
	suite.Nil(err)
	testgetbyuser := pb.GetByUserRequest{
		Usrid: "newuser",
	}
	res, err := suite.server.GetLogsByUser(context.Background(), &testgetbyuser)
	suite.Nil(err)
	fmt.Println(res.Logs)
	suite.Equal(len(res.Logs), 1)

}
