"use client";

import { React } from "react";
import style from "../styles/HomePage.module.css";
import ProfileWindow from "./components/Home/ProfileWindow";
import PostsFeed from "./components/Home/PostsFeed";
import Explore from "./components/Home/Explore";

export default function Home() {
  return (
    <div className="container max-w-full flex justify-between">
      <div className={style.userProfileDiv}>
        <ProfileWindow />
      </div>
      <div className={style.postsDiv}>
        <PostsFeed isGroup={false} />
      </div>
      <div className={style.searchDiv}>
        <Explore />
      </div>
    </div>
  );
}
