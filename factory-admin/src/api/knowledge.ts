import client from './client';

export interface KnowledgeStats {
  total_chunks: number;
  total_documents: number;
  collection_name: string;
}

export interface DocumentItem {
  source: string;
  format: string;
  chunk_count: number;
  indexed_at: string;
}

export interface DocumentListResponse {
  documents: DocumentItem[];
  total: number;
}

export interface UploadResponse {
  source: string;
  chunks: number;
  format: string;
  status: string;
}

export interface DeleteResponse {
  source: string;
  deleted_chunks: number;
}

export async function uploadDocument(file: File): Promise<UploadResponse> {
  const formData = new FormData();
  formData.append('file', file);
  const res = await client.post('/knowledge/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
  return res.data.data ?? res.data;
}

export async function getDocuments(): Promise<DocumentListResponse> {
  const res = await client.get('/knowledge/documents');
  return res.data.data ?? res.data;
}

export async function deleteDocument(source: string): Promise<DeleteResponse> {
  const res = await client.delete(`/knowledge/documents/${encodeURIComponent(source)}`);
  return res.data.data ?? res.data;
}

export async function getKnowledgeStats(): Promise<KnowledgeStats> {
  const res = await client.get('/knowledge/stats');
  return res.data.data ?? res.data;
}
