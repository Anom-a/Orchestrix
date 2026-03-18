import faiss
import numpy as np
import os
import pickle
INDEX_PATH = "vector.index"
DOC_PATH = "documents.pk1"
index = None
documents = []

def load_index():
    global index, documents
    if os.path.exists(INDEX_PATH):
        index = faiss.read_index(INDEX_PATH)
    if os.path.exists(DOC_PATH):
        with open(DOC_PATH, "rb") as f:
            documents = pickle.load(f)
def save_index():
    global index, documents
    if index is not None:
        faiss.write_index(index, INDEX_PATH)
    with open(DOC_PATH, "wb") as f:
        pickle.dump(documents, f)
def store_embeddings(embeddings, chunks):
    global index, documents

    embeddings = np.array(embeddings).astype("float32")

    if index is None:
        dim = embeddings.shape[1]
        index = faiss.IndexFlatL2(dim)

    index.add(embeddings)
    documents.extend(chunks)

    save_index()

def search(query_embedding, k=3):
    if index is None:
        return []
    D, I = index.search(np.array([query_embedding]).astype("float32"), k)
    return [documents[i] for i in I[0]]