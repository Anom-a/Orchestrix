from fastapi import APIRouter
from app.models.schemas import ProcessDocumentRequest, QueryRequest, QueryResponse
from app.services.document_processor import process_document
from app.services.qa_chain import answer_question

router = APIRouter()

@router.get("/health")
def health_check():
    return {"status": "ok", "service": "orchestrix-ai"}

@router.post("/ai/process-document")
def process_doc(req: ProcessDocumentRequest):
    result = process_document(req.document_id, req.file_path)
    return result

@router.post("/ai/query", response_model=QueryResponse)
def query_ai(req: QueryRequest):
    answer = answer_question(req.document_id, req.question)
    return {"answer": answer}