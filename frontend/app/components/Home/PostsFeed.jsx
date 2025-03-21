import React, { useState, useCallback, useEffect } from "react";
import styles from "@/styles/PostsFeed.module.css";
import { invokeAPI } from "@/utils/invokeAPI";
import CreateNewPost from "./CreateNewPost";
import Link from "next/link";
import PostActionButtons from "./PostActionButtons";
import PostLoader from "../loaders/PostLoader";

const PostsFeed = ({ isGroup, groupID, isProfile, profileID, myProfile }) => {
  const [activeTab, setActiveTab] = useState("latest");
  const [createNewPost, setCreateNewPost] = useState(false);
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const toggles = [
    { id: "latest", label: "Latest Posts" },
    { id: "trending", label: "Trending" },
    { id: "my-posts", label: "My Posts" },
  ];

  const fetchPosts = useCallback(async () => {
    try {
      setLoading(true);
      let response;
      if (isGroup) {
        if (!groupID) {
          setError("Failed to fetch posts");
          return;
        }
        response = await invokeAPI(`groupPosts/${groupID}`, null, "GET");
      } else if (isProfile) {
        if (!profileID) {
          setError("Failed to fetch posts");
          return;
        }
        response = await invokeAPI(`profilePosts/${profileID}`, null, "GET");
      } else {
        const queryParams = {
          tab: activeTab,
        };
        response = await invokeAPI("posts", null, "GET", null, queryParams);
      }
      if (response.code === 200) {
        setPosts(response.data);
        console.log(response.data);
      } else {
        setError("Failed to fetch posts");
      }
    } catch (error) {
      setError("Something went wrong");
      console.error("Error fetching posts:", error);
    } finally {
      setLoading(false);
    }
  }, [activeTab]);

  useEffect(() => {
    fetchPosts();
  }, [fetchPosts]);

  const refreshPosts = useCallback(() => {
    setActiveTab("latest");
    fetchPosts();
  }, [fetchPosts]);

  // Show create post button based on context
  const shouldShowCreatePost = () => {
    if (!isGroup && !isProfile) return true; // Home page
    if (isGroup) return true; // Group page
    if (isProfile && myProfile) return true; // My profile page
    return false;
  };


  if (error) {
    return <div className={styles.error}>{error}</div>;
  }

  return (
    <div className={styles.feedContainer}>
      <div className={styles.headerContainer}>
        <h1 className={styles.title}>
          {" "}
          {isGroup || isProfile ? "Activity" : "Posts"}
        </h1>
        <nav className={styles.toggleNav}>
          {!isGroup && !isProfile && (
            <>
              <div className={styles.toggleButtons}>
                {toggles.map((toggle) => (
                  <button
                    key={toggle.id}
                    className={`${styles.toggleButton} ${
                      activeTab === toggle.id ? styles.activeToggle : ""
                    }`}
                    onClick={() => setActiveTab(toggle.id)}
                  >
                    {toggle.label}
                  </button>
                ))}
              </div>
            </>
          )}
          {shouldShowCreatePost() && (
            <button
              className={styles.createPostButton}
              onClick={() => setCreateNewPost(true)}
            >
              Create Post
            </button>
          )}
        </nav>

        {createNewPost && shouldShowCreatePost() && (
          <CreateNewPost
            onClose={() => setCreateNewPost(false)}
            onPostCreated={refreshPosts}
            isGroup={isGroup}
            groupID={groupID}
          />
        )}
      </div>

      <div className={styles.postsGrid}>
        {posts &&
          !loading &&
          posts.map((post) => (
            <div key={post.post_id} className={styles.postCard}>
              {post.post_image && (
                <div className={styles.postImageContainer}>
                  <img
                    src={`data:${post.image_mime_type};base64,${post.post_image}`}
                    alt={post.post_title}
                    className={styles.postImage}
                  />
                </div>
              )}

              <div className={styles.postContent}>
                <div className={styles.postHeader}>
                  <h2 className={styles.postTitle}>{post.post_title}</h2>
                  <div className={styles.authorInfo}>
                    <Link href={`/profile/${post.user_id}`}>
                      <span className={styles.postAuthor}>
                        @{post.user_nickname}
                      </span>
                    </Link>
                    <span className={styles.postDate}>{post.created_at}</span>
                  </div>
                </div>
                <p className={styles.postCaption}>{post.content}</p>
              </div>

              <PostActionButtons postID={post.post_id} isGroup={isGroup} />
            </div>
          ))}

        {loading && !error && <PostLoader />}
      </div>
    </div>
  );
};

export default PostsFeed;
