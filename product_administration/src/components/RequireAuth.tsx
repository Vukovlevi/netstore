import { Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const LOGIN_PATH = `${import.meta.env.VITE_STORE_ADMIN_URL}/login`;

export default function RequireAuth() {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    window.location.href = LOGIN_PATH;
    return null;
  }

  return <Outlet />;
}
