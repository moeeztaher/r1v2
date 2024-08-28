package Handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	//"github.com/gorilla/mux"
	"GUI-GO/Apis"
	//"io/ioutil"
	//"os"
	//"path/filepath"
	//"log"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func PublishServiceHandler(serviceCollection *mongo.Collection, rappCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var api Apis.PublishServiceAPI

		vars := mux.Vars(r)
		apfId := vars["apfId"]

		rappFilter := bson.D{{"apf_id", apfId}}
		rappCursor, err := rappCollection.Find(context.TODO(), rappFilter)
		if err != nil {
			http.Error(w, "Failed to retrieve documents from MongoDB", http.StatusInternalServerError)
			return
		}
		defer rappCursor.Close(context.TODO())

		var rapp Apis.Rapp
		rappExists := false

		for rappCursor.Next(context.TODO()) {
			if err = rappCursor.Decode(&rapp); err != nil {
				http.Error(w, "Failed to decode document", http.StatusInternalServerError)
				return
			}
			rappExists = true
		}

		if !rappExists {
			http.Error(w, fmt.Sprintf("Invoker ID: %v not found", apfId), http.StatusNotFound)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&api)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if api.APIName == "" {
			http.Error(w, "apiName field is required", http.StatusBadRequest)
			return
		}

		if api.APIID == "" {
			api.APIID = uuid.New().String()
		}

		// TODO: better naming for filters/cursors. maybe move this to separate function
		rappFilter2 := bson.M{"apf_id": apfId, "authorized_services": api.APIID}
		rappCursor2, err := rappCollection.CountDocuments(context.TODO(), rappFilter2)
		if err != nil {
			http.Error(w, "Failed to retrieve documents from MongoDB", http.StatusInternalServerError)
			return
		}
		if rappCursor2 == 1 {
			http.Error(w, fmt.Sprintf("The service: %v already exists for rapp: %v.", api.APIID, apfId), http.StatusInternalServerError)
			return
		}

		rappUpdate := bson.D{
			{"$push", bson.D{
				{"authorized_services", api.APIID},
			}},
		}
		_, err = rappCollection.UpdateOne(context.TODO(), rappFilter, rappUpdate)
		if err != nil {
			http.Error(w, "Failed to insert document into MongoDB", http.StatusInternalServerError)
			return
		}

		_, err = serviceCollection.InsertOne(context.TODO(), api)
		if err != nil {
			http.Error(w, "Failed to insert document into MongoDB", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// TODO: Change the location to the actual location of created service
		w.Header().Set("Location", "https://api.service-apis/"+api.APIID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Service API published successfully"})
	}
}

func GetServiceAPIsHandler(serviceCollection *mongo.Collection, rappCollection *mongo.Collection) func(r *http.Request) ([]byte, error) {
	return func(r *http.Request) ([]byte, error) {
		var serviceAPIs []Apis.GetServiceAPI
		invokerId := r.URL.Query().Get("api-invoker-id")

		rappFilter := bson.D{{"apf_id", invokerId}}
		rappCursor, err := rappCollection.Find(context.TODO(), rappFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve documents from MongoDB: %w", err)
		}
		defer rappCursor.Close(context.TODO())

		var rapp Apis.Rapp
		rappExists := false

		for rappCursor.Next(context.TODO()) {
			if err = rappCursor.Decode(&rapp); err != nil {
				return nil, fmt.Errorf("failed to decode document: %w", err)
			}
			rappExists = true
		}

		if !rappExists {
			return nil, fmt.Errorf("invoker ID not found")
		}

		if rapp.IsAuthorized == false {
			return nil, fmt.Errorf("unauthorized")
		}

		serviceFilter := bson.D{{"apiid", bson.D{{"$in", rapp.AuthorizedServices}}}}

		serviceCursor, err := serviceCollection.Find(context.TODO(), serviceFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve documents from MongoDB: %w", err)
		}
		defer serviceCursor.Close(context.TODO())

		for serviceCursor.Next(context.TODO()) {
			var api Apis.GetServiceAPI
			if err = serviceCursor.Decode(&api); err != nil {
				return nil, fmt.Errorf("failed to decode document: %w", err)
			}
			serviceAPIs = append(serviceAPIs, api)
		}

		if err := serviceCursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}

		jsonData, err := json.Marshal(serviceAPIs)
		if err != nil {
			return nil, fmt.Errorf("failed to encode JSON: %w", err)
		}

		return jsonData, nil
	}
}

func GetSpecificServiceAPIHandler(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		//apfId := params["apfId"]
		serviceApiId := params["serviceApiId"]

		var serviceAPI Apis.GetServiceAPI

		filter := bson.M{"apiId": serviceApiId}
		err := collection.FindOne(context.TODO(), filter).Decode(&serviceAPI)
		if err != nil {
			http.Error(w, "Failed to retrieve document from MongoDB", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(serviceAPI)
	}
}

func UpdateServiceAPIHandler(serviceCollection *mongo.Collection, rappCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		apfId := vars["apfId"]
		serviceApiId := vars["serviceApiId"]

		// TODO: update this to use one filter
		rappFilter := bson.D{{"apf_id", apfId}}
		rappCursor, err := rappCollection.CountDocuments(context.TODO(), rappFilter)
		if err != nil {
			http.Error(w, "Failed to retrieve documents from MongoDB", http.StatusInternalServerError)
			return
		}
		if rappCursor == 0 {
			http.Error(w, fmt.Sprintf("The specified rapp: %v does not exist.", apfId), http.StatusInternalServerError)
			return
		}
		rappFilter2 := bson.M{"apf_id": apfId, "authorized_services": serviceApiId}
		rappCursor, err = rappCollection.CountDocuments(context.TODO(), rappFilter2)
		if rappCursor == 0 {
			http.Error(w, "The specified service does not exist or it is not authorized for this rapp.", http.StatusInternalServerError)
			return
		}

		var apiData Apis.PublishServiceAPI
		err = json.NewDecoder(r.Body).Decode(&apiData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filter := bson.M{"apiId": serviceApiId}
		update := bson.M{
			"$set": apiData,
		}

		_, err = serviceCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{"message": fmt.Sprintf("Service API %s updated", serviceApiId)}
		json.NewEncoder(w).Encode(response)
	}
}

func PatchServiceAPIHandler(serviceCollection *mongo.Collection, rappCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		apfId := vars["apfId"]
		serviceApiId := vars["serviceApiId"]

		rappFilter := bson.D{{"apf_id", apfId}}
		rappCursor, err := rappCollection.CountDocuments(context.TODO(), rappFilter)
		if err != nil {
			http.Error(w, "Failed to retrieve documents from MongoDB", http.StatusInternalServerError)
			return
		}
		if rappCursor == 0 {
			http.Error(w, fmt.Sprintf("The specified rapp: %v does not exist.", apfId), http.StatusInternalServerError)
			return
		}
		rappFilter2 := bson.M{"apf_id": apfId, "authorized_services": serviceApiId}
		rappCursor, err = rappCollection.CountDocuments(context.TODO(), rappFilter2)
		if rappCursor == 0 {
			http.Error(w, "The specified service does not exist or it is not authorized for this rapp.", http.StatusInternalServerError)
			return
		}

		var patchReq bson.M
		if err := json.NewDecoder(r.Body).Decode(&patchReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filter := bson.M{"apiId": serviceApiId}

		update := bson.M{"$set": patchReq}
		result, err := serviceCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if result.MatchedCount == 0 {
			http.Error(w, "No documents found with the given ID", http.StatusNotFound)
			return
		}

		response := map[string]string{"message": fmt.Sprintf("Service API %s updated", serviceApiId)}
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteServiceAPIHandler(serviceCollection *mongo.Collection, rappCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceApiId := vars["serviceApiId"]
		rappId := vars["apfId"]

		serviceFilter := bson.M{"apiId": serviceApiId}
		// TODO: add transaction to make these operations atomic
		result, err := serviceCollection.DeleteOne(context.Background(), serviceFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			http.Error(w, "No documents found with the given ID", http.StatusNotFound)
			return
		}

		rappFilter := bson.M{"apf_id": rappId, "authorized_services": serviceApiId}
		update := bson.M{"$pull": bson.M{"authorized_services": serviceApiId}}

		_, err = rappCollection.UpdateOne(context.TODO(), rappFilter, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{"message": fmt.Sprintf("Service API %s deleted", serviceApiId)}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
