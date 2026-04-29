import React, { useState } from 'react';
import { FileText, MessageSquare, LogOut, Menu, X, Database } from 'lucide-react';
import DocumentsPanel from '../components/DocumentsPanel';
import ChatPanel from '../components/ChatPanel';
import './Dashboard.css';

export default function Dashboard({ setIsAuthenticated }) {
  const [activeTab, setActiveTab] = useState('documents');
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const handleLogout = () => {
    localStorage.removeItem('token');
    setIsAuthenticated(false);
  };

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <div className="app-container">
      {/* Mobile Header */}
      <div className="mobile-header">
        <div className="mobile-brand">
          <Database size={24} className="text-primary" />
          <span>Orchestrix</span>
        </div>
        <button className="menu-btn" onClick={toggleSidebar}>
          {isSidebarOpen ? <X size={24} /> : <Menu size={24} />}
        </button>
      </div>

      {/* Sidebar Overlay for Mobile */}
      {isSidebarOpen && <div className="sidebar-overlay" onClick={() => setIsSidebarOpen(false)}></div>}

      {/* Sidebar */}
      <aside className={`sidebar ${isSidebarOpen ? 'open' : ''}`}>
        <div className="sidebar-header">
          <Database size={28} className="sidebar-logo-icon" />
          <span className="sidebar-logo-text">Orchestrix</span>
        </div>

        <nav className="sidebar-nav">
          <button 
            className={`nav-item ${activeTab === 'documents' ? 'active' : ''}`}
            onClick={() => { setActiveTab('documents'); setIsSidebarOpen(false); }}
          >
            <FileText size={20} />
            <span>Documents</span>
          </button>
          
          <button 
            className={`nav-item ${activeTab === 'chat' ? 'active' : ''}`}
            onClick={() => { setActiveTab('chat'); setIsSidebarOpen(false); }}
          >
            <MessageSquare size={20} />
            <span>Chat</span>
          </button>
        </nav>

        <div className="sidebar-footer">
          <button className="nav-item logout-btn" onClick={handleLogout}>
            <LogOut size={20} />
            <span>Logout</span>
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="main-content">
        {activeTab === 'documents' && <DocumentsPanel />}
        {activeTab === 'chat' && <ChatPanel />}
      </main>
    </div>
  );
}
