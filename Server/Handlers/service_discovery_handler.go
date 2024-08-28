package Handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var serviceAPIsCollection *mongo.Collection

func InitServiceAPIsCollection(collection *mongo.Collection) {
	serviceAPIsCollection = collection
}

func ServiceDiscoveryHandler(serviceCollection *mongo.Collection, rappCollection *mongo.Collection, apiInvokerId *string, apiName *string,
	apiVersion *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rappFilter := bson.M{}
		serviceFilter := bson.M{}

		if *apiInvokerId != "" {
			rappFilter["apf_id"] = *apiInvokerId
		} else {
			http.Error(w, "The api-invoker-id parameter is required.", http.StatusInternalServerError)
			return
		}

		// TODO: replace with FindOne
		rappCursor, findErr := rappCollection.Find(context.TODO(), rappFilter)
		if findErr != nil {
			http.Error(w, findErr.Error(), http.StatusInternalServerError)
		}
		if !rappCursor.TryNext(context.TODO()) {
			http.Error(w, fmt.Sprintf("The specified rapp: %v does not exist.", rappFilter["apf_id"]), http.StatusNotFound)
			return
		}

		var results []bson.M
		if err := rappCursor.All(context.TODO(), &results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if *apiName != "" {
			serviceFilter["apiName"] = *apiName
		}

		// if *apiVersion != "" {
		// 	filter["aefprofiles.versions.apiversion"] = apiVersion
		// }

		serviceCursor, err := serviceCollection.Find(context.TODO(), serviceFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer serviceCursor.Close(context.TODO())

		// Check if any services were found
		if !serviceCursor.TryNext(context.TODO()) {
			http.Error(w, "No documents found", http.StatusNotFound)
			return
		}

		// Process the documents
		var serviceResults []bson.M
		if err := serviceCursor.All(context.TODO(), &serviceResults); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		respErr := json.NewEncoder(w).Encode(serviceResults)
		if respErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func handleError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorMessage := fmt.Sprintf(`{"error": "%s", "message": "%s"}`, message, err.Error())
	w.Write([]byte(errorMessage))
}
