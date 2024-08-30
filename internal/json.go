package internal

import (
	"encoding/json"
	"os"
)

func ReadJson[T any](path string, obj *T) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, obj)
	if err != nil {
		return err
	}
	return nil
}
