from fastapi import APIRouter
from pydantic import BaseModel
from app.services.rag_pipeline import answer_query

router = APIRouter()

class QueryRequest(BaseModel):
    question: str

@router.post("/query")
def query(req: QueryRequest):
    answer = answer_query(req.question)
    return {"answer": answer}