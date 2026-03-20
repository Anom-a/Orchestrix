from app.services.embedding_service import chunk_text, embed_chunks
from app.services.vector_store import store_embeddings, search
from sentence_transformers import SentenceTransformer
model = SentenceTransformer("all-MiniLM-L6-v2")

def ingest_document(text: str):
    chunks = chunk_text(text)
    embeddings = embed_chunks(chunks)
    store_embeddings(embeddings, chunks)

def answer_query(document_id,  query: str):
    query_embedding = model.encode(query)
    retrieved_chunks = search(query_embedding, document_id)
    if not retrieved_chunks:
        return {
            "answer": "No relevant information found for this document.",
            "sources": []
        }
    context = "\n".join([chunk["text"] for chunk in retrieved_chunks])
    answer = f"Answer based on retrieved context; \n{context}"
    return {
        "answer": answer,
        "sources": retrieved_chunks
    }