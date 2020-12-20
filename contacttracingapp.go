package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
	return []byte(stamp), nil
}

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Dob         string             `json:"dob,omitempty" bson: "dob,omitempty"`
	PhoneNumber string             `json:"phone_num,omitempty" bson: "phone_num,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Time_stamp  time.Time          `json:"time_stamp" bson: "time_stamp"`
}
type Contact struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId1    string             `json:"user_id_1,omitempty" bson:"user_id_1,omitempty"`
	UserId2    string             `json:"user_id_2,omitempty" bson: "user_id_2,omitempty"`
	Time_stamp time.Time          `json:"time_stamp" bson: "time_stamp"`
}

// type User struct {

func main() {
	fmt.Println("Starting the application...")
	// for connecting to mongodb_server
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://new-Shashank:ShAsHaNk@cluster0.loytd.mongodb.net/<dbname>?retryWrites=true&w=majority")
	client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()

	router.HandleFunc("/users", AllUsers).Methods("GET")
	router.HandleFunc("/users", CreatenewUser).Methods("POST")

	router.HandleFunc("/users/{id}", FindUser).Methods("GET")

	router.HandleFunc("/contacts", AllContacts).Methods("GET")
	router.HandleFunc("/contacts", CreatenewContact).Methods("POST")

	router.HandleFunc("/contacts?users={id}&infection_timestamp={ts}", findContact).Methods("GET")

	http.ListenAndServe(":4000", router)
}

func FindUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])
	var users User
	collection := client.Database("newFormData").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err := collection.FindOne(ctx, User{ID: id}).Decode(&users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(users)
}

func AllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var users []User
	collection := client.Database("newFormData").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var users User
		cursor.Decode(&users)
		users = append(users, users)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(users)
}

func CreatenewContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var contact Contact
	_ = json.NewDecoder(r.Body).Decode(&contact)

	contact.Time_stamp = time.Now()

	collection := client.Database("newFormData").Collection("contacts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, _ := collection.InsertOne(ctx, contact)
	json.NewEncoder(w).Encode(result)
}

func findContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])
	var contact Contact
	collection := client.Database("newFormData").Collection("contacts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err := collection.FindOne(ctx, Contact{ID: id}).Decode(&contact)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(contact)
}

func AllContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var contacts []Contact
	collection := client.Database("newFormData").Collection("contacts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var contact Contact
		cursor.Decode(&contact)
		contacts = append(contacts, contact)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(contacts)
}

func CreatenewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var users User
	_ = json.NewDecoder(r.Body).Decode(&users)

	users.Time_stamp = time.Now()

	collection := client.Database("newFormData").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, _ := collection.InsertOne(ctx, users)
	json.NewEncoder(w).Encode(result)
}
