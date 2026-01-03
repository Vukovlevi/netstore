import { useState, useEffect } from 'react';
import { Outlet } from 'react-router-dom';

const AUTH_API_PATH = './api/auth';
const LOGIN_PATH = 'http://localhost:8000/login';

const useAuthStatus = () => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await fetch(AUTH_API_PATH);
        console.log(res);
        setIsAuthenticated(true);
      } catch (error) {
        console.error("Auth check failed:", error);
        setIsAuthenticated(true);
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []); 

  return { isAuthenticated, isLoading };
};

interface RequireAuthProps {
}

export default function RequireAuth({}: RequireAuthProps) {
  const { isAuthenticated, isLoading } = useAuthStatus();

  if (isLoading) {
    return <div>Loading...</div>; 
  }

  if (!isAuthenticated) {
    return window.location.href=LOGIN_PATH;
  }

  return <Outlet />;
}