import { ReactElement, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Stack, Typography, CircularProgress } from '@mui/material';
import { rootPaths } from 'routes/paths';

const OAuthCallbackPage = (): ReactElement => {
    const navigate = useNavigate();

    useEffect(() => {
      const urlParams = new URLSearchParams(window.location.search);
      const token = urlParams.get('token');
      const error = urlParams.get('error');
  
      if (token) {
        // Успішна авторизація
        localStorage.setItem('authToken', token);
        setTimeout(() => {
          navigate(rootPaths.homeRoot);
        }, 2000);
      } else if (error) {
        navigate('/login', { state: { error: 'Авторизацію скасовано' } });
      } else {
        navigate('/error');
      }
    }, [navigate]);

  return (
    <Stack
      minHeight="100vh"
      width="fit-content"
      mx="auto"
      justifyContent="center"
      alignItems="center"
      gap={4}
      py={12}
    >
      <Typography variant="h4" color="text.primary">
        Авторизація успішна
      </Typography>
      <Typography
        variant="body1"
        color="text.secondary"
        maxWidth={400}
        textAlign="center"
      >
        Ваш токен отримано. Зараз вас буде перенаправлено на головну сторінку.
      </Typography>
      <CircularProgress size={60} />
    </Stack>
  );
};

export default OAuthCallbackPage;