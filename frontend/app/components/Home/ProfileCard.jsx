import React from "react";
import style from "@/styles/ProfileCard.module.css";

const Card = ({ imageSrc }) => {
  return (
    <div className={style.card}>
      <div className={style.profileImage}>
        <img src={imageSrc} alt="Profile" />
      </div>
      <div className={style.textContainer}>
        <p className={style.name}>@Pepper Potts</p>
        <div className={style.stats}>
          <span>
            <p className={style.statNumber}>120</p>
            <p className={style.statLabel}>Posts</p>
          </span>
          <span>
            <p className={style.statNumber}>5.2K</p>
            <p className={style.statLabel}>Followers</p>
          </span>
          <span>
            <p className={style.statNumber}>200</p>
            <p className={style.statLabel}>Following</p>
          </span>
        </div>
        <button className={style.viewProfileButton}>View Profile</button>
      </div>
    </div>
  );
};

export default Card;
