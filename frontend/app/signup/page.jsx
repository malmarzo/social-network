"use client";
import Link from "next/link";
import { React, useEffect, useState } from "react";
import styles from "@/styles/SignUpForm.module.css";
import { invokeAPI } from "@/utils/invokeAPI";
import FailAlert from "../components/Alerts/FailAlert";
import {
  validateDOB,
  validateNickname,
  validatePassword,
  validateTextOnly,
} from "@/utils/formValidators";
import { useRouter } from "next/navigation";
import LoadingSpinner from "../components/LoadingSpinner";

const SignUpForm = () => {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [dob, setDob] = useState("");
  const [nickname, setNickname] = useState("");
  const [avatar, setAvatar] = useState("");
  const [aboutMe, setAboutMe] = useState("");
  const [success, setSuccess] = useState();
  const [errorMsg, setErrorMsg] = useState("");
  const [firstErr, setFirstErr] = useState("");
  const [lastErr, setLastErr] = useState("");
  const [nickErr, setNickErr] = useState("");
  const [emailErr, setEmailErr] = useState("");
  const [dateErr, setDateErr] = useState("");
  const [avatarErr, setAvatarErr] = useState("");
  const [aboutErr, setAboutErr] = useState("");
  const [passErr, setPassErr] = useState("");
  const [loading, setLoading] = useState(false);

  const router = useRouter();

  //Validates the first name on every change
  const handleFirstNameChange = (e) => {
    const value = e.target.value;
    setFirstName(value);
    if (!validateTextOnly(value) && value) {
      setFirstErr("letters only!");
    } else {
      setFirstErr("");
    }
  };

  //Validates the last name on every change
  const handleLastNameChange = (e) => {
    const value = e.target.value;
    setLastName(value);
    if (!validateTextOnly(value) && value) {
      setLastErr("letters only!");
    } else {
      setLastErr("");
    }
  };

  //Sets the email value
  const handleEmailChange = (e) => setEmail(e.target.value);

  //Validate the password
  const handlePasswordChange = (e) => {
    const value = e.target.value;
    setPassword(value);
    if (!validatePassword(value) && value) {
      setPassErr("min 8 chars, 1 upper, 1 lower, 1 number");
    } else {
      setPassErr("");
    }
  };

  //Validates the dob
  const handleDobChange = (e) => {
    const value = e.target.value;
    setDob(value);
    if (!validateDOB(value) && value) {
      setDateErr("You must be at least 15 years old to signup.");
    } else {
      setDateErr("");
    }
  };

  //validates the nickname
  const handleNicknameChange = (e) => {
    const value = e.target.value;
    setNickname(value);
    if (!validateNickname(value) && value) {
      setNickErr("letters, numbers, and underscores only!");
    } else {
      setNickErr("");
    }
  };

  //Handles avatar change
  const handleAvatarChange = (e) => {
    const file = e.target.files[0];
    const allowedExtensions = ["image/jpeg", "image/png", "image/gif"];
    if (file && !allowedExtensions.includes(file.type)) {
      setAvatarErr("Avatar should be a JPEG, PNG, or GIF image.");
      setAvatar("");
    } else {
      setAvatarErr("");
      setAvatar(file);
    }
  };

  //Sets the about me value
  const handleAboutMeChange = (e) => setAboutMe(e.target.value);

  //Handles form submission
  const handleFormSubmit = async (e) => {
    e.preventDefault();

    setLoading(true);

    const formData = new FormData();
    formData.append("first_name", firstName);
    formData.append("last_name", lastName);
    formData.append("email", email);
    formData.append("password", password);
    formData.append("dob", dob);
    formData.append("nickname", nickname);
    formData.append("about_me", aboutMe);

    if (
      !firstName ||
      !lastName ||
      !email ||
      !password ||
      !dob ||
      !nickname ||
      (!validateDOB(dob) && dob)
    ) {
      setErrorMsg("Please fill in all required fields.");
      setSuccess(false);
      return;
    }

    if (avatar) {
      formData.append("avatar", avatar);
    }

    try {
      setTimeout(() => {});
      const response = await invokeAPI("signup", formData, "POST");
      if (response.code === 200) {
        setSuccess(true);

        // Clear the form
        setFirstName("");
        setLastName("");
        setEmail("");
        setPassword("");
        setDob("");
        setNickname("");
        setAvatar("");
        setAboutMe("");
        router.push("/login");
      } else {
        setSuccess(false);
        // Display the error message
        setErrorMsg(response.error_msg);
      }
    } catch (error) {
      setErrorMsg("Something went wrong. Please try again later.");
      setSuccess(false);
      setLoading(false);
      return;
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className={styles.main}>
      <div className={styles.formCard}>
        <h1 className={styles.title}>Create Account</h1>
        {!success && errorMsg && <FailAlert msg={errorMsg} />}

        <form className={styles.form} onSubmit={handleFormSubmit}>
          <div className={styles.flex}>
            <div className={styles.inputGroup}>
              <input
                className={styles.input}
                type="text"
                placeholder="First Name"
                value={firstName}
                onChange={handleFirstNameChange}
                required
              />
              {firstErr && <span className={styles.errMsg}>{firstErr}</span>}
            </div>
            <div className={styles.inputGroup}>
              <input
                className={styles.input}
                type="text"
                placeholder="Last Name"
                value={lastName}
                onChange={handleLastNameChange}
                required
              />
              {lastErr && <span className={styles.errMsg}>{lastErr}</span>}
            </div>
          </div>

          <div className={styles.inputGroup}>
            <input
              className={styles.input}
              type="text"
              placeholder="Nickname"
              value={nickname}
              onChange={handleNicknameChange}
              required
            />
            {nickErr && <span className={styles.errMsg}>{nickErr}</span>}
          </div>

          <div className={styles.inputGroup}>
            <input
              className={styles.input}
              type="email"
              placeholder="Email"
              value={email}
              onChange={handleEmailChange}
              required
            />
            {emailErr && <span className={styles.errMsg}>{emailErr}</span>}
          </div>

          <div className={styles.inputGroup}>
            <div className={styles.specialInput}>
              <small className={styles.specialLabel}>Date of Birth</small>
              <input
                className={styles.input}
                type="date"
                value={dob}
                onChange={handleDobChange}
                required
              />
            </div>
            {dateErr && <span className={styles.errMsg}>{dateErr}</span>}
          </div>

          <div className={styles.inputGroup}>
            <div className={styles.specialInput}>
              <small className={styles.specialLabel}>Avatar (Optional)</small>
              <input
                className={`${styles.input} ${styles.fileInput}`}
                type="file"
                accept="image/jpeg, image/png, image/gif"
                onChange={handleAvatarChange}
              />
            </div>
            {avatarErr && <span className={styles.errMsg}>{avatarErr}</span>}
          </div>

          <div className={styles.inputGroup}>
            <textarea
              className={`${styles.input} ${styles.textarea}`}
              placeholder="About (Optional)"
              value={aboutMe}
              onChange={handleAboutMeChange}
              maxLength="50"
            />
            {aboutErr && <span className={styles.errMsg}>{aboutErr}</span>}
          </div>

          <div className={styles.inputGroup}>
            <input
              className={styles.input}
              type="password"
              placeholder="Password"
              value={password}
              onChange={handlePasswordChange}
              required
            />
            {passErr && <span className={styles.errMsg}>{passErr}</span>}
          </div>

          <button
            className={`${styles.button} ${loading ? styles.loading : ""}`}
            type="submit"
            disabled={loading}
          >
            {loading ? (
              <div className={styles.buttonContent}>
                <LoadingSpinner />
              </div>
            ) : (
              "Create Account"
            )}
          </button>

          <p className={styles.signin}>
            Already have an account? <Link href="/login">Sign in</Link>
          </p>
        </form>
      </div>
    </main>
  );
};

export default SignUpForm;
