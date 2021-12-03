package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
var mongoURL = "mongodb://localhost:27017"
var collection *mongo.Collection

type Employee struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Role string             `json:"role,omitempty" bson:"role,omitempty"`
	Flag string              `json:"flag,omitempty" bson:"flag,omitempty"`
}

type Employees []Employee

func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprint(response, "Welcome to Homepage")
}

func CreateEmployee(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var employee Employee
	json.NewDecoder(request.Body).Decode(&employee)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, employee)
	json.NewEncoder(response).Encode(result)
}

func GetEmployees(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var employees []Employee
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": " `+err.Error()+`"}`))
		return 
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var employee Employee
		cursor.Decode(&employee)
		employees = append(employees, employee)
	}
	if err:= cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": " `+err.Error()+`"}`))
		return
	}
	json.NewEncoder(response).Encode(employees)
}

func connectMongo() {
	//initialize new mongo with options
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Connected to MongoDB server: ", mongoURL)
	collection = client.Database("golangproject").Collection("employees")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/empcreate", CreateEmployee).Methods("POST")
	myRouter.HandleFunc("/showemp", GetEmployees).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	fmt.Println("Starting Application")
	connectMongo()
	handleRequests()
}
