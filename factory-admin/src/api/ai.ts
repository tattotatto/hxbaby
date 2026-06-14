import client from './client';

interface AIGenerateResponse {
  content: string;
  tokens_used: number;
  model: string;
}

export async function generateArticle(topic: string, category?: string): Promise<AIGenerateResponse> {
  const res = await client.post('/ai/generate-article', { topic, category });
  return res.data.data;
}

export async function generateSummary(content: string): Promise<AIGenerateResponse> {
  const res = await client.post('/ai/generate-summary', { content });
  return res.data.data;
}

export async function generateActivityCopy(title: string, description?: string): Promise<AIGenerateResponse> {
  const res = await client.post('/ai/generate-activity-copy', { title, description });
  return res.data.data;
}

export async function generateSellingPoints(name: string, description?: string): Promise<AIGenerateResponse> {
  const res = await client.post('/ai/generate-selling-points', { name, description });
  return res.data.data;
}
