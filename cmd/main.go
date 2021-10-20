package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/sinoturaev/todo-app"
	"github.com/sinoturaev/todo-app/pkg/handler"
	"github.com/sinoturaev/todo-app/pkg/repository"
	"github.com/sinoturaev/todo-app/pkg/service"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "sino2611",
		DBName:   "todo-db",
		SSLMode:  "disable",

	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo_app.Server)
	go func() {
		if err := srv.Run("8000", handlers.InitRouters()); err != nil {
			logrus.Fatalf("Error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occured on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("Error occured on db connecting close: %s", err.Error())
	}
}


