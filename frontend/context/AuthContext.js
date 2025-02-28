"use client"; // Required for Next.js App Router (if using app directory)

import { createContext, useContext, useEffect, useState } from "react";
import { invokeAPI } from "../utils/invokeAPI";

const AuthContext = createContext({ isLoggedIn: false });

export function AuthProvider({ children }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function checkAuth() {
      try {
        const response = await invokeAPI("session", null, "POST");
        if (response.code === 200) {
          setIsLoggedIn(true);
        } else {
          setIsLoggedIn(false);
        }
      } catch (error) {
        console.error("Auth check failed:", error);
        setIsLoggedIn(false);
      }
      setLoading(false);
    }

    checkAuth();
  }, []);

  return (
    <AuthContext.Provider value={{ isLoggedIn, setIsLoggedIn, loading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
