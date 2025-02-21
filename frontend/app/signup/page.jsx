import Link from "next/link";
import React from "react";
import styles from "./SignUpForm.module.css";

const SignUpForm = () => {
  return (
    <div className={styles.wrapper}>
      <form className={styles.form}>
        <p className={styles.title}>Signup</p>
        <div className={styles.flex}>
          <label>
            <input className={styles.input} type="text" required />
            <span>First Name</span>
          </label>
          <label>
            <input className={styles.input} type="text" required />
            <span>Last Name</span>
          </label>
        </div>
        <label>
          <input className={styles.input} type="email" required />
          <span>Email</span>
        </label>
        <label>
          <input className={styles.input} type="password" required />
          <span>Password</span>
        </label>
        <label>
          <input className={styles.input} type="date" required />
          <span>Date of Birth</span>
        </label>
        <label>
          <input className={styles.input} type="text" />
          <span>Nickname</span>
        </label>
        <label>
          <input className={styles.input} type="file" accept="image/*" />
          <span>Avatar/Image (Optional)</span>
        </label>
        <label>
          <textarea
            className={styles.input}
            rows="2"
            style={{ resize: "none" }}
          ></textarea>
          <span>About Me (Optional)</span>
        </label>
        <button className={styles.submit}>Submit</button>
        <p className={styles.signin}>
          Already have an account? <Link href="/">Signin</Link>
        </p>
      </form>
    </div>
  );
};

export default SignUpForm;
