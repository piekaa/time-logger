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
	"go.mongodb.org/mongo-driver/bson"
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
	e.Static("/", "static")

	// Start server
	e.Logger.Fatal(e.Start(":8882"))
}

type Server struct {
	client    *mongo.Client
	checkouts *mongo.Collection
	status    *mongo.Collection
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
	s.checkouts = s.client.Database("time_logger").Collection("checkouts")
	s.status = s.client.Database("time_logger").Collection("status")
}

func (s *Server) startCounter(c echo.Context) error {
	m := echo.Map{}
	c.Bind(&m)
	name := fmt.Sprintf("%v", m["name"])
	fmt.Println(name)
	//todo handle errors
	result, _ := s.checkouts.InsertOne(s.getCtx(), bson.M{"name": name, "start": currentTimeMillis(), "_id": currentTimeMillis()})
	_, err := s.status.UpdateOne(s.getCtx(), bson.M{"name": name}, bson.M{"$set": bson.M{"name": name, "status": "working", "id": result.InsertedID}}, upsert())

	if err != nil {
		fmt.Println(err)
	}

	return c.String(http.StatusOK, "{}")
}

func (s *Server) stopCounter(c echo.Context) error {

	m := echo.Map{}
	c.Bind(&m)
	name := fmt.Sprintf("%v", m["name"])
	fmt.Println(name)
	response := s.status.FindOne(s.getCtx(), bson.M{"name": name})

	result := struct {
		Id int64 `json:"_id"`
	}{}
	response.Decode(&result)
	s.checkouts.UpdateOne(s.getCtx(), bson.M{"_id": result.Id}, bson.M{"$set": bson.M{"stop": currentTimeMillis()}})
	s.status.UpdateOne(s.getCtx(), bson.M{"name": name}, bson.M{"$set": bson.M{"name": name, "status": "not working", "id": ""}}, upsert())
	return c.String(http.StatusOK, "{}")
}

func (s *Server) getCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func upsert() *options.UpdateOptions {
	b := true
	return &options.UpdateOptions{Upsert: &b}
}
