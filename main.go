package main

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

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

func main() {
	db, err := createSQLDB(DIALECT, DBFILE)
	if err != nil {
		panic(err)
	}
	defer db.db.Close()
	db.db.AutoMigrate(&Log{})
	s := NewService(&db)

	viper.SetDefault("port", "8080")
	viper.AutomaticEnv()

	e := echo.New()
	//e.Use(middleware.Logger())

	e.GET("/logs", s.GetHandler)
	e.POST("/logs", authorization(s.PostHandler))
	e.GET("/logs/:usrid", s.GetHandlerByUser)
	e.PUT("/logs/:id", s.PutHandler)
	e.POST("/login", s.LoginHandler)
	e.Logger.Fatal(e.Start(":" + viper.GetString("port")))
}
