package storage

import (
	"context"
	"encoding/json"
	"os"
	"slackwatch/backend/internal/model"
)

type FileStorage struct {
	FilePath string
}

func (fs *FileStorage) Save(ctx context.Context, data model.Data) error {
	file, err := os.OpenFile(fs.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func (fs *FileStorage) Get(ctx context.Context, id string) (model.Data, error) {
	// Implementation for retrieving data by ID from file
	// This is a simplified example. A real implementation might need to read the entire file and decode each entry to find the matching ID.
	return model.Data{}, nil
}
