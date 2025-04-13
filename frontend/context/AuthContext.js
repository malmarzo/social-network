"use client"; // Required for Next.js App Router (if using app directory)

import { createContext, useContext, useEffect, useState } from "react";
import { invokeAPI } from "../utils/invokeAPI";

//Setting value to be false at the start
const AuthContext = createContext({
  isLoggedIn: false,
  userID: "",
  userNickname: "",
});

export function AuthProvider({ children }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [userID, setUserID] = useState("");
  const [userNickname, setUserNickname] = useState("");
  const [loading, setLoading] = useState(true);

  const updateAuthState = async (newState, userData = null) => {
    setIsLoggedIn(newState);
    if (userData && newState === true) {
      setUserID(userData.user_id);
      setUserNickname(userData.user_nickname);
    } else {
      setUserID("");
      setUserNickname("");
    }
  };

  useEffect(() => {
    console.log("Auth State Updated:", isLoggedIn);
  }, [isLoggedIn]);

  useEffect(() => {
    //Validating user authentication from the backend
    async function checkAuth() {
      try {
        const response = await invokeAPI("session", null, "GET");
        if (response.code === 200) {
          await updateAuthState(true, response.data);
        } else {
          await updateAuthState(false);
        }
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
        userID,
        userNickname,
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
