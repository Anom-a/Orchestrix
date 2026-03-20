import sys
import types

from fastapi import FastAPI
from fastapi.testclient import TestClient


fake_rag_pipeline = types.ModuleType("app.services.rag_pipeline")
fake_rag_pipeline.answer_query = lambda _document_id, _question: ""
fake_rag_pipeline.ingest_document = lambda _document_id, _text: None
sys.modules["app.services.rag_pipeline"] = fake_rag_pipeline

from app.routes import ingest, query


def _build_app() -> FastAPI:
    app = FastAPI()
    app.include_router(ingest.router, prefix="/ai")
    app.include_router(query.router, prefix="/ai")
    return app


def test_ingest_success(monkeypatch):
    app = _build_app()
    client = TestClient(app)

    calls = {}

    def fake_load_document_from_path(path: str):
        calls["path"] = path
        return "example text"

    def fake_ingest_document(document_id: str, text: str):
        calls["document_id"] = document_id
        calls["text"] = text

    monkeypatch.setattr(ingest, "load_document_from_path", fake_load_document_from_path)
    monkeypatch.setattr(ingest, "ingest_document", fake_ingest_document)

    response = client.post(
        "/ai/ingest",
        json={"document_id": "doc-1", "file_path": "storage/uploads/file.txt"},
    )

    assert response.status_code == 200
    assert response.json() == {"status": "document ingested", "document_id": "doc-1"}
    assert calls["path"] == "../backend/storage/uploads/file.txt"
    assert calls["document_id"] == "doc-1"
    assert calls["text"] == "example text"


def test_ingest_missing_file_returns_404(monkeypatch):
    app = _build_app()
    client = TestClient(app)

    def fake_load_document_from_path(path: str):
        raise FileNotFoundError(path)

    monkeypatch.setattr(ingest, "load_document_from_path", fake_load_document_from_path)

    response = client.post(
        "/ai/ingest",
        json={"document_id": "doc-2", "file_path": "storage/uploads/missing.txt"},
    )

    assert response.status_code == 404


def test_query_success(monkeypatch):
    app = _build_app()
    client = TestClient(app)

    def fake_answer_query(document_id: str, question: str):
        assert document_id == "doc-9"
        assert question == "What is this about?"
        return "This is a test answer."

    monkeypatch.setattr(query, "answer_query", fake_answer_query)

    response = client.post(
        "/ai/query",
        json={"document_id": "doc-9", "question": "What is this about?"},
    )

    assert response.status_code == 200
    assert response.json() == {"answer": "This is a test answer."}
