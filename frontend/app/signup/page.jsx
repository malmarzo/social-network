"use client";
import Link from "next/link";
import { React, useState } from "react";
import styles from "./SignUpForm.module.css";
import { invokeAPI } from "@/utils/invokeAPI";
import SuccessAlert from "../components/Alerts/SuccessAlert";
import FailAlert from "../components/Alerts/FailAlert";
import {
  validateDOB,
  validateNickname,
  validatePassword,
  validateTextOnly,
} from "@/utils/formValidators";

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

    const formData = new FormData();
    formData.append("first_name", firstName);
    formData.append("last_name", lastName);
    formData.append("email", email);
    formData.append("password", password);
    formData.append("dob", dob);
    formData.append("nickname", nickname);
    formData.append("about_me", aboutMe);

    if (!firstName || !lastName || !email || !password || !dob || !nickname || (!validateDOB(dob) && dob)) {
      setErrorMsg("Please fill in all required fields.");
      setSuccess(false);
      return;
    }

    if (avatar) {
      formData.append("avatar", avatar);
    }

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
    } else {
      setSuccess(false);
      // Display the error message
      setErrorMsg(response.error_msg);
    }
  };

  return (
    <>
      {success && (
        <SuccessAlert
          msg="Signup successful! Please signin."
          link={"/"}
          linkText={"Login"}
        />
      )}
      {!success && errorMsg && <FailAlert msg={errorMsg} />}
      <div className={styles.wrapper}>
        <form className={styles.form} onSubmit={handleFormSubmit}>
          <p className={styles.title}>Signup</p>
          <div className={styles.flex}>
            <label>
              <input
                className={styles.input}
                type="text"
                value={firstName}
                onChange={handleFirstNameChange}
                required
              />
              <span>
                First Name <span className={styles.errMsg}>{firstErr}</span>
              </span>
            </label>
            <label>
              <input
                className={styles.input}
                type="text"
                required
                value={lastName}
                onChange={handleLastNameChange}
              />
              <span>
                Last Name <span className={styles.errMsg}>{lastErr}</span>
              </span>
            </label>
          </div>
          <label>
            <input
              className={styles.input}
              type="text"
              value={nickname}
              onChange={handleNicknameChange}
              required
            />
            <span>
              Nickname <span className={styles.errMsg}>{nickErr}</span>
            </span>
          </label>
          <label>
            <input
              className={styles.input}
              type="email"
              required
              value={email}
              onChange={handleEmailChange}
            />
            <span>
              Email <span className={styles.errMsg}>{emailErr}</span>
            </span>
          </label>

          <label>
            <input
              className={styles.input}
              type="date"
              required
              value={dob}
              onChange={handleDobChange}
            />
            <span>
              Date of Birth <span className={styles.errMsg}>{dateErr}</span>
            </span>
          </label>

          <label>
            <input
              className={styles.input}
              type="file"
              accept="image/jpeg, image/png, image/gif"
              onChange={handleAvatarChange}
            />
            <span>
              Avatar (Optional){" "}
              <span className={styles.errMsg}>{avatarErr}</span>
            </span>
          </label>
          <label>
            <textarea
              className={styles.input}
              rows="1"
              style={{ resize: "none" }}
              value={aboutMe}
              onChange={handleAboutMeChange}
              maxLength="100"
            ></textarea>
            <span>
              About (Optional) <span className={styles.errMsg}>{aboutErr}</span>
            </span>
          </label>
          <label>
            <input
              className={styles.input}
              type="password"
              required
              value={password}
              onChange={handlePasswordChange}
            />
            <span>
              Password <span className={styles.errMsg}>{passErr}</span>
            </span>
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
