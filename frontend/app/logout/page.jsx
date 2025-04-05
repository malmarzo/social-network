"use client";
import { React, useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";
import FailAlert from "../components/Alerts/FailAlert";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { useAlert } from "../components/Alerts/PopUp";

const LogoutButton = () => {
  const [errorMsg, setErrorMsg] = useState("");
  const router = useRouter();
  const { isLoggedIn, setIsLoggedIn } = useAuth();
  const { showAlert } = useAlert();

  const handleLogout = async (e) => {
    e.preventDefault();

    showAlert({
      type: "confirm",
      message: "Are you sure you want to logout?",
      action: async () => {
        try {
          const response = await invokeAPI("logout", null, "POST");
          if (response.code === 200) {
            setIsLoggedIn(false); // Ensure state updates
            router.replace("/login"); // Navigate after state update
          } else {
            setErrorMsg(response.error_msg);
          }
        } catch (error) {
          console.error("Logout failed:", error);
          setErrorMsg("Logout failed. Please try again.");
        }
      },
    });
  };

  return (
    <>
      {errorMsg && <FailAlert msg={errorMsg} />}
      <div>
        <button
          onClick={handleLogout}
          className="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200"
        >
          Logout
        </button>
      </div>
    </>
  );
};

export default LogoutButton;
