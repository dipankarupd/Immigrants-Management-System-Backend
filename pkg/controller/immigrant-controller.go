package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dipankarupd/immigrant-management-system/pkg/config"
	"github.com/dipankarupd/immigrant-management-system/pkg/model"
	"github.com/dipankarupd/immigrant-management-system/pkg/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func GetImmigrants(w http.ResponseWriter, r *http.Request) {
	// Get a reference to the MongoDB client from the config package
	client := config.Client
	// Get a reference to the "immigrants" collection
	collection := config.OpenCollection(client, "demo")

	// Create a slice to store the retrieved immigrants
	var immigrants []model.Immigrant

	// Define a filter (if needed) to query specific data from the collection
	filter := bson.M{} // Empty filter to retrieve all documents

	// Execute the find operation and store the result in the immigrants slice
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.TODO())

	// context.TODO() returns a non nil empty context
	// this way we can terminate the loop until all the cursors are iterated
	for cur.Next(context.TODO()) {
		var immigrant model.Immigrant
		err := cur.Decode(&immigrant)
		if err != nil {
			log.Fatal(err)
		}
		immigrants = append(immigrants, immigrant)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Convert the immigrants slice to JSON
	response, err := json.Marshal(immigrants)
	if err != nil {
		log.Fatal(err)
	}

	// Set the content type and write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func GetImmigrant(w http.ResponseWriter, r *http.Request) {
	client := config.Client
	collection := config.OpenCollection(client, "demo")

	var immigrant model.Immigrant

	vars := mux.Vars(r)
	passportnum := vars["passportno"]

	// convert the passport num gotten from the url to int from string

	ppn, err := strconv.Atoi(passportnum)
	if err != nil {
		http.Error(w, "Invalid passport number", http.StatusBadRequest)
	}
	filter := bson.M{"passportno": ppn}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	if cur.Next(context.TODO()) {
		if err := cur.Decode(&immigrant); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Immigrant not found", http.StatusNotFound)
		return
	}

	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(immigrant)
	if err != nil {
		http.Error(w, "Error while marshaling the response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func CreateImmigrant(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a model.Immigrant struct
	// unmarshalling
	immigrant := model.Immigrant{}
	util.ParseBody(r, &immigrant)
	immigrant.SetDefaultValue()

	immigrant.ID = bson.NewObjectId()

	// Get the MongoDB collection
	client := config.Client
	collection := config.OpenCollection(client, "demo")

	// check if the index value is unique or not:
	errr := util.CreateUniqueIndex(collection, "passportno")
	if errr != nil {
		http.Error(w, errr.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the immigrant document into the collection
	_, err := collection.InsertOne(context.TODO(), immigrant)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, "Passport number already exists.", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(immigrant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(response)

}

func AcceptImmigrant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	passportnum := vars["passportno"]
	ppn, er := strconv.Atoi(passportnum)

	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	client := config.Client
	collection := config.OpenCollection(client, "demo")

	filter := bson.M{"passportno": ppn}
	immigrant := model.Immigrant{}

	err := collection.FindOne(context.TODO(), filter).Decode(&immigrant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Immigrant not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get the desired approval status from the request body
	approval, err := util.GetApprovalFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the immigrant's approval status
	immigrant.Approval = &approval

	// Update the document in the collection
	update := bson.M{"$set": bson.M{"approval": immigrant.Approval}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated immigrant object as the response
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(immigrant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
