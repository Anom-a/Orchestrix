import React, { useState, useEffect } from 'react';
import { api } from '../api';
import { UploadCloud, FileText, CheckCircle, Clock, AlertCircle, Play } from 'lucide-react';
import './Panels.css';
import UploadModal from './UploadModal';

export default function DocumentsPanel() {
  const [documents, setDocuments] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);
  const [error, setError] = useState('');
  
  const fetchDocuments = async () => {
    try {
      const data = await api.getDocuments();
      setDocuments(data || []);
    } catch (err) {
      setError('Failed to fetch documents');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchDocuments();
    // Poll for status updates
    const interval = setInterval(fetchDocuments, 10000);
    return () => clearInterval(interval);
  }, []);

  const handleProcess = async (id) => {
    try {
      await api.processDocument(id);
      fetchDocuments();
    } catch (err) {
      alert(err.message || 'Failed to process document');
    }
  };

  const getStatusBadge = (status) => {
    switch(status) {
      case 'ready':
        return <span className="badge badge-success"><CheckCircle size={14} className="mr-1"/> Ready</span>;
      case 'processing':
        return <span className="badge badge-warning"><Clock size={14} className="mr-1"/> Processing</span>;
      case 'failed':
        return <span className="badge badge-danger"><AlertCircle size={14} className="mr-1"/> Failed</span>;
      default:
        return <span className="badge badge-info"><FileText size={14} className="mr-1"/> Uploaded</span>;
    }
  };

  const stats = {
    total: documents.length,
    ready: documents.filter(d => d.status === 'ready').length,
    processing: documents.filter(d => d.status === 'processing').length
  };

  return (
    <div className="panel fade-in">
      <div className="panel-header">
        <div>
          <h1>Documents</h1>
          <p className="panel-subtitle">Upload, process, and manage your knowledge base</p>
        </div>
        <button className="btn btn-primary" onClick={() => setIsUploadModalOpen(true)}>
          <UploadCloud size={18} />
          <span>Upload Document</span>
        </button>
      </div>

      <div className="stats-grid">
        <div className="stat-card glass-panel">
          <div className="stat-icon bg-primary-light">
            <FileText size={24} className="text-primary" />
          </div>
          <div className="stat-details">
            <span className="stat-value">{stats.total}</span>
            <span className="stat-label">Total Files</span>
          </div>
        </div>
        <div className="stat-card glass-panel">
          <div className="stat-icon bg-success-light">
            <CheckCircle size={24} className="text-success" />
          </div>
          <div className="stat-details">
            <span className="stat-value">{stats.ready}</span>
            <span className="stat-label">Ready</span>
          </div>
        </div>
        <div className="stat-card glass-panel">
          <div className="stat-icon bg-warning-light">
            <Clock size={24} className="text-warning" />
          </div>
          <div className="stat-details">
            <span className="stat-value">{stats.processing}</span>
            <span className="stat-label">Processing</span>
          </div>
        </div>
      </div>

      <div className="table-container glass-panel">
        {isLoading ? (
          <div className="loading-state">
            <div className="spinner"></div>
            <p>Loading documents...</p>
          </div>
        ) : documents.length === 0 ? (
          <div className="empty-state">
            <FileText size={48} className="empty-icon" />
            <h3>No documents found</h3>
            <p>Upload your first document to get started</p>
            <button className="btn btn-primary mt-4" onClick={() => setIsUploadModalOpen(true)}>
              Upload Now
            </button>
          </div>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>File Name</th>
                <th>Status</th>
                <th>Date</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {documents.map((doc) => (
                <tr key={doc.ID}>
                  <td className="font-medium">
                    <div className="flex items-center gap-2">
                      <FileText size={16} className="text-muted" />
                      {doc.filename || `Document #${doc.ID}`}
                    </div>
                  </td>
                  <td>{getStatusBadge(doc.status)}</td>
                  <td className="text-muted">
                    {doc.created_at ? new Date(doc.created_at).toLocaleDateString() : 'N/A'}
                  </td>
                  <td>
                    {(doc.status === 'uploaded' || doc.status === 'failed') && (
                      <button 
                        className="btn btn-secondary btn-sm"
                        onClick={() => handleProcess(doc.ID)}
                        title="Process Document"
                      >
                        <Play size={14} />
                        Process
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {isUploadModalOpen && (
        <UploadModal 
          onClose={() => setIsUploadModalOpen(false)} 
          onSuccess={() => {
            setIsUploadModalOpen(false);
            fetchDocuments();
          }}
        />
      )}
    </div>
  );
}
