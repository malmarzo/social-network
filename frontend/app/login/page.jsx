"use client";
//import Link from "next/link";
import { React, useState } from "react";
//import styles from "./SignUpForm.module.css";
import { invokeAPI } from "@/utils/invokeAPI";
import SuccessAlert from "../components/Alerts/SuccessAlert";
import FailAlert from "../components/Alerts/FailAlert";
import { useRouter } from "next/navigation";

const LogInForm = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [errorMsg, setErrorMsg] = useState("");
    const [success, setSuccess] = useState("");
    const router = useRouter();

const handleLogin = async (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append("email",email);
    formData.append("password",password)

    if (!email || !password){
        setErrorMsg("Please fill in all required fields.");
        setSuccess(false);
        return
    }
    const response = await invokeAPI("login", formData, "POST");
    if (response.code === 200){
        setSuccess(true);
        setEmail("");
        setPassword("");
        router.push("/"); // Redirect to homepage
       
    }else{
        setSuccess(false);
        setErrorMsg(response.error_msg);
    }


};

return (
    <>
    {/* {success && (
        <SuccessAlert
          msg="Login successful!."
          link={"/"}
          linkText={"Home"}
        />
      )} */}
      {!success && errorMsg && <FailAlert msg={errorMsg} />}
    <div>
    <h2>Login</h2>
    <input type="text" placeholder="email" value={email} onChange={(e) => setEmail(e.target.value)} />
    <input type="password" placeholder="password" value={password} onChange={(e) => setPassword(e.target.value)} />
    <button onClick={handleLogin}>Login</button>
  </div>
  </>
);
};
export default LogInForm;