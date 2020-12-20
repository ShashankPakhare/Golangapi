package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type users struct {
	//ID       primitive.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Id    string `json: "username" `
	Name  string `json: "username" `
	DOB   string `json: "username" `
	PNO   string `json: "username" `
	Email string `json: "username" `
	Times int64  `json: "username" `
}

func main() {
	http.HandleFunc("/", alluser)
	http.HandleFunc("/user", user)
	http.HandleFunc("/user/", userres)
	http.HandleFunc("/contact", Contact)
	http.ListenAndServe(":3000", nil)
}

func alluser(w http.ResponseWriter, r *http.Request) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://new-Shashank:ShAsHaNk@cluster0.loytd.mongodb.net/<dbname>?retryWrites=true&w=majority"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	newdata := client.Database("newdata")
	userCollection := newdata.Collection("userDetails")

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var episodes []bson.M
	if err = cursor.All(ctx, &episodes); err != nil {
		log.Fatal(err)
	}

	for _, episodes := range episodes {
		we, err := json.Marshal(episodes)
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-type", "application/json")
		w.Write(we)
	}

}

func user(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)

	} else {
		r.ParseForm()

		id := r.FormValue("id")
		name := r.FormValue("name")
		dob := r.FormValue("dob")
		pno := r.FormValue("pno")
		email := r.FormValue("email")
		times := time.Now().Unix()

		obj := users{

			Id:    id,
			Name:  name,
			DOB:   dob,
			PNO:   pno,
			Email: email,
			Times: times,
		}

		userDetails, err := json.Marshal(obj)
		if err != nil {
			fmt.Println(err)
		}

		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://new-Shashank:ShAsHaNk@cluster0.loytd.mongodb.net/<dbname>?retryWrites=true&w=majority"))

		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)

		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)

		newdata := client.Database("newdata")
		userCollection := newdata.Collection("userDetails")

		userResult, err := userCollection.InsertOne(ctx, obj)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(userResult.InsertedID)

		w.Header().Set("Content-type", "application/json")
		w.Write(userDetails)
	}
}

func userres(w http.ResponseWriter, r *http.Request) {
	t := strings.Trim(r.URL.Path, "/user/")
	fmt.Println(t)
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://new-Shashank:ShAsHaNk@cluster0.loytd.mongodb.net/<dbname>?retryWrites=true&w=majority"))

	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	newdata := client.Database("newdata")
	userCollection := newdata.Collection("userDetails")

	cursor, err := userCollection.Find(ctx, bson.M{"id": t})
	if err != nil {
		log.Fatal(err)
	}

	var episodes []bson.M
	if err = cursor.All(ctx, &episodes); err != nil {
		log.Fatal(err)
	}

	for _, episodes := range episodes {
		we, err := json.Marshal(episodes)
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-type", "application/json")
		w.Write(we)
	}

}

func Contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		userid := r.FormValue("user")
		times := r.FormValue("infection_timestamp")
		if userid == "" {
			t, _ := template.ParseFiles("contact.html")
			t.Execute(w, nil)

		} else {

			fmt.Println(userid, times)

			fmt.Println(reflect.TypeOf(times))
			t, err := time.Parse("2006-01-02", times)

			if err != nil {
				fmt.Println(err)
			}

			timespam := t.Add(-24 * 14 * time.Hour).Unix()
			//{ $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }

			client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://new-Shashank:ShAsHaNk@cluster0.loytd.mongodb.net/<dbname>?retryWrites=true&w=majority"))

			if err != nil {
				log.Fatal(err)
			}
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			err = client.Connect(ctx)

			if err != nil {
				log.Fatal(err)
			}
			defer client.Disconnect(ctx)

			newdata := client.Database("newdata")
			userCollection := newdata.Collection("contactDetails")

			opts := options.Find()
			opts.SetSort(bson.D{{"duration", 1}})

			sortCursor, err := userCollection.Find(ctx, bson.D{
				{"time", bson.D{
					{"$gt", timespam},
				}}, // {"useridOne", userid},
			}, opts)

			var episodesSorted []bson.M
			if err = sortCursor.All(ctx, &episodesSorted); err != nil {
				log.Fatal(err)
			}

			for _, episodesSorted := range episodesSorted {

				fmt.Println(episodesSorted["useridTwo"])
				fmt.Println(reflect.TypeOf(episodesSorted["useridTwo"]))

				Collect := newdata.Collection("userDetails")

				if episodesSorted["useridTwo"] == userid {
					cursor, err := Collect.Find(ctx, bson.M{"id": episodesSorted["useridOne"]})
					if err != nil {
						log.Fatal(err)
					}

					var epis []bson.M
					if err = cursor.All(ctx, &epis); err != nil {
						log.Fatal(err)
					}

					for _, epis := range epis {
						wer, err := json.Marshal(epis)
						if err != nil {
							fmt.Println(err)
						}

						w.Header().Set("Content-type", "application/json")
						w.Write(wer)
					}
				} else if episodesSorted["useridOne"] == userid {

					cursor, err := Collect.Find(ctx, bson.M{"id": episodesSorted["useridTwo"]})
					if err != nil {
						log.Fatal(err)
					}

					var epis []bson.M
					if err = cursor.All(ctx, &epis); err != nil {
						log.Fatal(err)
					}

					for _, epis := range epis {
						wer, err := json.Marshal(epis)
						if err != nil {
							fmt.Println(err)
						}

						w.Header().Set("Content-type", "application/json")
						w.Write(wer)
					}
				}

			}
		}

	} else {
		r.ParseForm()

		userOne := r.FormValue("useridOne")
		userTwo := r.FormValue("useridTwo")
		timestamp := r.FormValue("Timestamp")

		t, err := time.Parse("2006-01-02", timestamp)
		if err != nil {
			log.Fatal(err)
		}

		times := t.Unix()

		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://new-Shashank:ShAsHaNk@cluster0.loytd.mongodb.net/<dbname>?retryWrites=true&w=majority"))

		if err != nil {
			log.Fatal(err)
		}

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)

		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)

		newdata := client.Database("newdata")
		userCollection := newdata.Collection("contactDetails")

		userResult, err := userCollection.InsertOne(ctx, bson.D{
			{"useridOne", userOne},
			{"useridTwo", userTwo},
			{"time", times},
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(userResult.InsertedID)

		cursor, err := userCollection.Find(ctx, bson.M{"_id": userResult.InsertedID})
		if err != nil {
			log.Fatal(err)
		}

		var episodes []bson.M
		if err = cursor.All(ctx, &episodes); err != nil {
			log.Fatal(err)
		}

		for _, episodes := range episodes {
			we, err := json.Marshal(episodes)
			if err != nil {
				fmt.Println(err)
			}

			w.Header().Set("Content-type", "application/json")
			w.Write(we)
		}
	}
}
