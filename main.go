package main

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/MarekSalgovic/internshiptask/grpc"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"

	//"github.com/labstack/echo/middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
)

type Log struct {
	gorm.Model
	Id                          string `json:"id" gorm:"column:usrid"`
	Log                         string `json:"log"`
	Timestamp                   int64  `json:"timestamp"`
	Unique_phrase               string `json:"unique_phrase"`
	Notification_email          string `json:"notification_email" valid:"email"`
	Notification_email_optional string `json:"notification_email_optional" valid:"email, optional"`
}

type user struct {
	Id string `json:"id"`
}

type Accessor interface {
	Get() ([]Log, error)
	GetByUser(id string) ([]Log, error)
	Create(log Log) (Log, error)
	Update(id int, log Log) (Log, error)
}

type Service struct {
	access Accessor
}

func NewService(access Accessor) *Service {
	return &Service{access: access}
}

func (s *Service) GetHandlerByUser(c echo.Context) error {
	var logs []Log
	usrid := c.Param("id")
	logs, err := s.access.GetByUser(usrid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, logs)
}

func (s *Service) GetHandler(c echo.Context) error {
	var logs []Log
	logs, err := s.access.Get()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, logs)
}

func (s *Service) PostHandler(c echo.Context) error {
	var log Log
	userId := c.Get("usrid")
	if err := c.Bind(&log); err != nil {
		return err
	}
	if userId == nil {
		return errors.New("")
	}
	log.Id = userId.(string)
	_, err := govalidator.ValidateStruct(&log)
	if err != nil {
		return err
	}
	log, err = s.access.Create(log)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, log)
}

func (s *Service) PutHandler(c echo.Context) error {
	var log Log
	idstr := c.Param("usrid")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return err
	}
	if err := c.Bind(&log); err != nil {
		return err
	}
	_, err = govalidator.ValidateStruct(log)
	if err != nil {
		return err
	}
	log, err = s.access.Update(id, log)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, log)
}

var jwtKey = []byte("jwtauth")

func authorization(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authorizationHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		token, err := jwt.Parse(authorizationHeader, func(token *jwt.Token) (i interface{}, e error) {
			return jwtKey, nil
		})
		if err != nil {
			return err
		}
		if !token.Valid {
			return errors.New("invalid token")
		}
		userId := token.Claims.(jwt.MapClaims)["usrid"]
		c.Set("usrid", userId)
		return handlerFunc(c)
	}
}

func (s *Service) LoginHandler(c echo.Context) error {
	var user user
	if err := c.Bind(&user); err != nil {
		return err
	}
	jti := generateRandomString(8)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":   fmt.Sprintf("%x", jti),
		"usrid": user.Id,
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, tokenString)
}

type server struct {
	svc *Service
}

func (s *server) GetLogs(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	var logs []Log
	logs, err := s.svc.access.Get()
	if err != nil {
		return &pb.GetResponse{}, err
	}
	var response []*pb.Log
	for i := 0; i < len(logs); i++ {
		log := logToGLog(&logs[i])
		response = append(response, &log)
	}
	return &pb.GetResponse{Logs: response}, nil
}

func (s *server) GetLogsByUser(c context.Context, in *pb.GetByUserRequest) (*pb.GetResponse, error) {
	var logs []Log
	usrid := in.Usrid
	logs, err := s.svc.access.GetByUser(usrid)
	if err != nil {
		return &pb.GetResponse{}, err
	}
	var response []*pb.Log
	for i := 0; i < len(logs); i++ {
		glog := logToGLog(&logs[i])
		response = append(response, &glog)
	}
	return &pb.GetResponse{Logs: response}, nil
}

func (s *server) PostLog(c context.Context, in *pb.PostRequest) (*pb.PostResponse, error) {
	var log Log
	log = gNewLogToLog(in.Log)
	log, err := s.svc.access.Create(log)
	if err != nil {
		return &pb.PostResponse{}, err
	}
	return &pb.PostResponse{Log: in.Log}, nil
}

func (s *server) UpdateLog(c context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {
	id := in.LogId
	log := gLogToLog(in.Log)
	log, err := s.svc.access.Update(int(id), log)
	if err != nil {
		return &pb.PutResponse{}, err
	}
	return &pb.PutResponse{Log: in.Log}, nil
}

func Interceptor(c context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Println("log")
	return handler(c, req)
}

func main() {

	//db := createInMemoryDB()
	db, err := createSQLDB(DIALECT, DBFILE)
	if err != nil {
		panic(err)
	}
	defer db.db.Close()
	db.db.AutoMigrate(&Log{})

	s := NewService(&db)

	viper.SetDefault("port", "8080")
	viper.AutomaticEnv()

	srvr := grpc.NewServer(grpc.UnaryInterceptor(Interceptor))
	pb.RegisterLoggerServer(srvr, &server{svc: s})

	//e.Use(middleware.Logger())
	/*
		e := echo.New()

		e.GET("/logs", s.GetHandler)
		e.POST("/logs", authorization(s.PostHandler))
		e.GET("/logs/:usrid", s.GetHandlerByUser)
		e.PUT("/logs/:id", s.PutHandler)
		e.POST("/login", s.LoginHandler)
		e.Logger.Fatal(e.Start(":" + viper.GetString("port")))
	*/
	lis, err := net.Listen("tcp", "localhost:3000")
	fmt.Println("Listening on port 3000...")
	if err != nil {
		panic(err)
	}
	if err := srvr.Serve(lis); err != nil {
		panic(err)
	}
}

func logToGLog(log *Log) pb.Log {
	return pb.Log{
		Id:                        int32(log.Model.ID),
		Usrid:                     log.Id,
		Log:                       log.Log,
		Timestamp:                 int32(log.Timestamp),
		UniquePhrase:              log.Unique_phrase,
		NotificationEmail:         log.Notification_email,
		NotificationEmailOptional: log.Notification_email_optional,
	}
}

func gNewLogToLog(log *pb.NewLog) Log {
	return Log{
		Id:                          log.Usrid,
		Log:                         log.Log,
		Notification_email:          log.NotificationEmail,
		Notification_email_optional: log.NotificationEmailOptional,
	}
}

func gLogToLog(log *pb.Log) Log {
	return Log{
		Id:                          log.Usrid,
		Log:                         log.Log,
		Timestamp:                   int64(log.Timestamp),
		Unique_phrase:               log.UniquePhrase,
		Notification_email:          log.NotificationEmail,
		Notification_email_optional: log.NotificationEmailOptional,
	}
}
