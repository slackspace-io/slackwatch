package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slackwatch/backend/internal/model"
	"slackwatch/backend/internal/storage"
	"strings"
)

type DataService struct {
    Storage storage.Storage
}

func (ds *DataService) SaveData(ctx context.Context, data interface{}) error {
    file, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    return encoder.Encode(data)
}

func (ds *DataService) GetData(ctx context.Context) ([]map[string]string, error) {
    file, err := os.Open("data.json")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var data []map[string]string
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&data)
    if err != nil {
        return nil, err
    }

    return data, nil
}

// HandleSaveData handles the HTTP request for saving data
func (ds *DataService) HandleSaveData(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var data model.Data
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := ds.SaveData(context.Background(), data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

// HandleGetData handles the HTTP request for retrieving data by ID
func (ds *DataService) HandleGetData(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    id := strings.TrimPrefix(r.URL.Path, "/api/data/")
    log.Printf("Getting data for ID: %s", id)
    data, err := ds.GetData(context.Background())
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
