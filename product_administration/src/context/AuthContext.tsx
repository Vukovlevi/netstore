import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';

const AUTH_API_PATH = './api/auth';
const LOGIN_PATH = `${import.meta.env.VITE_STORE_ADMIN_URL}/login`;

export const ROLES = {
  UZLETVEZETO: 'Üzletvezető',
  HR: 'HR-es',
  RAKTARVEZETO: 'Raktárvezető',
  RAKTARKEZELO: 'Raktárkezelő',
  PENZTAROS: 'Pénztáros',
  EGYEB: 'Egyéb dolgozó',
} as const;

export type Role = typeof ROLES[keyof typeof ROLES];

interface User {
  id: number;
  username: string;
  name: string;
}

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  role: Role | null;
  canWrite: (resource: string) => boolean;
  canAccess: (resource: string) => boolean;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await fetch(AUTH_API_PATH, { credentials: 'include' });
        if (res.ok) {
          const userData = await res.json();
          setUser(userData);
        } else {
          setUser(null);
        }
      } catch (error) {
        console.error("Auth check failed:", error);
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

  const role = user?.name as Role | null;

  const canWrite = (resource: string): boolean => {
    if (!role) return false;

    if (role === ROLES.UZLETVEZETO) return true;

    switch (resource) {
      case 'category':
      case 'sub_category':
      case 'product_type':
      case 'storing_condition':
      case 'brand':
      case 'ingredient':
        return role === ROLES.RAKTARVEZETO;

      case 'product':
        return role === ROLES.RAKTARVEZETO || role === ROLES.RAKTARKEZELO;

      case 'user':
      case 'contract':
      case 'contract_type':
        return role === ROLES.HR;

      default:
        return false;
    }
  };

  const canAccess = (resource: string): boolean => {
    if (!role) return false;

    if (resource === 'search' || resource === 'dashboard') return true;

    return canWrite(resource) || role === ROLES.UZLETVEZETO;
  };

  const logout = () => {
    document.cookie = 'auth_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    window.location.href = LOGIN_PATH;
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        role,
        canWrite,
        canAccess,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
