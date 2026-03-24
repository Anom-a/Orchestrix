package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AIQueryRequest struct {
	DocumentID string `json:"document_id"`
	Question   string `json:"question"`
}

type AISourceChunk struct {
	DocumentID string  `json:"document_id"`
	ChunkIndex int     `json:"chunk_index"`
	Text       string  `json:"text"`
	Distance   float64 `json:"distance"`
}

type AIQueryResponse struct {
	Answer  string          `json:"answer"`
	Sources []AISourceChunk `json:"sources"`
}
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

func QueryAIService(documentID uint, question string) (*AIQueryResponse, error) {
	reqBody := AIQueryRequest{
		DocumentID: fmt.Sprintf("%d", documentID),
		Question:   question,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		"http://localhost:8000/ai/query",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ai service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var aiResp AIQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, err
	}

	return &aiResp, nil
}

type AIProcessRequest struct {
	DocumentID string `json:"document_id"`
	FilePath   string `json:"file_path"`
}
func ProcessDocumentAI(documentID uint, filePath string) error {
	reqBody := AIProcessRequest{
		DocumentID: fmt.Sprintf("%d", documentID),
		FilePath:   filePath,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Post(
		"http://localhost:9000/ai/process-document",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("AI service failed: %s", string(body))
	}

	return nil
}