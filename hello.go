package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"Name,omitempty" bson:"Name,omitempty"`
	Email    string             `json:"Email,omitempty" bson:"Email,omitempty"`
	Password string             `json:"Password,omitempty" bson:"Password,omitempty"`
}

type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption   string             `json:"Caption,omitempty" bson:"Caption,omitempty"`
	ImageUrl  string             `json:"ImageUrl,omitempty" bson:"ImageUrl,omitempty"`
	Timestamp string             `json:"Timestamp,omitempty" bson:"Timestamp,omitempty"`
	User_id   string             `json:"User_id,omitempty" bson:"User_id,omitempty"`
}

func CreateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var U1 User
	_ = json.NewDecoder(request.Body).Decode(&U1)
	var collection = ConnectDBuser()
	result, er := collection.InsertOne(context.TODO(), U1)
	if er != nil {
		fmt.Println(er)
	}
	json.NewEncoder(response).Encode(result)
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}
	message, _ := json.Marshal(response)
	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
func GetUser(response http.ResponseWriter, request *http.Request) {
	var id2 = request.URL.Query().Get("id")
	id, _ := primitive.ObjectIDFromHex(id2)
	response.Header().Set("content-type", "application/json")
	var users []User
	var collection = ConnectDBuser()
	cur, err := collection.Find(context.TODO(), bson.M{"_id": id})
	if err != nil {
		GetError(err, response)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var us User
		// & character returns the memory address of the following variable.
		err := cur.Decode(&us) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		users = append(users, us)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(users)
}
func CreatePost(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var P1 Post
	_ = json.NewDecoder(request.Body).Decode(&P1)
	var collection = ConnectDBpost()
	result, er := collection.InsertOne(context.TODO(), P1)
	if er != nil {
		fmt.Println(er)
	}
	json.NewEncoder(response).Encode(result)
}
func GetPost(response http.ResponseWriter, request *http.Request) {
	var id2 = request.URL.Query().Get("id")
	id, _ := primitive.ObjectIDFromHex(id2)
	response.Header().Set("content-type", "application/json")
	var posts []Post
	var collection = ConnectDBpost()
	cur, err := collection.Find(context.TODO(), bson.M{"_id": id})
	if err != nil {
		GetError(err, response)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {

		var p Post
		err := cur.Decode(&p) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		posts = append(posts, p)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(posts)
}
func UserPost(response http.ResponseWriter, request *http.Request) {

	var id2 = request.URL.Query().Get("userid")
	response.Header().Set("content-type", "application/json")
	var posts []Post
	var collection = ConnectDBpost()
	cur, err := collection.Find(context.TODO(), bson.M{"User_id": id2})
	if err != nil {
		GetError(err, response)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {

		var p Post
		err := cur.Decode(&p) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		posts = append(posts, p)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(posts)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/user", CreateUser)
	http.HandleFunc("/user/", GetUser)
	http.HandleFunc("/posts", CreatePost)
	http.HandleFunc("/posts/", GetPost)
	http.HandleFunc("/posts/user/", UserPost)
	err := http.ListenAndServe(":6969", nil)
	if err != nil {
		panic(err)
	}
}
func ConnectDBuser() *mongo.Collection {
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("API_KEY"))
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection := client.Database("mydb").Collection("instauser")
	return collection
}
func ConnectDBpost() *mongo.Collection {
	// Set client optionss
	clientOptions := options.Client().ApplyURI("mongodb+srv://siddharth:Naruto2619*@cluster0.gqdbh.mongodb.net/mydb?retryWrites=true&w=majority")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection := client.Database("mydb").Collection("instapost")
	return collection
}

func main() {
	fmt.Println("Hello nigga")
	handleRequests()
}
