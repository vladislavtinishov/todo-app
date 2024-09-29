package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vladislavtinishov/todo-app"
	"github.com/vladislavtinishov/todo-app/pkg/handler"
	"github.com/vladislavtinishov/todo-app/pkg/repository"
	"github.com/vladislavtinishov/todo-app/pkg/service"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error initializing env: %s", err.Error())
	}

	db, err := repository.NewMysqlDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	})

	if err != nil {
		logrus.Fatalf("Error opening db connection: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	src := new(todo.Server)

	if err := src.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while runnig http server %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
