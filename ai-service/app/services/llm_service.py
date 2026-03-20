import os
from langchain_google_genai import ChatGoogleGenerativeAI
from langchain_core.messages import HumanMessage, SystemMessage

def get_llm():
    return ChatGoogleGenerativeAI(
        model="gemini-2.5-flash",
        temperature=0
    )

def generate_answer(context: str, question: str) -> str:
    llm = get_llm()

    system_prompt = (
        "You are a document question-answering assistant. "
        "Answer only using the provided context. "
        "If the answer is not in the context, say that the document does not contain enough information."
    )

    user_prompt = f"""
Context:
{context}

Question:
{question}
"""

    response = llm.invoke([
        SystemMessage(content=system_prompt),
        HumanMessage(content=user_prompt)
    ])

    return response.content