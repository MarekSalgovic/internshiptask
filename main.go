package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

type Message struct {
	gorm.Model
	Id        string `json:"id" gorm:"column:idmsg"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Accessor interface {
	Get() ([]Message, error)
	Create(message Message) (Message, error)
}

type Service struct {
	access Accessor
}

func NewService(access Accessor) *Service {
	return &Service{access: access}
}

func (s *Service) GetHandler(c echo.Context) error {
	var messages []Message
	messages, err := s.access.Get()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, messages)
}

func (s *Service) PostHandler(c echo.Context) error {
	var message Message
	if err := c.Bind(&message); err != nil {
		return err
	}
	message, err := s.access.Create(message)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, message)
}

func main() {
	db, err := createSQLDB()
	if err != nil {
		panic(err)
	}
	defer db.db.Close()

	if !db.db.HasTable(&Message{}) {
		db.db.CreateTable(&Message{})
	}


	s := NewService(db)
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/", s.GetHandler)
	e.POST("/", s.PostHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
