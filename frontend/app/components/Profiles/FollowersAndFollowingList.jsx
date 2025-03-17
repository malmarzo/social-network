import React, { useState, useEffect, useCallback } from "react";
import Link from "next/link";
import { invokeAPI } from "@/utils/invokeAPI";
import debounce from "lodash/debounce";
import styles from "@/styles/Explore.module.css";
import UserLoader from "../loaders/UserLoader";
import { useAlert } from "@/app/components/Alerts/PopUp";

const FollowersAndFollowingList = ({ profileID, myProfile }) => {
  const [toggleType, setToggleType] = useState("followers");
  const [followersList, setFollowersList] = useState([]);
  const [followingList, setFollowingList] = useState([]);
  const [requestsList, setRequestsList] = useState([]);
  const [searchValue, setSearchValue] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filteredFollowers, setFilteredFollowers] = useState([]);
  const [filteredFollowing, setFilteredFollowing] = useState([]);

  const { showAlert } = useAlert();

  //Fetches the followers and following and pending requets on component mount
  async function fetchData() {
    setError("");
    try {
      const response = await invokeAPI(
        `profileUsersLists/${profileID}`,
        null,
        "GET"
      );
      if (response.code === 200) {
        console.log(response.data);
        setFollowersList(response.data.followers_list);
        setFollowingList(response.data.following_list);
        setRequestsList(response.data.requests_list);
        setFilteredFollowers(response.data.followers_list);
        setFilteredFollowing(response.data.following_list);
        setLoading(false);
      } else {
        setError("Failed to fetch...");
          setLoading(false);
          console.log(response);
      }
    } catch (error) {
      console.error("Failed to fetch users data:", error);
      setError("Failed to fetch...");
    }
  }

  useEffect(() => {
    fetchData();
  }, []);

  // Debounced search function will run every 300ms after the user changed the input
  const debouncedSearch = useCallback(
    debounce((searchTerm) => {
      if (searchTerm.trim() === "") {
        //If the inout is empty then show all the list items
        setFilteredFollowers(followersList);
        setFilteredFollowing(followingList);
        return;
      }

      const term = searchTerm.toLowerCase();

      if (toggleType === "followers") {
        const filtered = followersList.filter((user) =>
          user.nickname.toLowerCase().includes(term)
        );
        setFilteredFollowers(filtered);
      } else if (toggleType === "following") {
        const filtered = followingList.filter((user) =>
          user.nickname.toLowerCase().includes(term)
        );
        setFilteredFollowing(filtered);
      } else {
        setFilteredFollowers(followersList);
        setFilteredFollowing(followingList);
      }
    }, 300),
    [followersList, followingList, toggleType]
  );

  // Update search when input changes
  useEffect(() => {
    debouncedSearch(searchValue);

    // Will clean up the debounced function
    return () => {
      debouncedSearch.cancel();
    };
  }, [searchValue, debouncedSearch]);

  // Update filtered lists when toggle changes
  useEffect(() => {
    if (searchValue.trim() === "") {
      setFilteredFollowers(followersList);
      setFilteredFollowing(followingList);
    } else {
      debouncedSearch(searchValue);
    }
  }, [toggleType]);

  const handleFollowRequest = async (requestID, action) => {
    try {
      const response = await invokeAPI(
        `followRequest/${action}/${requestID}`,
        null,
        "POST"
      );
      if (response.code === 200) {
        // Remove the request from the list
        setRequestsList((prev) => prev.filter((req) => req.id !== requestID));
        showAlert({
          type: "success",
          message: `Request ${action}ed successfully`,
        });
      }
    } catch (error) {
      console.error(`Failed to ${action} request:`, error);
      showAlert({
        type: "error",
        message: `Failed to ${action} request`,
      });
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.toggleContainer}>
        <button
          onClick={() => setToggleType("followers")}
          className={
            toggleType === "followers" ? styles.active : styles.inactive
          }
        >
          Followers
        </button>
        <button
          onClick={() => setToggleType("following")}
          className={
            toggleType === "following" ? styles.active : styles.inactive
          }
        >
          Following
        </button>
        {myProfile && (
          <button
            onClick={() => setToggleType("requests")}
            className={
              toggleType === "requests" ? styles.active : styles.inactive
            }
          >
            Requests
          </button>
        )}
      </div>

      {toggleType !== "requests" && (
        <input
          type="text"
          placeholder={`Search ${toggleType}...`}
          className={styles.searchInput}
          value={searchValue}
          onChange={(e) => setSearchValue(e.target.value)}
        />
      )}

      {loading && !error && (
        <div className={styles.loaderContainer}>
          <UserLoader />
          <UserLoader />
        </div>
      )}

      {!loading && error && <p className={styles.error}>{error}</p>}

      {!loading && !error && (
        <div className={styles.usersList}>
          {toggleType === "followers" &&
            Array.isArray(filteredFollowers) &&
            filteredFollowers.map((user) => (
              <Link key={user.id} href={`/profile/${user.id}`}>
                <div className={styles.userCard}>
                  <img
                    src={user.avatar_url || "/imgs/defaultAvatar.jpg"}
                    alt={user.nickname}
                    className={styles.userImage}
                  />
                  <div className={styles.userInfo}>
                    <h3 className={styles.userName}>@{user.nickname}</h3>
                  </div>
                </div>
              </Link>
            ))}

          {toggleType === "following" &&
            Array.isArray(filteredFollowing) &&
            filteredFollowing.map((user) => (
              <Link key={user.id} href={`/profile/${user.id}`}>
                <div className={styles.userCard}>
                  <img
                    src={user.avatar_url || "/imgs/defaultAvatar.jpg"}
                    alt={user.nickname}
                    className={styles.userImage}
                  />
                  <div className={styles.userInfo}>
                    <h3 className={styles.userName}>@{user.nickname}</h3>
                  </div>
                </div>
              </Link>
            ))}

          {toggleType === "requests" &&
            Array.isArray(requestsList) &&
            requestsList.map((request) => (
              <div key={request.request_id} className={styles.userCard}>
                <img
                  src={"/imgs/defaultAvatar.jpg"}
                  alt={request.nickname}
                  className={styles.userImage}
                />
                <div className={styles.userInfo}>
                  <h3 className={styles.userName}>@{request.nickname}</h3>
                  <div className={styles.requestActions}>
                    <button
                      onClick={() => handleFollowRequest(request.request_id, "accept")}
                      className={styles.acceptButton}
                    >
                      Accept
                    </button>
                    <button
                      onClick={() => handleFollowRequest(request.request_id, "reject")}
                      className={styles.rejectButton}
                    >
                      Reject
                    </button>
                  </div>
                </div>
              </div>
            ))}

          {toggleType !== "requests" &&
            (toggleType === "followers"
              ? filteredFollowers?.length === 0
              : filteredFollowing?.length === 0) && (
              <p className={styles.emptyMessage}>No {toggleType} found...</p>
            )}

          {toggleType === "requests" &&
            (!requestsList?.length || requestsList.length === 0) && (
              <p className={styles.emptyMessage}>No pending requests</p>
            )}
        </div>
      )}
    </div>
  );
};

export default FollowersAndFollowingList;
