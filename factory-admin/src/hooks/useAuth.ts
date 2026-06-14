import { useState, useEffect, useCallback } from 'react';
import type { Customer } from '../api/auth';

export function useAuth() {
  const [customer, setCustomer] = useState<Customer | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const saved = localStorage.getItem('factory_customer');
    if (saved) {
      try { setCustomer(JSON.parse(saved) as Customer); } catch { /* ignore parse errors */ }
    }
    setLoading(false);
  }, []);

  const login = useCallback((customer: Customer, token: string) => {
    localStorage.setItem('factory_token', token);
    localStorage.setItem('factory_customer', JSON.stringify(customer));
    setCustomer(customer);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem('factory_token');
    localStorage.removeItem('factory_customer');
    setCustomer(null);
  }, []);

  return { customer, loading, login, logout, isLoggedIn: !!customer };
}
