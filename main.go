package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	if err := c.Bind(&log); err != nil {
		return err
	}
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

func main() {
	db, err := createSQLDB(DIALECT, DBFILE)
	if err != nil {
		panic(err)
	}
	defer db.db.Close()

	db.db.AutoMigrate(&Log{})

	s := NewService(&db)
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/logs", s.GetHandler)
	e.POST("/logs", s.PostHandler)
	e.GET("/logs/:usrid", s.GetHandlerByUser)
	e.PUT("/logs/:id", s.PutHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
