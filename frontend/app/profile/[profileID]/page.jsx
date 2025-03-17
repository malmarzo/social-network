"use client";

import React, { useState, useEffect, act } from "react";
import { useParams, useRouter } from "next/navigation";
import { invokeAPI } from "@/utils/invokeAPI";
import styles from "@/styles/ProfilePage.module.css";
import { LockClosedIcon, LockOpenIcon } from "@heroicons/react/24/outline";
import PostsFeed from "@/app/components/Home/PostsFeed";
import Explore from "@/app/components/Home/Explore";
import { useAlert } from "@/app/components/Alerts/PopUp";
import { useWebSocket } from "@/context/Websocket";
import { set } from "lodash";

const ProfilePage = () => {
  const { showAlert } = useAlert();
  const router = useRouter();
  const { profileID } = useParams();
  const [user, setUser] = useState({
    id: "",
    nickname: "",
    first_name: "",
    last_name: "",
    email: "",
    dob: "",
    avatar_url: "",
    avatar_mime_type: "",
    aboutMe: "",
    is_private: true,
    num_of_followers: 0,
    num_of_following: 0,
    num_of_posts: 0,
  });
  const [posts, setPosts] = useState([]);
  const [isMyProfile, setIsMyProfile] = useState(false);
  const [followedByMe, setFollowedByMe] = useState(false);
  const [followsMe, setFollowsMe] = useState(false);
  const [imageSrc, setImageSrc] = useState("/imgs/defaultAvatar.jpg");
  const [isRequestSent, setIsRequestSent] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const { addMessageHandler } = useWebSocket();

  addMessageHandler("new_post", async () => {
    await fetchStats();
  });

  useEffect(() => {
    fetchUserData();
    // fetchUserPosts();
  }, []);

  const fetchStats = async () => {
    try {
      const response = await invokeAPI(
        `profileStats/${profileID}`,
        null,
        "GET"
      );
      if (response.code === 200) {
        setUser((prev) => ({
          ...prev,
          num_of_followers: response.data.num_of_followers,
          num_of_following: response.data.num_of_following,
          num_of_posts: response.data.num_of_posts,
        }));
      }
    } catch (error) {
      console.error("Failed to fetch stats:", error);
      setError("Failed to fetch stats");
    }
  };

  const fetchUserData = async () => {
    console.log(profileID);
    try {
      const response = await invokeAPI(`profile/${profileID}`, null, "GET");
      if (response.code === 200) {
        console.log(response.data);
        setUser(response.data);
        setIsMyProfile(response.data.is_my_profile);
        setFollowedByMe(response.data.followed_by_me);
        setFollowsMe(response.data.is_following_me);
        setIsRequestSent(response.data.is_request_sent);
      } else if (
        response.code === 404 &&
        response.error_msg === "user does not exist"
      ) {
        router.push("/");
      }
    } catch (error) {
      setError("Failed to load profile");
    } finally {
      setLoading(false);
    }
  };

  const fetchUserPosts = async () => {
    try {
      const response = await invokeAPI(`posts/${profileID}`, null, "GET");
      if (response.code === 200) {
        setPosts(response.data);
      }
    } catch (error) {
      console.error("Failed to fetch posts:", error);
    }
  };

  const togglePrivacy = async () => {
    console.log(user.is_private);
    showAlert({
      type: "confirm",
      message: `Are you sure you want to make your profile ${
        user.is_private ? "public" : "private"
      }?`,
      action: async () => {
        try {
          const response = await invokeAPI(
            "updatePrivacy",
            {
              is_private: !user.is_private,
            },
            "POST"
          );
          if (response.code === 200) {
            setUser((prev) => ({ ...prev, is_private: !prev.is_private }));
            showAlert({
              type: "success",
              message: `Profile is now ${
                !user.is_private ? "private" : "public"
              }`,
            });
          }
        } catch (error) {
          console.error("Failed to update privacy:", error);
          showAlert({
            type: "error",
            message: "Failed to update privacy settings",
          });
        }
      },
    });
  };

  const handleFollow = async () => {
    try {
      const response = await invokeAPI(`follow/${profileID}`, {}, "POST");
      if (response.code === 200) {
        setFollowedByMe(!followedByMe);
        setUser((prev) => ({
          ...prev,
          num_followers: followedByMe
            ? prev.num_followers - 1
            : prev.num_followers + 1,
        }));
      }
    } catch (error) {
      console.error("Failed to follow/unfollow:", error);
    }
  };

  if (loading) return <div className={styles.loading}>Loading...</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.container}>
      <div className={styles.profileDetails}>
        <div className={styles.profileCard}>
          <div className={styles.profileImage}>
            <img
              src={
                user.avatar_url
                  ? `data:${user.avatar_mime_type};base64,${user.avatar_url}`
                  : "/imgs/defaultAvatar.jpg"
              }
              alt="Profile"
            />
          </div>

          <div className={styles.info}>
            <div className={styles.header}>
              <div className={styles.nameSection}>
                <div className={styles.nicknameRow}>
                  <h1 className={styles.nickname}>@{user.nickname}</h1>
                  {user.is_private ? (
                    <LockClosedIcon
                      className={`${styles.lockIcon} ${styles.private}`}
                    />
                  ) : (
                    <LockOpenIcon className={styles.lockIcon} />
                  )}
                </div>
                <h2 className={styles.name}>
                  {user.first_name} {user.last_name}
                </h2>
              </div>
            </div>

            <div className={styles.details}>
              <p>{user.email}</p>
              <p>{user.dob}</p>
            </div>
            {user.aboutMe && (
              <div className={styles.aboutMe}>
                <p>{user.aboutMe}</p>
              </div>
            )}

            <div className={styles.stats}>
              <div className={styles.statItem}>
                <span className={styles.statNumber}>{user.num_of_posts}</span>
                <span className={styles.statLabel}>Posts</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statNumber}>
                  {user.num_of_followers}
                </span>
                <span className={styles.statLabel}>Followers</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statNumber}>
                  {user.num_of_following}
                </span>
                <span className={styles.statLabel}>Following</span>
              </div>
            </div>

            {isMyProfile ? (
              <button
                onClick={togglePrivacy}
                className={`${styles.actionButton} ${styles.editButton}`}
              >
                Change Privacy
              </button>
            ) : (
              <button
                onClick={handleFollow}
                className={`${styles.actionButton} ${
                  isRequestSent
                    ? styles.requestedButton
                    : followedByMe
                    ? styles.unfollowButton
                    : styles.followButton
                }`}
                disabled={isRequestSent}
              >
                {isRequestSent
                  ? "Requested"
                  : followedByMe
                  ? "Following"
                  : "Follow"}
              </button>
            )}
          </div>
        </div>
      </div>
      <div className={styles.postsFeed}>
        <PostsFeed
          isProfile={true}
          profileID={profileID}
          myProfile={isMyProfile}
        />
      </div>
      <div className={styles.followersFollowing}>
        {" "}
        <Explore />{" "}
      </div>
    </div>
  );
};

export default ProfilePage;
