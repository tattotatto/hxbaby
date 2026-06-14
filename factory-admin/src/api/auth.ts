import client from './client';

export interface Customer {
  id: number;
  name: string;
  phone: string;
  company_name?: string;
  plan: string;
  max_projects: number;
}

export async function login(phone: string, password: string): Promise<{ customer: Customer; token: string }> {
  const res = await client.post('/factory/auth/login', { phone, password });
  return res.data.data;
}

export async function register(phone: string, password: string, name: string): Promise<{ customer: Customer; token: string }> {
  const res = await client.post('/factory/auth/register', { phone, password, name });
  return res.data.data;
}
