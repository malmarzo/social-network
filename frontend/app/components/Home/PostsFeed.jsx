import React, { useState } from "react";
import styles from "@/styles/PostsFeed.module.css";
import {
  HeartIcon,
  ChatBubbleLeftIcon,
  HandThumbDownIcon,
} from "@heroicons/react/24/outline";
import { HeartIcon as HeartSolidIcon } from "@heroicons/react/24/solid";

const PostsFeed = () => {
  const [activeTab, setActiveTab] = useState("latest");
  const [expandedPost, setExpandedPost] = useState(null);
  const [likes, setLikes] = useState({});

  const toggles = [
    { id: "latest", label: "Latest Posts" },
    { id: "oldest", label: "Oldest Posts" },
    { id: "trending", label: "Trending" },
    { id: "my-posts", label: "My Posts" },
  ];

  // Example posts data
  const posts = [
    {
      id: 1,
      image: "https://www.w3schools.com/howto/img_avatar.png",
      title: "Amazing sunset",
      caption: "Captured this beautiful moment...",
      comments: [],
    },
    {
      id: 2,
      image: "https://www.w3schools.com/howto/img_avatar.png",
      title: "Amazing sunset",
      caption: "Captured this beautiful moment...",
      comments: [],
    },
    {
      id: 3,
      image: "https://www.w3schools.com/howto/img_avatar.png",
      title: "Amazing sunset",
      caption: "Captured this beautiful moment...",
      comments: [],
    },
    {
      id: 4,
      image: "https://www.w3schools.com/howto/img_avatar.png",
      title: "Amazing sunset",
      caption: "Captured this beautiful moment...",
      comments: [],
    },
    {
      id: 5,
      image: "https://www.w3schools.com/howto/img_avatar.png",
      title: "Amazing sunset",
      caption: "Captured this beautiful moment...",
      comments: [],
    },
    // Add more posts as needed
  ];

  const handleLike = (postId) => {
    setLikes((prev) => ({
      ...prev,
      [postId]: !prev[postId],
    }));
  };

  const toggleComments = (postId) => {
    setExpandedPost(expandedPost === postId ? null : postId);
  };

  return (
    <div className={styles.feedContainer}>
      <h1 className={styles.title}>Posts</h1>

      <nav className={styles.toggleNav}>
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
      </nav>

      <div className={styles.postsGrid}>
        {posts.map((post) => (
          <div key={post.id} className={styles.postCard}>
            <img
              src={post.image}
              alt={post.title}
              className={styles.postImage}
            />

            <div className={styles.postContent}>
              <h2 className={styles.postTitle}>{post.title}</h2>
              <p className={styles.postCaption}>{post.caption}</p>
            </div>

            <div className={styles.actionButtons}>
              <button
                className={styles.actionButton}
                onClick={() => handleLike(post.id)}
              >
                {likes[post.id] ? (
                  <HeartSolidIcon className="w-6 h-6 text-red-500" />
                ) : (
                  <HeartIcon className="w-6 h-6" />
                )}
                <span>Like</span>
              </button>

              <button className={styles.actionButton}>
                <HandThumbDownIcon className="w-6 h-6" />
                <span>Dislike</span>
              </button>

              <button
                className={styles.actionButton}
                onClick={() => toggleComments(post.id)}
              >
                <ChatBubbleLeftIcon className="w-6 h-6" />
                <span>Comment</span>
              </button>
            </div>

            <div
              className={`${styles.commentsSection} ${
                expandedPost === post.id ? styles.expanded : ""
              }`}
            >
              <div className={styles.commentsList}>
                {post.comments.map((comment, index) => (
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
