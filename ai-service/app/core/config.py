import os
from dotenv import load_dotenv

load_dotenv()

AI_SERVICE_NAME = "orchestrix-ai"
HOST = os.getenv("HOST", "0.0.0.0")
PORT = int(os.getenv("PORT", "8000"))
CORS_ORIGINS = os.getenv("CORS_ORIGINS", "*").split(",")
STORAGE_PATH = os.getenv("STORAGE_PATH", "../backend")
GEMINI_MODEL = os.getenv("GEMINI_MODEL", "gemini-2.5-flash")