import React, { useState, useRef } from 'react';
import { api } from '../api';
import { X, UploadCloud, File } from 'lucide-react';
import './UploadModal.css';

export default function UploadModal({ onClose, onSuccess }) {
  const [file, setFile] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const fileInputRef = useRef(null);

  const handleDrop = (e) => {
    e.preventDefault();
    const droppedFile = e.dataTransfer.files[0];
    validateAndSetFile(droppedFile);
  };

  const handleDragOver = (e) => {
    e.preventDefault();
  };

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    validateAndSetFile(selectedFile);
  };

  const validateAndSetFile = (selectedFile) => {
    if (!selectedFile) return;
    
    const validTypes = ['application/pdf', 'text/plain', 'text/markdown'];
    const validExtensions = ['.pdf', '.txt', '.md'];
    const extension = selectedFile.name.substring(selectedFile.name.lastIndexOf('.')).toLowerCase();
    
    if (validTypes.includes(selectedFile.type) || validExtensions.includes(extension)) {
      setFile(selectedFile);
      setError('');
    } else {
      setError('Please upload a PDF, TXT, or MD file.');
      setFile(null);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!file) {
      setError('Please select a file to upload');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      await api.uploadDocument(file);
      onSuccess();
    } catch (err) {
      setError(err.message || 'Failed to upload document');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal-card fade-in">
        <div className="modal-header">
          <h2>Upload Document</h2>
          <button className="close-btn" onClick={onClose}><X size={20} /></button>
        </div>
        
        <form onSubmit={handleSubmit} className="modal-body">
          <div 
            className={`dropzone ${file ? 'has-file' : ''}`}
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onClick={() => fileInputRef.current.click()}
          >
            <input 
              type="file" 
              ref={fileInputRef} 
              onChange={handleFileChange} 
              accept=".pdf,.txt,.md" 
              className="hidden-input" 
            />
            
            {!file ? (
              <div className="dropzone-content">
                <div className="dropzone-icon">
                  <UploadCloud size={40} />
                </div>
                <p className="dropzone-title">Click to upload or drag and drop</p>
                <p className="dropzone-hint">PDF, TXT, or MD files supported</p>
              </div>
            ) : (
              <div className="file-preview">
                <File size={32} className="text-primary" />
                <div className="file-info">
                  <span className="file-name">{file.name}</span>
                  <span className="file-size">{(file.size / 1024 / 1024).toFixed(2)} MB</span>
                </div>
                <button 
                  type="button" 
                  className="remove-file-btn" 
                  onClick={(e) => { e.stopPropagation(); setFile(null); }}
                >
                  <X size={16} />
                </button>
              </div>
            )}
          </div>
          
          {error && <div className="form-error mt-4">{error}</div>}
          
          <div className="modal-footer mt-6">
            <button type="button" className="btn btn-secondary" onClick={onClose}>Cancel</button>
            <button type="submit" className="btn btn-primary" disabled={!file || isLoading}>
              {isLoading ? <div className="spinner"></div> : 'Upload File'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
