import React, { useEffect, useState } from "react";
import style from "@/styles/ProfileCard.module.css";
import { invokeAPI } from "@/utils/invokeAPI";

const Card = () => {
  //Data for the profile card
  const [imageSrc, setImageSrc] = useState("/imgs/defaultAvatar.jpg");
  const [loading, setIsLoading] = useState(true);
  const [error, setError] = useState(false);
  const [username, setUsername] = useState("");
  const [posts, setPosts] = useState(0);
  const [followers, setFollowers] = useState(0);
  const [following, setFollowing] = useState(0);

  useEffect(() => {
    //Fetches the user's data for the card
    async function fetchData() {
      try {
        const response = await invokeAPI("profileCard", null, "GET");

        if (!response || response.code !== 200) {
          setError(true);
          return;
        }

        const profileData = response.data;
        console.log("Profile data:", profileData);

        if (profileData.avatar) {
          const imageDataUrl = `data:${profileData.avatar_mime_type};base64,${profileData.avatar}`;
          setImageSrc(imageDataUrl);
        }

        setUsername(profileData.nickname || "");
        setPosts(profileData.posts || 0);
        setFollowers(profileData.followers || 0);
        setFollowing(profileData.following || 0);
      } catch (error) {
        console.error("Failed to fetch profile data:", error);
        setError(true);
      } finally {
        setIsLoading(false);
      }
    }

    fetchData();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error loading profile</div>;
  }

  return (
    <div className={style.card}>
      <div className={style.profileImage}>
        <img src={imageSrc} alt="Profile" />
      </div>
      <div className={style.textContainer}>
        <p className={style.name}>@{username}</p>
        <div className={style.stats}>
          <span>
            <p className={style.statNumber}>{posts}</p>
            <p className={style.statLabel}>Posts</p>
          </span>
          <span>
            <p className={style.statNumber}>{followers}</p>
            <p className={style.statLabel}>Followers</p>
          </span>
          <span>
            <p className={style.statNumber}>{following}</p>
            <p className={style.statLabel}>Following</p>
          </span>
        </div>
        <button className={style.viewProfileButton}>View Profile</button>
      </div>
    </div>
  );
};

export default Card;
