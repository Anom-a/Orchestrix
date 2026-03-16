from fastapi import APIRouter
from pydantic import BaseModel
from app.services.rag_pipeline  import ingest_document

router = APIRouter()
class IngestRequest(BaseModel):
    text: str

@router.post("/ingest")
def ingest(req: IngestRequest):
    ingest_document(req.text)
    return {"status": "document ingested"}
