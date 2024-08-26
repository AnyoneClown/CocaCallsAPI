import React, { ReactNode, useEffect, useState } from 'react';
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
    let isMounted = true;

    const checkAuth = async () => {
      const token = getToken();
      if (!token) {
        console.log('No token found, redirecting to login...');
        if (isMounted) {
          setIsAuthenticated(false);
          navigate(paths.login);
        }
        return;
      }

      try {
        const isValid = await verifyToken(token);
        if (isValid) {
          if (isMounted) {
            setIsAuthenticated(true);
          }
        } else {
          console.log('Invalid token, redirecting to login...');
          if (isMounted) {
            setIsAuthenticated(false);
            navigate(paths.login);
          }
        }
      } catch (error) {
        console.error('Error verifying token:', error);
        if (isMounted) {
          setIsAuthenticated(false);
          navigate(paths.login);
        }
      }
    };

    checkAuth();

    return () => {
      isMounted = false;
    };
  }, [navigate]);

  if (isAuthenticated === null) {
    return null;
  }

  return isAuthenticated ? <>{children}</> : null;
};

export default AuthGuard;