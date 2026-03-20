from fastapi import APIRouter
from pydantic import BaseModel
from app.services.rag_pipeline import answer_query

router = APIRouter()

class QueryRequest(BaseModel):
    document_id: str
    question: str

class SourceChunk(BaseModel):
    document_id: str
    chunk_index: int
    text: str

class QueryResponse(BaseModel):
    answer: str
    sources: list[SourceChunk]

@router.post("/query")
def query(req: QueryRequest):
    answer = answer_query(req.document_id, req.question)
    return {"answer": answer}