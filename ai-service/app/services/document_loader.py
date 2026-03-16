from pypdf import pdfReader
from fastapi import UploadFile
async def load_document(file: UploadFile):
    if file.filename.endswith(".pdf"):
        reader  = pdfReader(file.file)
        text = ""
        for page in reader:
            text += page.extract_text()
        return text
    else:
        content = await file.read()
        return content.decode("utf-8")
    
    