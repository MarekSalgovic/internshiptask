package main

import (
	"context"
	"fmt"
	pb "github.com/MarekSalgovic/internshiptask/grpc"
	"google.golang.org/grpc"
)

const (
	address = "localhost:3000"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewLoggerClient(conn)
	r, err := c.GetLogs(context.Background(), &pb.GetRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("<----------------------Empty DB---------------------->")
	for i := 0; i < len(r.Logs); i++ {
		fmt.Println(r.Logs[i])
		printline()
	}
	for i := 0; i < 7; i++ {
		user := fmt.Sprintf("user%d", i%3)
		sprava := fmt.Sprintf("sprava cislo #%d", i)
		_, err = c.PostLog(context.Background(), &pb.PostRequest{
			Log: &pb.NewLog{
				Usrid:                     user,
				Log:                       sprava,
				NotificationEmail:         user + "@xxxx.sk",
				NotificationEmailOptional: user + "@mail.cz",
			},
		})
	}

	r, err = c.GetLogs(context.Background(), &pb.GetRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("<----------------------7 POST---------------------->")
	for i := 0; i < len(r.Logs); i++ {
		fmt.Println(r.Logs[i])
		printline()
	}

	r, err = c.GetLogsByUser(context.Background(), &pb.GetByUserRequest{Usrid: "user1"})
	if err != nil {
		panic(err)
	}
	fmt.Println("<----------------------GET BY ID: USER1---------------------->")
	for i := 0; i < len(r.Logs); i++ {
		fmt.Println(r.Logs[i])
		printline()
	}
	_, err = c.UpdateLog(context.Background(), &pb.PutRequest{
		LogId: 3,
		Log: &pb.Log{
			Id:                        3,
			Usrid:                     "user10",
			Log:                       "zmenena sprava 3",
			NotificationEmail:         "user10@abc.sk",
			NotificationEmailOptional: "user10@xxx.cz"}})
	if err != nil {
		panic(err)
	}
	fmt.Println("<----------------------UPDATE MESSAGE ID: 3---------------------->")
	r, err = c.GetLogs(context.Background(), &pb.GetRequest{})
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(r.Logs); i++ {
		fmt.Println(r.Logs[i])
		printline()
	}
}

func printline() {
	fmt.Println("----------------------------------------------")
}
