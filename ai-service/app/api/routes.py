from fastapi import APIRouter
from app.models.schemas import ProcessDocumentRequest, QueryRequest, QueryResponse
from app.services.rag_pipeline import answer_query, ingest_document

router = APIRouter()

@router.get("/health")
def health_check():
    return {"status": "ok", "service": "orchestrix-ai"}

@router.post("/ai/process-document")
def process_doc(req: ProcessDocumentRequest):
    result = ingest_document(req.document_id, req.file_path)
    return result

@router.post("/ai/query", response_model=QueryResponse)
def query_ai(req: QueryRequest):
    answer = answer_query(req.document_id, req.question)
    return {"answer": answer}