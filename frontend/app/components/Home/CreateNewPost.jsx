import React, { useEffect, useState } from "react";
import styles from "@/styles/CreateNewPost.module.css";
import { invokeAPI } from "@/utils/invokeAPI";
import FailAlert from "../Alerts/FailAlert";
import { useWebSocket } from "@/context/Websocket";
import { useAuth } from "@/context/AuthContext";

const CreateNewPost = ({ onClose, onPostCreated, isGroup, groupID }) => {
  // Add onPostCreated prop
  const [postTitle, setPostTitle] = useState("");
  const [postContent, setPostContent] = useState("");
  const [postImage, setPostImage] = useState(null);
  const [postPrivacy, setPostPrivacy] = useState("public");
  const [selectedFollowers, setSelectedFollowers] = useState([]);
  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [followersList, setFollowersList] = useState([]);
  const [fetchError, setFetchError] = useState("");

  const { sendMessage } = useWebSocket();
  const { userID } = useAuth();

  useEffect(() => {
    if (isGroup) {
      return;
    }
    async function getFollowers() {
      setFetchError("");
      try {
        const response = await invokeAPI("followersList", {}, "GET");
        if (response.code === 200) {
          console.log(response);
          const followers = response.data;
          setFollowersList(followers);
        } else {
          setFetchError("Error fetching data");
        }
      } catch (error) {
        console.error("Failed to fetch followers", error);
        setFetchError("Error fetching data");
      }
    }
    getFollowers();
  }, []);

  async function handleSubmit(e) {
    e.preventDefault();
    setIsSubmitting(true); // Add loading state
    const formData = new FormData();
    if (!postTitle || !postContent || fetchError) {
      setErrors({
        title: postTitle ? "" : "Title is required",
        content: postContent ? "" : "Content is required",
      });
      setIsSubmitting(false);
      return;
    }
    formData.append("title", postTitle);
    formData.append("content", postContent);
    formData.append("privacy", postPrivacy);
    if (postImage) {
      formData.append("image", postImage);
    }
    if (postPrivacy === "private") {
      if (selectedFollowers.length === 0) {
        setErrors({ submit: "Select at least one follower" });
        setIsSubmitting(false);
        return;
      }
      formData.append("followers", JSON.stringify(selectedFollowers));
    }

    try {
      let response;
      if (isGroup) {
        formData.append("groupID", groupID);
        response = await invokeAPI("createGroupPost", formData, "POST");
      } else {
        response = await invokeAPI("createPost", formData, "POST");
      }
      console.log(response);
      if (response.code === 200) {
        // Clear form
        setPostTitle("");
        setPostContent("");
        setPostImage(null);
        setPostPrivacy("public");
        setSelectedFollowers([]);

        // Call the onPostCreated callback
        if (onPostCreated) {
          onPostCreated();
        }
        onClose();

        if (!isGroup) {
          sendMessage({
            type: "new_post",
            userDetails: {
              id: userID,
            },
          });
        }
      }
    } catch (error) {
      console.error("Failed to create post", error);
      setErrors({ submit: "Something went wrong. Please try again." });
    } finally {
      setIsSubmitting(false);
    }
  }

  function handleImageChange(e) {
    const file = e.target.files[0];
    const allowedExtensions = ["image/jpeg", "image/png", "image/gif"];
    if (file && !allowedExtensions.includes(file.type)) {
      setErrors((prev) => ({
        ...prev,
        image: "Only JPG, PNG, and GIF files are allowed",
      }));
      setPostImage(null);
    } else {
      setErrors((prev) => ({ ...prev, image: "" }));
      setPostImage(file);
    }
  }

  const handleFollowerSelection = (followerId) => {
    setSelectedFollowers((prev) => {
      if (prev.includes(followerId)) {
        return prev.filter((id) => id !== followerId);
      } else {
        return [...prev, followerId];
      }
    });
  };

  return (
    <div className={styles.container}>
      <div className={styles.modal}>
        {fetchError && <FailAlert msg={fetchError} />}
        <div className={styles.modalHeader}>
          <span className={styles.modalTitle}>Create a Post</span>
          <button
            className={styles.closeButton}
            onClick={onClose}
            aria-label="Close"
          >
            ✖
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.modalBody}>
          <div className={styles.input}>
            <label className={styles.inputLabel}>Post Title</label>
            <input
              type="text"
              name="title"
              value={postTitle}
              onChange={(e) => {
                setPostTitle(e.target.value);
              }}
              required
              className={`${styles.inputField} ${
                errors.title ? styles.errorInput : ""
              }`}
              maxLength="32"
            />
            {errors.title && (
              <span className={styles.errorText}>{errors.title}</span>
            )}
          </div>

          <div className={styles.input}>
            <label className={styles.inputLabel}>Content</label>
            <textarea
              name="content"
              value={postContent}
              onChange={(e) => {
                setPostContent(e.target.value);
              }}
              className={`${styles.inputField} ${
                errors.content ? styles.errorInput : ""
              }`}
              rows="2"
              maxLength={"100"}
              style={{ resize: "none" }}
              required
            />
            {errors.content && (
              <span className={styles.errorText}>{errors.content}</span>
            )}
          </div>

          <div className={styles.input}>
            <label className={styles.inputLabel}>Image (optional)</label>
            <input
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={(e) => {
                handleImageChange(e);
              }}
              className={styles.fileInput}
            />
            {errors.image && (
              <span className={styles.errorText}>{errors.image}</span>
            )}
          </div>

          {!isGroup && (
            <div className={styles.input}>
              <label className={styles.inputLabel}>Post Privacy</label>
              <div className={styles.radioGroup}>
                <label
                  className={styles.radioLabel}
                  aria-checked={postPrivacy === "public"}
                >
                  <input
                    type="radio"
                    value="public"
                    checked={postPrivacy === "public"}
                    onChange={(e) => setPostPrivacy(e.target.value)}
                    name="privacy"
                  />
                  <span>Public</span>
                </label>
                <label
                  className={styles.radioLabel}
                  aria-checked={postPrivacy === "almost_private"}
                >
                  <input
                    type="radio"
                    value="almost_private"
                    checked={postPrivacy === "almost_private"}
                    onChange={(e) => setPostPrivacy(e.target.value)}
                    name="privacy"
                  />
                  <span>Followers</span>
                </label>
                {followersList && followersList.length !== 0 && (
                  <label
                    className={styles.radioLabel}
                    aria-checked={postPrivacy === "private"}
                  >
                    <input
                      type="radio"
                      value="private"
                      checked={postPrivacy === "private"}
                      onChange={(e) => setPostPrivacy(e.target.value)}
                      name="privacy"
                    />
                    <span>Private</span>
                  </label>
                )}
              </div>
            </div>
          )}

          {postPrivacy === "private" && !isGroup && (
            <div className={styles.followerSection}>
              <label className={styles.inputLabel}>Select Followers</label>
              <div className={styles.followerList}>
                {followersList.map((follower) => (
                  <label key={follower.id} className={styles.followerItem}>
                    <input
                      type="checkbox"
                      value={follower.id}
                      checked={selectedFollowers.includes(follower.id)}
                      onChange={() => handleFollowerSelection(follower.id)}
                    />
                    <span className={styles.followerName}>
                      @{follower.nickname}
                    </span>
                  </label>
                ))}
              </div>
            </div>
          )}

          <div className={styles.modalFooter}>
            <button
              type="button"
              className={styles.buttonSecondary}
              onClick={onClose}
            >
              Cancel
            </button>
            <button
              type="submit"
              className={styles.buttonPrimary}
              disabled={isSubmitting}
            >
              {isSubmitting ? "Creating..." : "Create Post"}
            </button>
          </div>

          {errors.submit && (
            <div className={styles.submitError}>{errors.submit}</div>
          )}
        </form>
      </div>
    </div>
  );
};

CreateNewPost.defaultProps = {
  onPostCreated: () => {}, // Add default prop
};

export default CreateNewPost;
