"use client";

import React, { useState, useEffect, act } from "react";
import { useParams, useRouter } from "next/navigation";
import { invokeAPI } from "@/utils/invokeAPI";
import styles from "@/styles/ProfilePage.module.css";
import { LockClosedIcon, LockOpenIcon } from "@heroicons/react/24/outline";
import PostsFeed from "@/app/components/Home/PostsFeed";
import { useAuth } from "@/context/AuthContext";
import { useAlert } from "@/app/components/Alerts/PopUp";
import { useWebSocket } from "@/context/Websocket";

import FollowersAndFollowingList from "@/app/components/Profiles/FollowersAndFollowingList";

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
  const [isMyProfile, setIsMyProfile] = useState(false);
  const [followedByMe, setFollowedByMe] = useState(false);
  const [followsMe, setFollowsMe] = useState(false);
  const [isRequestSent, setIsRequestSent] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const { addMessageHandler, sendMessage } = useWebSocket();
  const { userID } = useAuth();

  //Handles stats update when a new post is added
  addMessageHandler("new_post", async () => {
    await fetchStats();
  });

  //Handles stats update when a follow request is sent
  addMessageHandler("update_profile_stats", async () => {
    await fetchStats();
  });

  useEffect(() => {
    fetchUserData();
  }, []);

  //Fetches the stats of the user
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

  //Fetches the user data
  const fetchUserData = async () => {
    console.log(profileID);
    try {
      const response = await invokeAPI(`profile/${profileID}`, null, "GET");
      if (response.code === 200) {
        console.log(response.data);
        setUser(response.data);
        setIsMyProfile(response.data.is_my_profile);
        setFollowedByMe(response.data.is_following_him);
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

  //Toggles the privacy of the profile
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


  //Handles the follow/unfollow request
  const handleFollow = async () => {
    if (isRequestSent) {
      showAlert({
        type: "confirm",
        message: "Are you sure you want to cancel the follow request?",
        action: async () => {
          try {
            const response = await invokeAPI(
              `followRequest/cancel/${profileID}`,
              null,
              "POST"
            );
            if (response.code === 200) {
              setIsRequestSent(false);
              showAlert({
                type: "success",
                message: "Request cancelled successfully",
              });

              sendMessage({
                type: "cancel_follow_request",
                userDetails: {
                  id: "",
                  nickname: "",
                },
                content: "",
                followRequest: {
                  from: userID,
                  to: profileID,
                  senderNickname: "", // This will be set in the backend
                },
              });
              
            } else {
              showAlert({
                type: "error",
                message: response.error_msg || "Failed to cancel request",
              });
            }
          } catch (error) {
            console.error("Failed to cancel request:", error);
            showAlert({
              type: "error",
              message: "Failed to cancel request",
            });
          }
        },
      });
    } else if (followedByMe) {
      showAlert({
        type: "confirm",
        message: `Are you sure you want to unfollow @${user.nickname}?`,
        action: async () => {
          try {
            const response = await invokeAPI(
              `follow/${profileID}`,
              {},
              "DELETE"
            );
            if (response.code === 200) {
              setFollowedByMe(false);
              setUser((prev) => ({
                ...prev,
                num_followers: prev.num_followers - 1,
              }));

              sendMessage({
                type: "update_profile_stats",
                userDetails: {
                  id: userID,
                },
              });

              sendMessage({
                type: "update_follower_list",
                userDetails: {
                  id: userID,
                },
              });
            } else {
              showAlert({
                type: "error",
                message: response.error_msg || "Failed to unfollow user",
              });
            }
          } catch (error) {
            console.error("Failed to unfollow:", error);
            showAlert({
              type: "error",
              message: "Failed to unfollow user",
            });
          }
        },
      });
    } else if (!followedByMe && user.is_private) {
      try {
        const response = await invokeAPI(
          `followRequest/send/${profileID}`,
          {},
          "POST"
        );
        if (response.code === 200) {
          setIsRequestSent(true);
          showAlert({
            type: "success",
            message: "Request sent successfully",
          });

          sendMessage({
            type: "new_follow_request",
            userDetails: {
              id: "",
              nickname: "",
            },
            content: "",
            followRequest: {
              from: userID,
              to: profileID,
              senderNickname: "", // This will be set in the backend
            },
          });
        } else {
          showAlert({
            type: "error",
            message: response.error_msg || "Failed to send request",
          });
        }
      } catch (error) {
        console.error("Failed to send request:", error);
        showAlert({
          type: "error",
          message: "Failed to send follow request",
        });
      }
    } else {
      try {
        const response = await invokeAPI(`follow/${profileID}`, {}, "POST");
        if (response.code === 200) {
          setFollowedByMe(true);
          setUser((prev) => ({
            ...prev,
            num_followers: prev.num_followers + 1,
          }));

          sendMessage({
            type: "update_profile_stats",
            userDetails: {
              id: userID,
            },
          });

          sendMessage({
            type: "update_follower_list",
            userDetails: {
              id: userID,
            },
          });
        } else {
          showAlert({
            type: "error",
            message: response.error_msg || "Failed to follow user",
          });
        }
      } catch (error) {
        console.error("Failed to follow:", error);
        showAlert({
          type: "error",
          message: "Failed to follow user",
        });
      }
    }
  };

  if (loading) return <div className={styles.loading}>Loading...</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div
      className={`${styles.container} ${
        user.is_private && !followedByMe && !isMyProfile
          ? styles.privateContainer
          : ""
      }`}
    >
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
                {(!user.is_private || followedByMe || isMyProfile) && (
                  <h2 className={styles.name}>
                    {user.first_name} {user.last_name}
                  </h2>
                )}
              </div>
            </div>

            {(!user.is_private || followedByMe || isMyProfile) && (
              <>
                <div className={styles.details}>
                  <p className={styles.detailItem}>
                    <span className={styles.detailLabel}>Email:</span>{" "}
                    {user.email}
                  </p>
                  <p className={styles.detailItem}>
                    <span className={styles.detailLabel}>Date of Birth:</span>{" "}
                    {user.dob}
                  </p>
                </div>

                {user.aboutMe && (
                  <div className={styles.aboutMe}>
                    <h3 className={styles.aboutMeTitle}>About Me</h3>
                    <p>{user.aboutMe}</p>
                  </div>
                )}
              </>
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

            <div className={styles.actions}>
              {isMyProfile ? (
                <button
                  onClick={togglePrivacy}
                  className={`${styles.actionButton} ${styles.editButton}`}
                >
                  Change Privacy
                </button>
              ) : (
                <div className={styles.actionButtons}>
                  <button
                    onClick={handleFollow}
                    className={`${styles.actionButton} ${
                      isRequestSent
                        ? styles.requestedButton
                        : followedByMe
                        ? styles.unfollowButton
                        : styles.followButton
                    }`}
                  >
                    {isRequestSent
                      ? "Requested"
                      : followedByMe
                      ? "Following"
                      : "Follow"}
                  </button>
                  {(followedByMe || followsMe) && (
                    <button
                      onClick={() => router.push(`/messages/${profileID}`)}
                      className={`${styles.actionButton} ${styles.messageButton}`}
                    >
                      Message
                    </button>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
      {(!user.is_private ||
        isMyProfile ||
        (user.is_private && followedByMe)) && (
        <>
          <div className={styles.postsFeed}>
            <PostsFeed
              isProfile={true}
              profileID={profileID}
              myProfile={isMyProfile}
            />
          </div>
          <div className={styles.followersFollowing}>
            {" "}
            <FollowersAndFollowingList
              profileID={profileID}
              myProfile={isMyProfile}
              isPrivate={user.is_private}
            />{" "}
          </div>
        </>
      )}
    </div>
  );
};

export default ProfilePage;
