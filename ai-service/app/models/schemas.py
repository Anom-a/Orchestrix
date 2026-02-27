from pydantic import BaseModel

class ProcessDocumentRequest(BaseModel):
    document_id: str
    file_path: str

class QueryRequest(BaseModel):
    document_id: str
    question: str

class QueryResponse(BaseModel):
    answer: str