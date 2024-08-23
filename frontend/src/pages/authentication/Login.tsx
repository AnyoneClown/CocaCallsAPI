import {
  Box,
  Link,
  Paper,
  Stack,
  Button,
  Divider,
  Checkbox,
  FormGroup,
  TextField,
  IconButton,
  Typography,
  InputAdornment,
  FormControlLabel,
} from '@mui/material';
import IconifyIcon from 'components/base/IconifyIcon';
import { useState, ReactElement } from 'react';
import { useNavigate } from 'react-router-dom';
import { rootPaths } from 'routes/paths';
import { loginUser, setToken, UserLoginResponse } from 'api/auth'
import Image from 'components/base/Image';
import logoWithText from '/Logo-with-text.png';

const Login = (): ReactElement => {
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [message, setMessage] = useState<string | null>(null);

  const handleSubmit = async () => {
    if (!email || !password) {
      setMessage('Please fill in all fields');
      return;
    }

    try {
      const response: UserLoginResponse = await loginUser({ email, password });
      setToken(response.token);
      navigate(rootPaths.homeRoot);
    } catch (error) {
      console.error('Failed to login', error);
      setMessage(error.message);
    }
  };


  const handleClickShowPassword = () => {
    setShowPassword((prevShowPassword) => !prevShowPassword);
  };

  const handleGoogleSignIn = () => {
    window.location.href = "http://localhost:8080/api/auth/google/";
  };

  return (
    <>
      <Box component="figure" mb={5} mx="auto" textAlign="center">
        <Link href={rootPaths.homeRoot}>
          <Image src={logoWithText} alt="logo with text" height={200} />
        </Link>
      </Box>
      <Paper
        sx={{
          py: 6,
          px: { xs: 5, sm: 7.5 },
        }}
      >
        <Stack justifyContent="center" gap={5}>
          <Typography variant="h3" textAlign="center" color="text.secondary">
            Log In
          </Typography>
          <Typography variant="h6" fontWeight={500} textAlign="center" color="text.primary">
            Don’t have an account?{' '}
            <Link href="/authentication/sign-up" underline="none">
              Sign up
            </Link>
          </Typography>
          <TextField
            variant="filled"
            label="Email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            sx={{
              '.MuiFilledInput-root': {
                bgcolor: 'grey.A100',
                ':hover': {
                  bgcolor: 'background.default',
                },
                ':focus': {
                  bgcolor: 'background.default',
                },
                ':focus-within': {
                  bgcolor: 'background.default',
                },
              },
              borderRadius: 2,
            }}
          />
          <TextField
            variant="filled"
            label="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            type={showPassword ? 'text' : 'password'}
            sx={{
              '.MuiFilledInput-root': {
                bgcolor: 'grey.A100',
                ':hover': {
                  bgcolor: 'background.default',
                },
                ':focus': {
                  bgcolor: 'background.default',
                },
                ':focus-within': {
                  bgcolor: 'background.default',
                },
              },
              borderRadius: 2,
            }}
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton
                    aria-label="toggle password visibility"
                    onClick={handleClickShowPassword}
                    size="small"
                    edge="end"
                    sx={{
                      mr: 2,
                    }}
                  >
                    {showPassword ? (
                      <IconifyIcon icon="el:eye-open" color="text.secondary" />
                    ) : (
                      <IconifyIcon icon="el:eye-close" color="text.primary" />
                    )}
                  </IconButton>
                </InputAdornment>
              ),
            }}
          />
          {message && (
            <Typography variant="body1" color="error.main" textAlign="center" sx={{ mt: 2 }}>
              {message}
            </Typography>
          )}
          <Button
            onClick={handleSubmit}
            sx={{
              fontWeight: 'fontWeightRegular',
            }}
          >
            Log In
          </Button>
          <Divider />
          <Typography textAlign="center" color="text.secondary" variant="body1">
            Or sign in using:
          </Typography>
          <Stack gap={1.5} direction="row" justifyContent="space-between">
            <Button
              startIcon={<IconifyIcon icon="flat-color-icons:google" />}
              variant="outlined"
              onClick={handleGoogleSignIn}
              fullWidth
              sx={{
                fontSize: 'subtitle2.fontSize',
                fontWeight: 'fontWeightRegular',
              }}
            >
              Google
            </Button>
          </Stack>
        </Stack>
      </Paper>
    </>
  );
};

export default Login;
