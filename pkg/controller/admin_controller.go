// controller/user_controller.go

package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dipankarupd/immigrant-management-system/pkg/config"
	"github.com/dipankarupd/immigrant-management-system/pkg/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAdmin is a controller that retrieves the admin information by username
func GetAdmin(w http.ResponseWriter, r *http.Request) {
	client := config.Client
	collection := config.OpenCollection(client, "admin")

	var admin model.Admin

	vars := mux.Vars(r)
	username := vars["username"]

	filter := bson.M{"username": username}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	if cur.Next(context.TODO()) {
		if err := cur.Decode(&admin); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Admin not found", http.StatusNotFound)
		return
	}

	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(admin)
	if err != nil {
		http.Error(w, "Error while marshaling the response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
