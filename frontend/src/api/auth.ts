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
  user: {
    id: string;
    email: string;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string;
  };
  code: number;
  token: string;
}

export const loginUser = async (credentials: { email: string; password: string }): Promise<UserLoginResponse> => {
  const response = await fetch('http://localhost:8080/api/auth/login/', {
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

export const getToken = (): string | null => {
  return localStorage.getItem('authToken');
}