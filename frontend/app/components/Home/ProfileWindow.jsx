import React from "react";
import ProfileCard from "./ProfileCard";
import styles from "@/styles/ProfileWindow.module.css";

const ProfileWindow = () => {
  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Profile</h1>
      <ProfileCard />
    </div>
  );
};

export default ProfileWindow;
