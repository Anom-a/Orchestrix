package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ProccessRequest struct{
	DocumentID string `json:"document_id"`
	FilePath string `json:"file_path"`
}
func SendToAIService(docID, filepath string) error{
	reqBody := ProccessRequest{
		DocumentID: docID,
		FilePath: filepath,
	}
	jsonData, _ := json.Marshal(reqBody)
	_, err := http.Post(
		"http://localhost:9000/ai/process-document",
		"applicaton/json",
		bytes.NewBuffer(jsonData),
	)
	return err
}