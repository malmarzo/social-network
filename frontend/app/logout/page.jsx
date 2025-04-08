"use client";
import { React, useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";
import FailAlert from "../components/Alerts/FailAlert";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { useAlert } from "../components/Alerts/PopUp";
import styles from "@/styles/AuthButtons.module.css"; // Import the styles

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
      <button
        onClick={handleLogout}
        className={styles.button} // Add this new class
      >
        Logout
      </button>
    </>
  );
};

export default LogoutButton;
