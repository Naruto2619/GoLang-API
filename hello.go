package main

import (
	"Appointy/helper"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	key, err := helper.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	encPass, err := helper.Encrypt(key, []byte(U1.Password))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ciphertext: %s\n", hex.EncodeToString(encPass))
	U1.Password = hex.EncodeToString(encPass)
	var collection = helper.ConnectDBuser()
	result, er := collection.InsertOne(context.TODO(), U1)
	if er != nil {
		fmt.Println(er)
	}
	json.NewEncoder(response).Encode(result)
}

func GetUser(response http.ResponseWriter, request *http.Request) {
	var urlid = strings.TrimPrefix(request.URL.Path, "/user/")
	id, _ := primitive.ObjectIDFromHex(urlid)
	response.Header().Set("content-type", "application/json")
	var users []User
	var collection = helper.ConnectDBuser()
	cur, err := collection.Find(context.TODO(), bson.M{"_id": id})
	if err != nil {
		helper.GetError(err, response)
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
	P1.Timestamp = time.Now().String()
	var collection = helper.ConnectDBpost()
	result, er := collection.InsertOne(context.TODO(), P1)
	if er != nil {
		fmt.Println(er)
	}
	json.NewEncoder(response).Encode(result)
}
func GetPost(response http.ResponseWriter, request *http.Request) {
	var urlid = strings.TrimPrefix(request.URL.Path, "/posts/") // get id from http params
	id, _ := primitive.ObjectIDFromHex(urlid)                   // convert id from string to primitive type
	response.Header().Set("content-type", "application/json")
	var posts []Post
	var collection = helper.ConnectDBpost()
	cur, err := collection.Find(context.TODO(), bson.M{"_id": id}) //fetch posts filtering based on user id
	if err != nil {
		helper.GetError(err, response)
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
	var userid = strings.TrimPrefix(request.URL.Path, "/posts/user/") // get id from http params
	var limitstr = request.URL.Query().Get("limit")
	limit, err := strconv.ParseInt(limitstr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	response.Header().Set("content-type", "application/json")
	var posts []Post
	var collection = helper.ConnectDBpost()
	cur, err := collection.Find(context.TODO(), bson.M{"User_id": userid}) //fetch posts filtering based on user id
	if err != nil {
		helper.GetError(err, response)
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
	json.NewEncoder(response).Encode(posts[:limit])
}

func handleRequests() {
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

func main() {
	handleRequests()
}
