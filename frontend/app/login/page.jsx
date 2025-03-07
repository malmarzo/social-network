"use client";
import { React, useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";
import SuccessAlert from "../components/Alerts/SuccessAlert";
import FailAlert from "../components/Alerts/FailAlert";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";

const LogInForm = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState("");
  const [success, setSuccess] = useState(false);
  const router = useRouter();
  const { setIsLoggedIn } = useAuth();

  const handleLogin = async (e) => {
    e.preventDefault();
    if (!email || !password) {
      setErrorMsg("Please fill in all required fields.");
      setSuccess(false);
      return;
    }

    const formData = new FormData();
    formData.append("email_nickname", email);
    formData.append("password", password);

    const response = await invokeAPI("login", formData, "POST");
    if (response.code === 200) {
      setSuccess(true);
      setEmail("");
      setPassword("");
      setIsLoggedIn(true, response.data);
      router.push("/");
    } else {
      setSuccess(false);
      setErrorMsg(response.error_msg);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-lg w-full max-w-md">
        <h2 className="text-2xl font-bold text-center text-gray-800 mb-6">
          Login
        </h2>
        {errorMsg && !success && <FailAlert msg={errorMsg} />}
        <form onSubmit={handleLogin} className="space-y-4">
          <input
            type="text"
            placeholder="Nickname or Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            style={{ background: "white", color: "black" }}
          />
          <input
            type="password"
            placeholder="Password"
            style={{ background: "white", color: "black" }}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            type="submit"
            className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition"
          >
            Login
          </button>
        </form>
      </div>
    </div>
  );
};

export default LogInForm;
