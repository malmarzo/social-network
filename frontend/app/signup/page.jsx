"use client";
import Link from "next/link";
import { React, useState } from "react";
import styles from "./SignUpForm.module.css";
import { invokeAPI } from "@/utils/invokeAPI";
import SuccessAlert from "../components/Alerts/SuccessAlert";
import FailAlert from "../components/Alerts/FailAlert";

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

  const handleFirstNameChange = (e) => setFirstName(e.target.value);
  const handleLastNameChange = (e) => setLastName(e.target.value);
  const handleEmailChange = (e) => setEmail(e.target.value);
  const handlePasswordChange = (e) => setPassword(e.target.value);
  const handleDobChange = (e) => setDob(e.target.value);
  const handleNicknameChange = (e) => setNickname(e.target.value);
  const handleAvatarChange = (e) => setAvatar(e.target.files[0]);
  const handleAboutMeChange = (e) => setAboutMe(e.target.value);

  const handleFormSubmit = async (e) => {
    e.preventDefault();
    const body = {
      first_name: firstName,
      last_name: lastName,
      email: email,
      password: password,
      dob: dob,
      nickname: nickname,
      about_me: aboutMe,
      avatar: avatar,
    };

    const response = await invokeAPI("signup", body, "POST");
    console.log(response);
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
    } else {
      setSuccess(false);
      // Display the error message
      setErrorMsg(response.error_msg);

    }
  };

  return (
    <>
      {success && <SuccessAlert msg="Signup successful! Please signin." link={"/"} linkText={"Login"} />}
      {!success && errorMsg && <FailAlert msg={errorMsg} />}
      <div className={styles.wrapper}>
        <form className={styles.form} onSubmit={handleFormSubmit}>
          <p className={styles.title}>Signup</p>
          <div className={styles.flex}>
            <label>
              <input
                className={styles.input}
                type="text"
                required
                value={firstName}
                onChange={handleFirstNameChange}
              />
              <span>First Name</span>
            </label>
            <label>
              <input
                className={styles.input}
                type="text"
                required
                value={lastName}
                onChange={handleLastNameChange}
              />
              <span>Last Name</span>
            </label>
          </div>
          <label>
            <input
              className={styles.input}
              type="email"
              required
              value={email}
              onChange={handleEmailChange}
            />
            <span>Email</span>
          </label>
          <label>
            <input
              className={styles.input}
              type="password"
              required
              value={password}
              onChange={handlePasswordChange}
            />
            <span>Password</span>
          </label>
          <label>
            <input
              className={styles.input}
              type="date"
              required
              value={dob}
              onChange={handleDobChange}
            />
            <span>Date of Birth</span>
          </label>
          <label>
            <input
              className={styles.input}
              type="text"
              value={nickname}
              onChange={handleNicknameChange}
              required
            />
            <span>Nickname</span>
          </label>
          <label>
            <input
              className={styles.input}
              type="file"
              accept="image/*"
              onChange={handleAvatarChange}
            />
            <span>Avatar (Optional)</span>
          </label>
          <label>
            <textarea
              className={styles.input}
              rows="2"
              style={{ resize: "none" }}
              value={aboutMe}
              onChange={handleAboutMeChange}
            ></textarea>
            <span>About (Optional)</span>
          </label>
          <button className={styles.submit} type="submit">
            Submit
          </button>
          <p className={styles.signin}>
            Already have an account? <Link href="/">Signin</Link>
          </p>
        </form>
      </div>
    </>
  );
};

export default SignUpForm;
