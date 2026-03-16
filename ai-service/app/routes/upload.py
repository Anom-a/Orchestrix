from fastapi import APIRouter, UploadFile, File
from app.services.document_loader import load_document

router = APIRouter()
@router.post("/upload")
async def upload_file(file: UploadFile = File(...)):
    text = await load_document(file)
    return {"status": "file received", "length": len(text)}