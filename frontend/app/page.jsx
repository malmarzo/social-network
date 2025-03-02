"use client";

import { React } from "react";
import style from "../styles/HomePage.module.css";
import ProfileWindow from "./components/Home/ProfileWindow";
import PostsFeed from "./components/Home/PostsFeed";

export default function Home() {
  return (
    <div className="container max-w-full flex justify-between">
      <div className={style.userProfileDiv}>
        <ProfileWindow />
      </div>
      <div className={style.postsDiv}>
        <PostsFeed />
      </div>
      <div className={style.searchDiv}>Search</div>
    </div>
  );
}
