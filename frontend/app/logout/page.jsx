"use client";
import { React, useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";
import SuccessAlert from "../components/Alerts/SuccessAlert";
import FailAlert from "../components/Alerts/FailAlert";
import { useRouter } from "next/navigation";

const LogoutPage = () => {
  const [success, setSuccess] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const router = useRouter();

  const handleLogout = async (e) => {
    e.preventDefault();
    const response = await invokeAPI("logout", null, "POST");
    if (response.code === 200) {
      setSuccess(true);
      router.push("/login"); // Redirect to homepage
    } else {
      setSuccess(false);
      setErrorMsg(response.error_msg);
    }
  };

  return (
    <>
      {/* {success && (
        <SuccessAlert
          msg="Logout successful!"
          link="/"
          linkText="Home"
        />
      )} */}
      {!success && errorMsg && <FailAlert msg={errorMsg} />}
      
      <div>
        <button onClick={handleLogout} style={{ backgroundColor: "white", color: "black", border: "1px solid black", padding: "10px 20px", cursor: "pointer" }}>
            Logout
        </button>
      </div>
    </>
  );
};

export default LogoutPage;
