package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	resp, err := http.Post(
		"http://localhost:8000/ai/ingest",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ai service error: status %d body: %s", resp.StatusCode, string(body))
	}

	return nil
}
