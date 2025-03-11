import React, { useState, useCallback, useEffect } from "react";
import styles from "@/styles/PostsFeed.module.css";
import { invokeAPI } from "@/utils/invokeAPI";
import CreateNewPost from "./CreateNewPost";
import Link from "next/link";
import PostActionButtons from "./PostActionButtons";

const PostsFeed = ({isGroup, groupID}) => {
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
      } else {
        response = await invokeAPI("posts", null, "GET");
      }
      if (response.code === 200) {
        setPosts(response.data);
      } else {
        setError("Failed to fetch posts");
      }
    } catch (error) {
      setError("Something went wrong");
      console.error("Error fetching posts:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPosts();
  }, [fetchPosts]);

  const refreshPosts = useCallback(() => {
    fetchPosts();
  }, [fetchPosts]);



  if (loading) {
    return <div className={styles.loading}>Loading posts...</div>;
  }

  if (error) {
    return <div className={styles.error}>{error}</div>;
  }

  return (
    <div className={styles.feedContainer}>
      <h1 className={styles.title}>Posts</h1>
      <nav className={styles.toggleNav}>
      {!isGroup && (
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
          <button
            className={styles.createPostButton}
            onClick={() => setCreateNewPost(true)} // Add click handler
          >
            Create Post
          </button>
      </nav>

      {/* Keep only one instance of CreateNewPost */}
      {createNewPost && (
        <CreateNewPost
          onClose={() => setCreateNewPost(false)}
          onPostCreated={refreshPosts}
          isGroup={isGroup}
          groupID={groupID}
        />
      )}

      <div className={styles.postsGrid}>
        {posts &&
          posts.map((post) => (
            <div key={post.post_id} className={styles.postCard}>
              {post.post_image && (
                <img
                  src={`data:${post.image_mime_type};base64,${post.post_image}`}
                  alt={post.post_title}
                  className={styles.postImage}
                />
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

              <PostActionButtons postID={post.post_id} isGroup={isGroup}/>
            </div>
          ))}
      </div>
    </div>
  );
};

export default PostsFeed;
