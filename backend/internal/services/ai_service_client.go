package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// getAIServiceURL returns the base URL for the AI service from environment,
// falling back to http://localhost:8000 for local development.
func getAIServiceURL() string {
	url := os.Getenv("AI_SERVICE_URL")
	if url == "" {
		url = "http://localhost:8000"
	}
	return strings.TrimRight(url, "/")
}

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
		getAIServiceURL()+"/ai/ingest",
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
		getAIServiceURL()+"/ai/query",
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
		getAIServiceURL()+"/ai/ingest",
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