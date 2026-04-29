from app.services.embedding_service import chunk_text, embed_chunks
from app.services.vector_store import store_embeddings, search
from sentence_transformers import SentenceTransformer
from app.services.llm_service import generate_answer
model = SentenceTransformer("all-MiniLM-L6-v2")

def ingest_document(document_id: str, text: str):
    chunks = chunk_text(text)
    embeddings = embed_chunks(chunks)
    store_embeddings(document_id, embeddings, chunks)

def answer_query(document_id,  query: str):
    query_embedding = model.encode(query)
    retrieved_chunks = search(query_embedding, document_id)
    if not retrieved_chunks:
        return {
            "answer": "No relevant information found for this document.",
            "sources": []
        }
    best_distance = retrieved_chunks[0]["distance"]
    # FAISS L2: lower distance = better match. Reject if best match is too far.
    if best_distance > 1.5:
        return {
            "answer": "The document does not contain enough information to answer that",
            "sources": []
        }
    context = "\n".join([chunk["text"] for chunk in retrieved_chunks])
    answer = generate_answer(context, query)
    return {
        "answer": answer,
        "sources": retrieved_chunks
    }