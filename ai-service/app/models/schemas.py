from pydantic import BaseModel
from typing import List

class QueryRequest(BaseModel):
    document_id: str
    question: str

class SourceChunk(BaseModel):
    document_id: str
    chunk_index: int
    text: str

class QueryResponse(BaseModel):
    answer: str
    sources: List[SourceChunk]