from contextlib import asynccontextmanager
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.core.config import *  # noqa: F401,F403 — loads .env via dotenv
from app.core.config import HOST, PORT, CORS_ORIGINS
from app.routes import ingest, query
from app.services.vector_store import load_index


@asynccontextmanager
async def lifespan(app: FastAPI):
    load_index()
    yield

app = FastAPI(lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(ingest.router, prefix="/ai")
app.include_router(query.router, prefix="/ai")

@app.get("/health")
def health():
    return {
        "status": "ok",
        "service": "orchestrix-ai-service"
        }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("app.main:app", host=HOST, port=PORT, reload=False)