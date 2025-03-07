"use client";

import React, { useEffect, useState } from "react";
import {
  HeartIcon,
  ChatBubbleLeftIcon,
  HandThumbDownIcon,
} from "@heroicons/react/24/outline";
import { HeartIcon as HeartSolidIcon } from "@heroicons/react/24/solid";
import styles from "@/styles/PostsFeed.module.css";
import { invokeAPI } from "@/utils/invokeAPI";

const PostActionButtons = ({ postID }) => {
  const [likes, setLikes] = useState(0);
  const [dislikes, setDislikes] = useState(0);
  const [comments, setComments] = useState(0);
  const [error, setError] = useState("");

  async function fetchPostStats() {
    try {
      const response = await invokeAPI(`postInteractions/${postID}`, {}, "GET");
      if (response.code === 200) {
        console.log(response);
        setLikes(response.data.likes);
        setDislikes(response.data.dislikes);
        setComments(response.data.comments);
      } else {
        setError("Failed to fetch stats");
      }
    } catch (error) {
      console.error("Error fetching post stats:", error);
      setError("Failed to fetch stats");
    }
  }

  async function like() {
    try {
      const response = await invokeAPI(`like/${postID}`, {}, "POST");
      if (response.code === 200) {
        console.log(response);
        setLikes(response.data.likes);
        setDislikes(response.data.dislikes);
        setComments(response.data.comments);
      } else {
        setError("Error liking post");
      }
    } catch (error) {
      console.error("Error liking post:", error);
      setError("Error liking post");
    }
  }

  async function dislike() {
    try {
      const response = await invokeAPI(`dislike/${postID}`, {}, "POST");
      if (response.code === 200) {
        console.log(response);
        setLikes(response.data.likes);
        setDislikes(response.data.dislikes);
        setComments(response.data.comments);
      } else {
        setError("Error disliking post");
      }
    } catch (error) {
      console.error("Error disliking post:", error);
      setError("Error disliking post");
    }
  }

  useEffect(() => {
    fetchPostStats();
  }, []);

  return (
    <div className={styles.actionButtons}>
      <button className={styles.actionButton} onClick={like}>
        <HeartIcon className="w-6 h-6" />
        <span>{likes}</span>
      </button>

      <button className={styles.actionButton} onClick={dislike}>
        <HandThumbDownIcon className="w-6 h-6" />
        <span>{dislikes}</span>
      </button>

      <button className={styles.actionButton}>
        <ChatBubbleLeftIcon className="w-6 h-6" />
        <span>{comments}</span>
      </button>

      {error && <span className={styles.errorText}>{error}</span>}
    </div>
  );
};

export default PostActionButtons;
