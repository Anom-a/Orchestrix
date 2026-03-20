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
def store_embeddings(embeddings, chunks):
    global index, chunk_store

    embeddings = np.array(embeddings).astype("float32")

    if index is None:
        dim = embeddings.shape[1]
        index = faiss.IndexFlatL2(dim)

    index.add(embeddings)
    chunk_store.extend(chunks)

    save_index()

def search(query_embedding, k=3):
    if index is None:
        return []
    D, I = index.search(np.array([query_embedding]).astype("float32"), k)
    return [chunk_store[i] for i in I[0]]