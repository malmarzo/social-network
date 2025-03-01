"use client"; // Required for Next.js App Router (if using app directory)

import { createContext, useContext, useEffect, useState } from "react";
import { invokeAPI } from "../utils/invokeAPI";

//Setting value to be false at the start
const AuthContext = createContext({ isLoggedIn: false });

export function AuthProvider({ children }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loading, setLoading] = useState(true);

  const updateAuthState = async (newState) => {
    setIsLoggedIn(newState);
  };

  useEffect(() => {
    console.log("Auth State Updated:", isLoggedIn);
  }, [isLoggedIn]);

  useEffect(() => {
    //Validating user authentication from the backend
    async function checkAuth() {
      try {
        const response = await invokeAPI("session", null, "GET");
        await updateAuthState(response.code === 200);
      } catch (error) {
        console.error("Auth check failed:", error);
        await updateAuthState(false);
      } finally {
        setLoading(false);
      }
    }

    checkAuth();
  }, []);

  return (
    <AuthContext.Provider
      //Values that could be accessed through other components
      value={{
        isLoggedIn,
        setIsLoggedIn: updateAuthState,
        loading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
