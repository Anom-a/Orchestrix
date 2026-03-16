from app.services.embedding_service import chunk_text, embed_chunks
from app.services.vector_store import store_embeddings, search
from sentence_transformers import SentenceTransformer

model = SentenceTransformer("all-MiniLM-L6-v2")

def ingest_document(text: str):
    chunks = chunk_text(text)
    embeddings = embed_chunks(chunks)
    store_embeddings(embeddings, chunks)

def answer_query(query: str):
    query_embedding = model.encode(query)
    docs = search(query_embedding)
    context = "\n".join(docs)
    return f"Answer based on context:\n{context}"