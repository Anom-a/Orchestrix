from pypdf import PdfReader
from fastapi import UploadFile
import os

async def load_document(file: UploadFile):
    if file.filename.endswith(".pdf"):
        reader = PdfReader(file.file)
        text = ""
        for page in reader.pages:
            text += page.extract_text()
        return text
    else:
        content = await file.read()
        return content.decode("utf-8")

def load_document_from_path(file_path: str):
    if not os.path.exists(file_path):
        raise FileNotFoundError(f"File not found at {file_path}")
        
    if file_path.endswith(".pdf"):
        reader = PdfReader(file_path)
        text = ""
        for page in reader.pages:
            text += page.extract_text()
        return text
    else:
        with open(file_path, "r", encoding="utf-8") as f:
            return f.read()
    
    