package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type DataService struct {
	// Storage field is omitted for brevity as it's not directly related to the changes
}

func (ds *DataService) SaveContainerData(ctx context.Context, data interface{}) error {
	return ds.saveDataToFile(ctx, data, "containers.json")
}

func (ds *DataService) SaveImageData(ctx context.Context, data interface{}) error {
	return ds.saveDataToFile(ctx, data, "imageUpdates.json")
}

func (ds *DataService) saveDataToFile(ctx context.Context, data interface{}, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

// GetData now accepts a fileName parameter to specify which file to read from
func (ds *DataService) GetData(ctx context.Context, fileName string) ([]map[string]string, error) {
	file, err := os.Open(fileName)
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

func (ds *DataService) GetCombinedData(ctx context.Context) ([]map[string]interface{}, error) {
	containerData, err := ds.GetData(ctx, "containers.json")
	if err != nil {
		return nil, err
	}
	imageData, err := ds.GetData(ctx, "imageUpdates.json")
	if err != nil {
		return nil, err
	}
	combinedData := make([]map[string]interface{}, 0)
	imageMap := make(map[string]map[string]string)
	for _, image := range imageData {
		imageMap[image["containerName"]] = image
	}
	for _, container := range containerData {
		combined := make(map[string]interface{})
		for k, v := range container {
			combined[k] = v
		}
		if update, exists := imageMap[container["containerName"]]; exists {
			combined["updateAvailable"] = true
			for k, v := range update {
				combined[k] = v
			}
		} else {
			combined["updateAvailable"] = false
		}
		combinedData = append(combinedData, combined)
	}
	return combinedData, nil
}

func (ds *DataService) SaveData(ctx context.Context, dataType string, data interface{}) error {
	switch dataType {
	case "container":
		return ds.SaveContainerData(ctx, data)
	case "image":
		return ds.SaveImageData(ctx, data)
	default:
		return fmt.Errorf("unsupported data type: %s", dataType)
	}
}