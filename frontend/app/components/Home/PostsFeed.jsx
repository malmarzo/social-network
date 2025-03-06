import React, { useState, useCallback, useEffect } from "react";
import styles from "@/styles/PostsFeed.module.css";
import { invokeAPI } from "@/utils/invokeAPI";
import {
  HeartIcon,
  ChatBubbleLeftIcon,
  HandThumbDownIcon,
} from "@heroicons/react/24/outline";
import { HeartIcon as HeartSolidIcon } from "@heroicons/react/24/solid";
import CreateNewPost from "./CreateNewPost";
import Link from "next/link";

const PostsFeed = () => {
  const [activeTab, setActiveTab] = useState("latest");
  const [expandedPost, setExpandedPost] = useState(null);
  const [likes, setLikes] = useState({});
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
      const response = await invokeAPI("posts", null, "GET");
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

  const handleLike = (postId) => {
    setLikes((prev) => ({
      ...prev,
      [postId]: !prev[postId],
    }));
  };

  const toggleComments = (postId) => {
    setExpandedPost(expandedPost === postId ? null : postId);
  };

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

              <div className={styles.actionButtons}>
                <button
                  className={styles.actionButton}
                  onClick={() => handleLike(post.post_id)}
                >
                  {likes[post.post_id] ? (
                    <HeartSolidIcon className="w-6 h-6 text-red-500" />
                  ) : (
                    <HeartIcon className="w-6 h-6" />
                  )}
                  <span>Like ({post.num_of_likes})</span>
                </button>

                <button className={styles.actionButton}>
                  <HandThumbDownIcon className="w-6 h-6" />
                  <span>Dislike ({post.num_of_dislikes})</span>
                </button>

                <button
                  className={styles.actionButton}
                  onClick={() => toggleComments(post.post_id)}
                >
                  <ChatBubbleLeftIcon className="w-6 h-6" />
                  <span>Comment ({post.num_of_comments})</span>
                </button>
              </div>

              <div
                className={`${styles.commentsSection} ${
                  expandedPost === post.post_id ? styles.expanded : ""
                }`}
              >
                <div className={styles.commentsList}>
                  {post.comments?.map((comment, index) => (
                    <div key={index} className="p-2 border-b border-gray-100">
                      {comment}
                    </div>
                  ))}
                </div>
                <div className={styles.commentInput}>
                  <input
                    type="text"
                    placeholder="Add a comment..."
                    className="w-full p-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  />
                </div>
              </div>
            </div>
          ))}
      </div>
    </div>
  );
};

export default PostsFeed;
