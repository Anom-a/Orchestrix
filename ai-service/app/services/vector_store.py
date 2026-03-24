import faiss
import numpy as np
import os
import pickle
INDEX_PATH = "vector.index"
CHUNK_STORE_PATH = "documents.pk1"
index = None

chunk_store = []

def load_index():
    global index, chunk_store
    if os.path.exists(INDEX_PATH):
        index = faiss.read_index(INDEX_PATH)
    else:
        index = None
    if os.path.exists(CHUNK_STORE_PATH):
        with open(CHUNK_STORE_PATH, "rb") as f:
            chunk_store = pickle.load(f)
    else:
        chunk_store = []
def save_index():
    global index, chunk_store
    if index is not None:
        faiss.write_index(index, INDEX_PATH)
    with open(CHUNK_STORE_PATH, "wb") as f:
        pickle.dump(chunk_store, f)
def store_embeddings(document_id, embeddings, chunks):
    global index, chunk_store

    embeddings = np.array(embeddings).astype("float32")
    if len(embeddings.shape) != 2:
        raise ValueError("Embeddings must be 2D aray")
    dim = embeddings.shape[1]
    if index is None:
        index = faiss.IndexFlatL2(dim)
    index.add(embeddings)
    for i, chunk in enumerate(chunks):
        chunk_store.append({
            "document_id": document_id,
            "chunk_index": 1, 
            "text": chunk
        })
    save_index()



def search(query_embedding, document_id, k=3, fetch_k=10):
    global index, chunk_store

    if index is None or len(chunk_store) == 0:
        return []

    document_id = str(document_id)
    query_embedding = np.array([query_embedding]).astype("float32")
    distances, indices = index.search(query_embedding, fetch_k)

    results = []
    for rank, i in enumerate(indices[0]):
        if 0 <= i < len(chunk_store):
            chunk = chunk_store[i]
            if chunk["document_id"] == document_id:
                results.append({
                    "document_id": chunk["document_id"],
                    "chunk_index": chunk["chunk_index"],
                    "text": chunk["text"],
                    "distance": float(distances[0][rank]),
                })
        if len(results) == k:
            break

    return results