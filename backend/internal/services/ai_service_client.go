package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ProcessRequest struct {
	DocumentID string `json:"document_id"`
	FilePath   string `json:"file_path"`
}

func SendToAIService(docID, filepath string) error {
	reqBody := ProcessRequest{
		DocumentID: docID,
		FilePath:   filepath,
	}
	jsonData, _ := json.Marshal(reqBody)
	_, err := http.Post(
		"http://localhost:8000/ai/ingest",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	return err
}
