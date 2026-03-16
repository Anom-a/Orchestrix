from fastapi import FastAPI
from app.routes import ingest, query
from app.services.vector_store import load_index
app = FastAPI()
@app.on_event("startup")
def startup_event():
    load_index()
app.include_router(ingest.router, prefix="/ai")
app.include_router(query.router, prefix="/ai")

@app.get("/health")
def health():
    return {
        "status": "ok",
        "service": "orchestrix-ai-service"
        }