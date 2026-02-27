from fastapi import FastAPI
from app.api.routes import router

app = FastAPI(title="Orchestrix AI Service")

app.include_router(router)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("app.main:app", host="0.0.0.0", port=9000, reload=True)