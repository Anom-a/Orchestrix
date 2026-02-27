# 🚀 Orchestrix

<div align="center">
  <img src="https://img.shields.io/badge/Status-Active-success?style=for-the-badge" alt="Status" />
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/LangChain-121212?style=for-the-badge&logo=chainlink&logoColor=white" alt="LangChain" />
  <img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License" />
</div>

<br/>

> **“Your personal AI knowledge base server.”**

**Orchestrix** is a robust system that allows users to upload their personal or business documents and intelligently chat with them using AI. It handles the entire lifecycle of document ingestion, vector embedding, and Retrieval-Augmented Generation (RAG) to provide highly accurate, context-aware answers.


## 📑 Table of Contents
- [✨ Key Features](#-key-features)
-[🏗️ High-Level Architecture](#️-high-level-architecture)
- [🧠 AI Processing Pipeline](#-ai-processing-pipeline)
- [🗄️ Database Schema](#️-database-schema)
- [🔌 API Reference](#-api-reference)
- [🚀 Getting Started](#-getting-started)
- [🔮 Future Scope](#-future-scope)

---

## ✨ Key Features

### 👤 User Management
- **Authentication**: Secure JWT-based authentication for registration and login.
- **Profiles**: Manage user accounts and personal settings.

### 📁 Document Management
- **Uploads**: Support for diverse file formats (`.pdf`, `.txt`, `.md`).
- **Processing**: Automatic metadata extraction, text chunking, and preprocessing.
- **Status Tracking**: Monitor document pipelines (`processing`, `ready`, `failed`).

### 🧠 AI & Vector Processing
- **Embeddings**: Automatic generation of high-quality vector embeddings.
- **Vector Storage**: Integration with robust Vector databases (FAISS / Chroma).
- **RAG Engine**: Retrieves highly relevant document chunks to inject context into LLM prompts.

### 💬 Conversational System
- **Sessions**: Maintain continuous chat sessions.
- **History**: Persistent storage of chat history and multiple simultaneous conversations per user.

### ⚙️ System Reliability
- Structured logging, global error handling, rate limiting, background workers, and health checks.

---

## 🏗️ High-Level Architecture

Orchestrix utilizes a modular microservice-style architecture behind a unified API Gateway built with Go's **Gin** framework.

```mermaid
graph TD;
    Client([Client / Postman / Frontend]) -->|HTTP / REST| API[API Gateway Gin]
    
    subgraph Orchestrix Backend
        API --> Auth[Auth Service]
        API --> Doc[Document Service]
        API --> Chat[Chat Service]
        API --> AI[AI Service LangChain]
    end

    subgraph Data Layer
        Auth --> DB[(PostgreSQL)]
        Doc --> DB
        Chat --> DB
        AI --> VDB[(FAISS / Chroma Vector DB)]
    end

    subgraph External Services
        AI -->|API Call| LLM((LLM Provider))
    end