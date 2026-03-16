import faiss
import numpy as np

index = None
documents = []

def store_embeddings(embeddings, chunks):
    global index, documents
    dim = len(embeddings[0])
    if index is None:
        index = faiss.IndexFlatL2(dim)
    index.add(np.array(embeddings).astype("float32"))
    documents.extend(chunks)

def search(query_embedding, k=3):
    if index is None:
        return []
    D, I = index.search(np.array([query_embedding]).astype("float32"), k)
    return [documents[i] for i in I[0]]