"use client";

import { React } from "react";
import styles from "../styles/HomePage.module.css";
import ProfileWindow from "./components/Home/ProfileWindow";
import PostsFeed from "./components/Home/PostsFeed";
import Explore from "./components/Home/Explore";
import AuthButton from "./components/Buttons/AuthButtons";


export default function Home() {
  return (
    
    <div className={styles.container}>
      
      <div className={styles.userProfileDiv}>
        <ProfileWindow />
      </div>
      <div className={styles.postsDiv}>
        <PostsFeed isGroup={false} />
      </div>
      <div className={styles.searchDiv}>
        <Explore />
      </div>
     
    </div>
  );
}