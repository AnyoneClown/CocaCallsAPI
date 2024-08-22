import React, { ReactNode, useEffect, useState, startTransition  } from 'react';
import { useNavigate } from 'react-router-dom';
import { getToken, verifyToken } from 'api/auth';
import paths from 'routes/paths';

interface AuthGuardProps {
  children: ReactNode;
}

const AuthGuard: React.FC<AuthGuardProps> = ({ children }) => {
  const navigate = useNavigate();
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);

  useEffect(() => {
    const checkAuth = async () => {
      const token = getToken();
      if (!token) {
        console.log('No token found, redirecting to login...');
        setIsAuthenticated(false);
        navigate(paths.login);
        return;
      }

      try {
        const isValid = await verifyToken(token);
        if (isValid) {
          setIsAuthenticated(true);
        } else {
          console.log('Invalid token, redirecting to login...');
          setIsAuthenticated(false);
          navigate(paths.login);
        }
      } catch (error) {
        console.error('Error verifying token:', error);
        setIsAuthenticated(false);
        navigate(paths.login);
      }
    };

    checkAuth();
  }, [navigate]);

  if (isAuthenticated === null) {
    // Returning null ensures that nothing is rendered until authentication status is determined
    return null;
  }

  return isAuthenticated ? <>{children}</> : null;
};

export default AuthGuard;
