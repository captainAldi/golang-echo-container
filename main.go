package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	"golang-echo-container/utils/connection"
)

func main() {

	// Load ENV
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Error !, config file: %w \n", err))
	}

	// Define Variable

	// DB
	db_type := viper.GetString("database.type")
	db_host := viper.GetString("database.host")
	db_name := viper.GetString("database.name")
	db_user := viper.GetString("database.user")
	db_pass := viper.GetString("database.pass")
	db_port := viper.GetString("database.port")
	db_auth_source := viper.GetString("database.auth_source")

	// Redis
	redis_host := viper.GetString("redis.host")
	redis_port := viper.GetString("redis.port")
	redis_user := viper.GetString("redis.user")
	redis_pass := viper.GetString("redis.pass")

	// RabbitMQ
	rabbitmq_host := viper.GetString("rabbitmq.host")
	rabbitmq_port := viper.GetString("rabbitmq.port")
	rabbitmq_user := viper.GetString("rabbitmq.user")
	rabbitmq_pass := viper.GetString("rabbitmq.pass")

	// Create echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health", func(c echo.Context) error {

		// Check DB Type
		checkingDBStat := "status"

		if db_type == "mysql" {
			checkingDBStat = connection.CheckMySQL(db_user, db_pass, db_name, db_host, db_port)
		} else if db_type == "postgres" {
			checkingDBStat = connection.CheckPostgreSQL(db_user, db_pass, db_name, db_host, db_port)
		} else if db_type == "mongo" {
			checkingDBStat = connection.CheckMongo(db_user, db_pass, db_name, db_host, db_port, db_auth_source)
		} else {
			checkingDBStat = fmt.Sprintf("Type: [%s] Not Supported !", db_type)
		}

		// Struct JSON Response
		type Connection struct {
			Db_status            string `json:"db_status"`
			Redis_status         string `json:"redis_status"`
			Rabbitmq_status      string `json:"rabbitmq_status"`
			Version_status       string `json:"version_status"`
			Latest_update_status string `json:"latest_update_status"`
		}

		connectionStatus := &Connection{
			Db_status:            checkingDBStat,
			Redis_status:         connection.CheckRedis(redis_host, redis_port, redis_user, redis_pass),
			Rabbitmq_status:      connection.CheckRabbitMQ(rabbitmq_host, rabbitmq_port, rabbitmq_user, rabbitmq_pass),
			Version_status:       "v0.0.2",
			Latest_update_status: "CI/CD to Rara",
		}
		return c.JSON(http.StatusOK, connectionStatus)
	})

	// Start Echo
	e.Logger.Fatal(e.Start(":1223"))

}
