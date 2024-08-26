import { getToken, getUserIDFromToken } from "api/auth";
import { ReactNode, ReactElement, useState, useContext, useEffect } from "react";
import { User, UserContext } from "context/UserContext"


interface UserProviderProps {
    children: ReactNode;
  }
  
export const UserProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);

  const fetchUser = async () => {
    try {
      const token = getToken();
      const userID = getUserIDFromToken(token);
      if (!userID) throw new Error('Invalid token');

      const response = await fetch(`http://localhost:8080/api/users/${userID}`, {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) throw new Error('Failed to fetch user data');

      const result = await response.json();
      setUser(result.data);
    } catch (error) {
      console.error('Error fetching user data:', error);
    }
  };

  useEffect(() => {
    const initializeUser = async () => {
      await fetchUser();
    };
  
    initializeUser();
  }, []);  

  return (
    <UserContext.Provider value={{ user, fetchUser }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};