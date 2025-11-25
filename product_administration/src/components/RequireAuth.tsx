import { useState, useEffect } from 'react';
import { Outlet } from 'react-router-dom';

// Define the API path and Login path constants
const AUTH_API_PATH = './api/auth';
const LOGIN_PATH = 'http://localhost:8000/login';

// Custom hook to handle the async authentication logic
// This helps keep the main component cleaner
const useAuthStatus = () => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // We define an async function inside useEffect to handle the fetch
    const checkAuth = async () => {
      try {
        const res = await fetch(AUTH_API_PATH);
        // Assuming a status of 200-299 means success/authenticated
        setIsAuthenticated(res.ok);
      } catch (error) {
        // Handle network errors or other failures as unauthenticated
        console.error("Auth check failed:", error);
        setIsAuthenticated(false);
      } finally {
        // Always set loading to false once the check is complete
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []); // The empty dependency array [] ensures this effect runs only ONCE after the initial render

  return { isAuthenticated, isLoading };
};

interface RequireAuthProps {
  // You might add roles/permissions checks here later, but for now, we'll keep it simple
}

// The main component is NOT async and returns JSX
export default function RequireAuth({}: RequireAuthProps) {
  const { isAuthenticated, isLoading } = useAuthStatus();

  // 1. Show a loading screen while the API call is in progress
  if (isLoading) {
    return <div>Loading...</div>; // You can replace this with a proper spinner
  }

  // 2. Once loading is complete, check authentication status
  if (!isAuthenticated) {
    // Use <Navigate> component from react-router-dom for internal redirects.
    // This is the correct way to handle redirects within a React Router context.
    // The previous window.location.href method forces a full page reload, which is usually undesirable.
    return window.location.href=LOGIN_PATH;
  }

  // 3. If authenticated, render the child routes
  return <Outlet />;
}