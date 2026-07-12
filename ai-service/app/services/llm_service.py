import os
import logging
from langchain_google_genai import ChatGoogleGenerativeAI
from langchain_core.messages import HumanMessage, SystemMessage
from app.core.config import GEMINI_MODEL

logger = logging.getLogger(__name__)

def get_llm():
    api_key = os.getenv("GOOGLE_API_KEY") or os.getenv("GEMINI_API_KEY")
    if not api_key:
        raise RuntimeError("Gemini API key not configured")

    return ChatGoogleGenerativeAI(
        model=GEMINI_MODEL,
        temperature=0,
        api_key=api_key,
    )


def _fallback_answer(context: str) -> str:
    cleaned = (context or "").strip()
    if not cleaned:
        return "The document does not contain enough information to answer this question."

    lines = [line.strip() for line in cleaned.splitlines() if line.strip()]
    if not lines:
        return "The document does not contain enough information to answer this question."

    return f"Based on the document context: {lines[0]}"

def generate_answer(context: str, question: str) -> str:
    try:
        llm = get_llm()

        system_prompt = (
            "You are a document question-answering assistant. "
            "You are given extracted chunks from a document and must answer the user's question based ONLY on these chunks. "
            "When answering, consider the type of document (e.g., resume/CV, report, article, contract) based on the content structure. "
            "For example, if the chunks contain a person's name, skills, experience, and education, the document is likely a resume. "
            "Provide detailed, comprehensive answers. "
            "If the answer is not in the context, say that the document does not contain enough information."
        )

        user_prompt = f"""Document chunks:
{context}

Question: {question}"""

        response = llm.invoke([
            SystemMessage(content=system_prompt),
            HumanMessage(content=user_prompt)
        ])

        return response.content
    except Exception as e:
        logger.error("Gemini LLM call failed: %s", e, exc_info=True)
        return _fallback_answer(context)