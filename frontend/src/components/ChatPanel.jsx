import React, { useState, useEffect, useRef } from 'react';
import { api } from '../api';
import { Send, FileText, Database, Bot, User } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import './Panels.css';

export default function ChatPanel() {
  const [documents, setDocuments] = useState([]);
  const [selectedDoc, setSelectedDoc] = useState(null);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef(null);

  useEffect(() => {
    fetchReadyDocuments();
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const fetchReadyDocuments = async () => {
    try {
      const data = await api.getDocuments();
      setDocuments(data?.filter(d => d.status === 'ready') || []);
    } catch (err) {
      console.error('Failed to fetch documents', err);
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleDocSelect = (doc) => {
    setSelectedDoc(doc);
    setMessages([{
      id: 'welcome',
      role: 'bot',
      content: `I'm ready to answer questions about **${doc.filename}**. What would you like to know?`
    }]);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!input.trim() || !selectedDoc || isLoading) return;

    const userMessage = { id: Date.now().toString(), role: 'user', content: input };
    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);

    try {
      const res = await api.queryDocument(selectedDoc.ID, userMessage.content);
      
      const botMessage = {
        id: (Date.now() + 1).toString(),
        role: 'bot',
        content: res.answer,
        sources: res.sources
      };
      
      setMessages(prev => [...prev, botMessage]);
    } catch (err) {
      setMessages(prev => [...prev, {
        id: (Date.now() + 1).toString(),
        role: 'bot',
        content: `Error: ${err.message || 'Failed to get answer'}`,
        isError: true
      }]);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="chat-layout fade-in">
      {/* Sidebar for document selection */}
      <div className="chat-sidebar glass-panel">
        <h3 className="chat-sidebar-title">Select Document</h3>
        <div className="chat-doc-list">
          {documents.length === 0 ? (
            <p className="text-muted p-4 text-center">No ready documents available. Please upload and process a document first.</p>
          ) : (
            documents.map(doc => (
              <button
                key={doc.ID}
                className={`chat-doc-item ${selectedDoc?.ID === doc.ID ? 'active' : ''}`}
                onClick={() => handleDocSelect(doc)}
              >
                <FileText size={16} />
                <span className="truncate">{doc.filename}</span>
              </button>
            ))
          )}
        </div>
      </div>

      {/* Main chat area */}
      <div className="chat-area glass-panel">
        {!selectedDoc ? (
          <div className="chat-empty-state">
            <Database size={48} className="empty-icon" />
            <h2>Select a document to chat</h2>
            <p>Choose a document from the left sidebar to start asking questions.</p>
          </div>
        ) : (
          <>
            <div className="chat-header">
              <h2>{selectedDoc.filename}</h2>
              <span className="badge badge-success">Ready</span>
            </div>
            
            <div className="chat-messages-container">
              {messages.map((msg) => (
                <div key={msg.id} className={`chat-message ${msg.role}`}>
                  <div className="message-avatar">
                    {msg.role === 'user' ? <User size={18} /> : <Bot size={18} />}
                  </div>
                  <div className={`message-bubble ${msg.isError ? 'error' : ''}`}>
                    <div className="message-content">
                      {msg.role === 'bot' ? (
                        <ReactMarkdown>{msg.content}</ReactMarkdown>
                      ) : (
                        msg.content
                      )}
                    </div>
                    
                    {/* Sources section if available */}
                    {msg.sources && msg.sources.length > 0 && (
                      <div className="message-sources">
                        <div className="sources-title">Sources used:</div>
                        <ul className="sources-list">
                          {msg.sources.map((src, i) => (
                            <li key={i} className="source-item truncate" title={src.text}>
                              "{src.text.substring(0, 60)}..."
                            </li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                </div>
              ))}
              {isLoading && (
                <div className="chat-message bot">
                  <div className="message-avatar"><Bot size={18} /></div>
                  <div className="message-bubble loading">
                    <div className="typing-indicator">
                      <span></span><span></span><span></span>
                    </div>
                  </div>
                </div>
              )}
              <div ref={messagesEndRef} />
            </div>

            <form onSubmit={handleSubmit} className="chat-input-container">
              <input
                type="text"
                className="chat-input"
                placeholder="Ask a question about this document..."
                value={input}
                onChange={(e) => setInput(e.target.value)}
                disabled={isLoading}
              />
              <button 
                type="submit" 
                className="chat-send-btn"
                disabled={!input.trim() || isLoading}
              >
                <Send size={18} />
              </button>
            </form>
          </>
        )}
      </div>
    </div>
  );
}
