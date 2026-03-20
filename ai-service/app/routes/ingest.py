from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from app.services.rag_pipeline  import ingest_document
from app.services.document_loader import load_document_from_path
import os

router = APIRouter()
class IngestRequest(BaseModel):
    document_id: str
    file_path: str

@router.post("/ingest")
def ingest(req: IngestRequest):
    # Construct paths relative to the project root
    # Backend storage is in ../backend/storage/uploads/
    # The file_path from the backend is like "storage/uploads/unique_name"
    # So the full path from ai-service's perspective is:
    full_path = os.path.join("../backend", req.file_path)
    
    try:
        text = load_document_from_path(full_path)
        ingest_document(req.document_id, text)
        return {"status": "document ingested", "document_id": req.document_id}
    except FileNotFoundError as e:
        raise HTTPException(status_code=404, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
