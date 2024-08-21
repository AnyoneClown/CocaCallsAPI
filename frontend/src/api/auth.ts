export const registerUser = async (userData: { email: string; password: string }) => {
  const response = await fetch('http://localhost:8080/api/auth/register/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    throw new Error('Failed to register');
  }

  return response.json();
};
