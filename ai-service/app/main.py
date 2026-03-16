from fastapi import FastAPI
from app.routes import ingest, query

app = FastAPI()

app.include_router(ingest.router, prefix="/ai")
app.include_router(query.router, prefix="/ai")

@app.get("/health")
def health():
    return {
        "status": "ok",
        "service": "orchestrix-ai-service"
        }