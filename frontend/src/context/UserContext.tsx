import React, { createContext, useState, useContext, ReactNode } from 'react';

export interface User {
  ID: string;
  Email: string;
  GoogleID?: string;
  Picture?: string;
  Provider?: string;
  VerifiedEmail: boolean;
  IsAdmin: boolean;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt?: string | null;
  Subscription?: {
    ID: string;
    UserID: string;
    StartDate: string;
    EndDate: string;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt?: string | null;
  };
}


interface UserContextType {
  user: User | null;
  fetchUser: () => void;
}

export const UserContext = createContext<UserContextType | undefined>(undefined);

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};