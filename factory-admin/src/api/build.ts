import client from './client';

export interface BuildTask {
  id: number;
  project_id: number;
  status: string;
  output_zip_url: string;
  output_md5: string;
  error_log: string;
  duration_ms: number;
  created_at: string;
  completed_at: string | null;
}

export async function triggerBuild(projectId: number): Promise<{ build_id: number; status: string }> {
  const res = await client.post(`/factory/projects/${projectId}/build`);
  return res.data.data;
}

export async function getBuildStatus(id: number): Promise<BuildTask> {
  const res = await client.get(`/builds/${id}/status`);
  return res.data.data;
}

export async function getBuildHistory(projectId: number): Promise<BuildTask[]> {
  const res = await client.get(`/factory/projects/${projectId}/builds`);
  return res.data.data;
}

export function getDownloadUrl(id: number): string {
  return `${client.defaults.baseURL}/builds/${id}/download`;
}
