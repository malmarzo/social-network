import React, { useEffect, useState } from "react";
import styles from "@/styles/CreateNewPost.module.css";
import { invokeAPI } from "@/utils/invokeAPI";
import FailAlert from "../../components/Alerts/FailAlert"
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
  const [fetchError, setFetchError] = useState("");

  const { sendMessage } = useWebSocket();
  const { userID } = useAuth();

  useEffect(() => {
    if (isGroup) {
      return;
    }
   
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
    //formData.append("privacy", postPrivacy);
    if (postImage) {
      formData.append("image", postImage);
    }
    

    try {
      let response;
      if (isGroup) {
        formData.append("groupID", groupID);
        response = await invokeAPI(`groups/chat/${groupID}/createGroupPost`, formData, "POST");
      } else {
      
      }
      if (response.code === 200) {
        setPostTitle("");
        setPostContent("");
        setPostImage(null);
       

        // Call the onPostCreated callback
        if (onPostCreated) {
          onPostCreated();
        }
        onClose();
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
  onPostCreated: () => {}, 
};

export default CreateNewPost;
