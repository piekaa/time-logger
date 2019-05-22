package main

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"fmt"
)

func main() {
	// Echo instance

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	s := &Server{}

	s.initDb()

	// Routes
	e.POST("/start-counter", s.startCounter)
	e.POST("/stop-counter", s.stopCounter)
	e.Static("/static", "static")

	// Start server
	e.Logger.Fatal(e.Start(":8882"))
}

type Server struct {
	client *mongo.Client
	collection *mongo.Collection
}

func (s *Server) initDb() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	s.client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = s.client.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB")
	s.collection = s.client.Database("time_logger").Collection("checkouts")
}

func (*Server) startCounter(c echo.Context) error {
	return c.String(http.StatusOK, ":)")
}

func (*Server) stopCounter(c echo.Context) error {
	return c.String(http.StatusOK, ":)")
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
