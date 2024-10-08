import * as jwtDecode from "jwt-decode";

export const registerUser = async (userData: { email: string; password: string }) => {
  const response = await fetch('http://localhost:8080/api/auth/register/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    throw new Error('Failed to register user');
  }

  return response.json();
};

export interface UserLoginResponse {
  message: string;
  data: {
    user: {
      id: string;
      email: string;
      CreatedAt: string;
      UpdatedAt: string;
      DeletedAt: string;
    },
    token: string;
  };
  code: number;
}

export const loginUser = async (credentials: { email: string; password: string }): Promise<UserLoginResponse> => {
  const response = await fetch('http://localhost:8080/api/jwt/create/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message);
  }

  const data: UserLoginResponse = await response.json();
  return data;
};

export const setToken = (token: string) => {
  // Store the token in localStorage, cookies or elsewhere
  localStorage.setItem('authToken', token);
};

export const getToken = (): string => {
  return localStorage.getItem('authToken') ?? '';
}

export const verifyToken = async (token: string): Promise<boolean> => {
  try {
    const response = await fetch('http://localhost:8080/api/jwt/verify/', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    });
    return response.ok;
  } catch (error) {
    console.error('Error verifying token:', error);
    return false;
  }
};

export const getUserIDFromToken = (token: string): string | null => {
  const arrayToken = token.split('.');
  const decodedToken = JSON.parse(atob(arrayToken[1]));
  return decodedToken.userID;
};