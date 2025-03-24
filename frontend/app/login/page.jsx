"use client";
import { React, useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";
import SuccessAlert from "../components/Alerts/SuccessAlert";
import FailAlert from "../components/Alerts/FailAlert";
import LoadingSpinner from "../components/LoadingSpinner";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import styles from "@/styles/Login.module.css";

const LogInForm = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const router = useRouter();
  const { setIsLoggedIn } = useAuth();

  const handleLogin = async (e) => {
    e.preventDefault();
    setErrorMsg("");
    setLoading(true);
    if (!email || !password || email.trim() === "" || password.trim() === "") {
      setErrorMsg("Please fill in all required fields.");
      setSuccess(false);
      return;
    }

    const formData = new FormData();
    formData.append("email_nickname", email.trim());
    formData.append("password", password.trim());

    try {
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
        setLoading(false);
      }
    } catch (error) {
      setErrorMsg("An error occurred. Please try again later.");
      setSuccess(false);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className={styles.main}>
      <div className={styles.formCard}>
        <h1 className={styles.title}>Welcome Back</h1>
        {errorMsg && !success && <FailAlert msg={errorMsg} />}
        <form onSubmit={handleLogin} className={styles.form}>
          <div className={styles.inputGroup}>
            <input
              type="text"
              placeholder="Nickname or Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className={styles.input}
            />
          </div>
          <div className={styles.inputGroup}>
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={styles.input}
            />
          </div>
          <button
            type="submit"
            className={`${styles.button} ${loading ? styles.loading : ""}`}
            disabled={loading}
          >
            {loading ? (
              <div className={styles.buttonContent}>
                <LoadingSpinner />
              </div>
            ) : (
              "Sign In"
            )}
          </button>
        </form>
      </div>
    </main>
  );
};

export default LogInForm;
