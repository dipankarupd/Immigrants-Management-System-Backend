package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dipankarupd/immigrant-management-system/pkg/config"
	"github.com/dipankarupd/immigrant-management-system/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

func SendNotification(receiver string, message string) error {
	client := config.Client
	notificationCollection := config.OpenCollection(client, "notifications")

	notification := model.Notification{
		Receiver: receiver,
		Message:  message,
	}

	_, err := notificationCollection.InsertOne(context.Background(), notification)
	return err
}

func GetNotificationsForImmigrant(w http.ResponseWriter, r *http.Request) {
	// Get the email of the immigrant from the request or any other means
	// In this example, I'll assume the email is passed as a query parameter
	email := r.URL.Query().Get("receiver")

	// Get a reference to the MongoDB client from the config package
	client := config.Client
	// Get a reference to the "notifications" collection
	notificationCollection := config.OpenCollection(client, "notifications")

	// Define a filter to query documents where "receiver" field matches the immigrant's email
	filter := bson.M{"receiver": email}

	// Create a slice to store the retrieved notifications
	var notifications []model.Notification

	// Execute the find operation and store the result in the notifications slice
	cur, err := notificationCollection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())

	// Loop through the cursor and decode each notification into the slice
	for cur.Next(context.Background()) {
		var notification model.Notification
		err := cur.Decode(&notification)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notifications = append(notifications, notification)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the notifications slice to JSON
	response, err := json.Marshal(notifications)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type and write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
