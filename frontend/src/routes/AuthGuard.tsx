import React, { ReactNode, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getToken } from 'api/auth';
import paths from 'routes/paths';

interface AuthGuardProps {
  children: ReactNode;
}

const AuthGuard: React.FC<AuthGuardProps> = ({ children }) => {
  const navigate = useNavigate();
  const token = getToken();

  useEffect(() => {
    if (!token) {
      console.log('No token found, redirecting to login...');
      navigate(paths.login);
    }
  }, [token, navigate]);

  return token ? <>{children}</> : null;
};

export default AuthGuard;
