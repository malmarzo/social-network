"use client";
import { React, useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";
import SuccessAlert from "../components/Alerts/SuccessAlert";
import FailAlert from "../components/Alerts/FailAlert";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";

const LogoutButton = () => {
  const [success, setSuccess] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const router = useRouter();

  const { setIsLoggedIn } = useAuth();

  const handleLogout = async (e) => {
    e.preventDefault();
    const response = await invokeAPI("logout", null, "POST");
    if (response.code === 200) {
      setSuccess(true);
      setIsLoggedIn(false);
      router.push("/login"); // Redirect to homepage
    } else {
      setSuccess(false);
      setErrorMsg(response.error_msg);
    }
  };

  return (
    <>
      {!success && errorMsg && <FailAlert msg={errorMsg} />}

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
