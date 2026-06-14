import client from './client';

export interface Project {
  id: number;
  name: string;
  description: string;
  modules: string[];
  brand_config: Record<string, any>;
  api_key: string;
  status: string;
  created_at: string;
}

export async function getProjects(): Promise<Project[]> {
  const res = await client.get('/factory/projects');
  return res.data.data;
}

export async function createProject(data: { name: string; description: string; modules: string[] }): Promise<Project> {
  const res = await client.post('/factory/projects', data);
  return res.data.data;
}

export async function getProject(id: number): Promise<Project> {
  const res = await client.get(`/factory/projects/${id}`);
  return res.data.data;
}

export async function updateProject(id: number, data: { modules?: string[]; brand_config?: Record<string, any> }): Promise<void> {
  await client.put(`/factory/projects/${id}`, data);
}
