# Orchestrix Backend Documentation

This guide explains how to run and use the Go backend 
## What This Service Does

The backend provides:
- user authentication (register/login)
- JWT-protected document APIs
- document upload and status tracking
- query forwarding to the AI service

Base URL (local):
- http://localhost:8080

## Quick Start

## 1. Requirements
- Go 1.22+ (or compatible with this project)
- PostgreSQL database
- AI service running (for processing/query flows)

## 2. Environment Variables
Create or update `backend/.env`:

```env
PORT=8080
DB_URL=postgresql://<user>:<password>@<host>/<db>?sslmode=require
JWT_SECRET=your_jwt_secret
```

## 3. Run the Backend
From the `backend` folder:

```bash
go run ./cmd/api/main.go
```

Health check:

```bash
curl http://localhost:8080/api/health
```

Expected:

```json
{
  "status": "ok",
  "service": "orchestrix-backend"
}
```

## Authentication

Most document routes require a JWT token in the `Authorization` header.

Header format:

```text
Authorization: Bearer <token>
```

## Auth Endpoints

## Register
- Method: `POST`
- Path: `/auth/register`
- Auth: not required

Body:

```json
{
  "name": "Alice Doe",
  "email": "alice@example.com",
  "password": "strong-password"
}
```

Success response (`201`):

```json
{
  "message": "user registered"
}
```

## Login
- Method: `POST`
- Path: `/auth/login`
- Auth: not required

Body:

```json
{
  "email": "alice@example.com",
  "password": "strong-password"
}
```

Success response (`200`):

```json
{
  "token": "<jwt-token>"
}
```

## Document Endpoints

All routes below require JWT auth.

## Upload Document
- Method: `POST`
- Path: `/api/documents/upload`
- Content type: `multipart/form-data`
- Field name: `file`

Example:

```bash
curl -X POST http://localhost:8080/api/documents/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@/absolute/path/to/document.pdf"
```

Success response (`201`):

```json
{
  "message": "document uploaded",
  "document_id": 12,
  "file": "1730000000000_document.pdf",
  "status": "uploaded"
}
```

Note:
- Upload stores the document with status `uploaded`.
- Call `POST /api/documents/:id/process` to start AI processing.
- After processing starts, status transitions: `uploaded -> processing -> ready`.

## List My Documents
- Method: `GET`
- Path: `/api/documents/`

Example:

```bash
curl http://localhost:8080/api/documents/ \
  -H "Authorization: Bearer <token>"
```

## Get One Document By ID
- Method: `GET`
- Path: `/api/documents/:id`

Example:

```bash
curl http://localhost:8080/api/documents/12 \
  -H "Authorization: Bearer <token>"
```

Important behavior:
- This endpoint currently relies on a service that only returns documents in `ready` state.
- If document is still processing, you may receive an error message similar to:
  - `document is not ready for querying`

## Query a Document
- Method: `POST`
- Path: `/api/documents/:id/query`
- Content type: `application/json`

Body:

```json
{
  "question": "Summarize this document"
}
```

Example:

```bash
curl -X POST http://localhost:8080/api/documents/12/query \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"question":"Summarize this document"}'
```

Typical success response (`200`):

```json
{
  "document_id": 12,
  "answer": "...",
  "sources": [
    {
      "document_id": "12",
      "chunk_index": 0,
      "text": "...",
      "distance": 0.123
    }
  ]
}
```

## Manually Process a Document
- Method: `POST`
- Path: `/api/documents/:id/process`

Use this when a document is uploaded but not processed yet.

Example:

```bash
curl -X POST http://localhost:8080/api/documents/12/process \
  -H "Authorization: Bearer <token>"
```

## Document Status Lifecycle

Current statuses used by backend logic:
- `uploaded`: created but not processed
- `processing`: AI pipeline currently running
- `ready`: available for querying
- `failed`: processing failed

## Common Errors and Fixes

## 401 Unauthorized
Symptoms:
- `missing token`
- `invalid token`

Fix:
- send `Authorization: Bearer <token>`
- ensure `JWT_SECRET` is configured and consistent

## 409 / "document is still processing"
Symptoms:
- Query endpoint rejects request

Fix:
- poll `GET /api/documents/` until target document status is `ready`

## 500 From AI Service
Symptoms:
- `failed to query ai service`
- `processing failed`

Fix:
- ensure AI service is running
- verify AI service endpoint URLs/ports
- inspect backend and AI logs

## CORS Errors in Browser
Symptoms:
- frontend request blocked in browser devtools (CORS)

Fix:
- add CORS middleware in backend for your frontend origin
- ensure frontend is calling correct backend URL and route

## Current Integration Notes

The backend currently calls two AI service ports:
- ingest/query: `http://localhost:8000`
- process-document: `http://localhost:9000`

If only one AI service instance is running, align these URLs to the same active host/port.

## Route Map Summary

Public:
- `GET /api/health`
- `POST /auth/register`
- `POST /auth/login`

Protected (`Authorization: Bearer <token>`):
- `POST /api/documents/upload`
- `GET /api/documents/`
- `GET /api/documents/:id`
- `POST /api/documents/:id/query`
- `POST /api/documents/:id/process`
