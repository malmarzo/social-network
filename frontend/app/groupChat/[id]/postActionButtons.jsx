"use client";

import React, { useEffect, useState } from "react";
import {
  HeartIcon,
  ChatBubbleLeftIcon,
  HandThumbDownIcon,
  PaperAirplaneIcon,
  PhotoIcon,
} from "@heroicons/react/24/outline";
import styles from "@/styles/PostsFeedGroup.module.css";
import { invokeAPI } from "@/utils/invokeAPI";
import Link from "next/link";
import { comment } from "postcss";

const PostActionButtons = ({ postID, isGroup, groupID }) => {
  const [likes, setLikes] = useState(0);
  const [dislikes, setDislikes] = useState(0);
  const [comments, setComments] = useState(0);
  const [expandedPost, setExpandedPost] = useState(null);
  const [commentInput, setCommentInput] = useState("");
  const [commentsList, setCommentsList] = useState([]);
  const [error, setError] = useState("");
  const [commentImage, setCommentImage] = useState(null);
  const [imageError, setImageError] = useState("");
  const [imagePreviewUrl, setImagePreviewUrl] = useState(null);

  async function fetchPostStats() {
    try {
      let response;

      if (isGroup) {
        //console.log("yunes");
        response = await invokeAPI(
          `groups/chat/${groupID}/groupPostInteractions/${postID}`,
          {},
          "GET"
        );
      } else {
        response = await invokeAPI(`groups/chat/${groupID}/postInteractions/${postID}`, {}, "GET");
      }
      if (response.code === 200) {
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
      let response;
      if (isGroup) {
        response = await invokeAPI(`groups/chat/${groupID}/likeGroupPost/${postID}`, {}, "POST");
      } else {
        response = await invokeAPI(`like/${postID}`, {}, "POST");
      }
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
      let response;
      if (isGroup) {
        response = await invokeAPI(`groups/chat/${groupID}/dislikeGroupPost/${postID}`, {}, "POST");
      } else {
        response = await invokeAPI(`dislike/${postID}`, {}, "POST");
      }
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

  const toggleComments = async (postId) => {
    setExpandedPost(expandedPost === postId ? null : postId);
    await getPostComments(expandedPost);
  };

  async function postComment() {
    const formData = new FormData();

    if ((!commentInput || commentInput.trim() === "") && !commentImage) {
      return;
    }

    formData.append("postID", postID);
    formData.append("comment", commentInput);
    if (commentImage) {
      formData.append("image", commentImage);
    }

    try {
      let response;
      if (isGroup) {
        response = await invokeAPI(`groups/chat/${groupID}/groupComment`, formData, "POST");
      } else {
        response = await invokeAPI("comment", formData, "POST");
      }

      if (response.code === 200) {
        console.log(response);
        setLikes(response.data.stats.likes);
        setDislikes(response.data.stats.dislikes);
        setComments(response.data.stats.comments);

        // Initialize commentsList as empty array if null
        setCommentsList((prevComments) => {
          if (!prevComments) {
            return [response.data.comment];
          }
          return [...prevComments, response.data.comment];
        });

        setCommentInput("");
        setCommentImage(null);
        setImagePreviewUrl(null);
      } else {
        setError("Error posting comment");
      }
    } catch (error) {
      console.error("Error posting comment:", error);
      setError("Error posting comment");
    }
  }

  function handleImageChange(e) {
    const file = e.target.files[0];
    const allowedExtensions = ["image/jpeg", "image/png", "image/gif"];

    if (file && !allowedExtensions.includes(file.type)) {
      setImageError("Only JPG, PNG, and GIF files are allowed");
      setCommentImage(null);
      setImagePreviewUrl(null);
    } else {
      setImageError("");
      setCommentImage(file);
      // Create preview URL
      const url = URL.createObjectURL(file);
      setImagePreviewUrl(url);
    }
  }

  async function getPostComments(isOpen) {
    if (isOpen !== postID) {
      try {
        const response = await invokeAPI(`groups/chat/${groupID}/groupComments/${postID}`, {}, "GET");
        if (response.code === 200) {
          console.log("hello yunes",response);
          setCommentsList(response.data);
        }
      } catch (error) {
        console.log("Error fetching comments:", error);
      }
    }
  }
console.log(commentsList);
  // Clean up preview URL when component unmounts
  useEffect(() => {
    return () => {
      if (imagePreviewUrl) {
        URL.revokeObjectURL(imagePreviewUrl);
      }
    };
  }, [imagePreviewUrl]);

  useEffect(() => {
    fetchPostStats();
  }, []);

  return (
    <>
      <div className={styles.actionButtons}>
        <button className={styles.actionButton} onClick={like}>
          <HeartIcon className="w-6 h-6" />
          <span>{likes}</span>
        </button>

        <button className={styles.actionButton} onClick={dislike}>
          <HandThumbDownIcon className="w-6 h-6" />
          <span>{dislikes}</span>
        </button>

        <button
          className={styles.actionButton}
          onClick={() => toggleComments(postID)}
        >
          <ChatBubbleLeftIcon className="w-6 h-6" />
          <span>{comments}</span>
        </button>

        {error && <span className={styles.errorText}>{error}</span>}
      </div>

      <div
        className={`${styles.commentsSection} ${
          expandedPost === postID ? styles.expanded : ""
        }`}
      >
        <div className={styles.commentsList}>
          {commentsList?.map((comment, index) => (
            <div key={index} className={styles.commentItem}>
              <div className={styles.commentHeader}>
                <Link href={`/profile/${comment.user_id}`}>
                  <span className={styles.commentUser}>
                    @{comment.user_nickname}
                  </span>
                </Link>
                <span className={styles.commentDate}>{comment.created_at}</span>
              </div>
              <p className={styles.commentText}>{comment.comment_text}</p>
              {comment.comment_image && (
                <div className={styles.commentImageContainer}>
                  <img
                    src={`data:${comment.image_mime_type};base64,${comment.comment_image}`}
                    alt="Comment attachment"
                    className={styles.commentImage}
                  />
                </div>
              )}
            </div>
          ))}
        </div>

        <div className={styles.commentInputContainer}>
          {imagePreviewUrl && (
            <div className={styles.imagePreviewContainer}>
              <img
                src={imagePreviewUrl}
                alt="Preview"
                className={styles.imagePreview}
              />
              <button
                onClick={() => {
                  setCommentImage(null);
                  setImagePreviewUrl(null);
                }}
                className={styles.removeImageButton}
              >
                ×
              </button>
            </div>
          )}
          <div className={styles.commentInput}>
            <input
              type="text"
              placeholder="Add a comment..."
              onChange={(e) => setCommentInput(e.target.value)}
              value={commentInput}
            />
            <label className={styles.fileInputLabel}>
              <input
                type="file"
                className={styles.fileInput}
                onChange={handleImageChange}
                accept="image/jpeg,image/png,image/gif"
              />
              <PhotoIcon className="w-6 h-6 text-gray-500 hover:text-black transition-colors duration-200" />
            </label>
            <button className={styles.sendButton} onClick={postComment}>
              <PaperAirplaneIcon className="w-6 h-5 block text-gray-500 hover:text-black transition-colors duration-200" />
            </button>
          </div>
        </div>
      </div>
    </>
  );
};

export default PostActionButtons;
