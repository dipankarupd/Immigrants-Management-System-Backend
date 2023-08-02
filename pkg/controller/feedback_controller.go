package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dipankarupd/immigrant-management-system/pkg/config"
	"github.com/dipankarupd/immigrant-management-system/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateFeedback(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into the Feedback struct
	feedback := model.Feedback{}
	err := json.NewDecoder(r.Body).Decode(&feedback)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Check if the provided passportno exists in the "demo" collection
	client := config.Client
	immigrantCollection := config.OpenCollection(client, "demo")
	filter := bson.M{"passportno": feedback.PassportNo}
	count, err := immigrantCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		http.Error(w, "Error checking passportno existence", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Passport number not found. Feedback can only be given for existing immigrants.", http.StatusNotFound)
		return
	}

	// Retrieve the Immigrant with the provided passportno
	immigrant := model.Immigrant{}
	err = immigrantCollection.FindOne(context.Background(), filter).Decode(&immigrant)
	if err != nil {
		http.Error(w, "Error retrieving immigrant details", http.StatusInternalServerError)
		return
	}

	// Set the ImmigrantID in the Feedback struct
	feedback.ImmigrantID = immigrant.ID

	// Get the Feedback collection
	feedbackCollection := config.OpenCollection(client, "feedbacks")

	// Insert the feedback document into the collection
	_, err = feedbackCollection.InsertOne(context.Background(), feedback)
	if err != nil {
		log.Printf("Error inserting feedback: %v", err)
		http.Error(w, "Error inserting feedback", http.StatusInternalServerError)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshal the feedback as the response
	response, err := json.Marshal(feedback)
	if err != nil {
		log.Printf("Error marshaling the response: %v", err)
		http.Error(w, "Error marshaling the response", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
