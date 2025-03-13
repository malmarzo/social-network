import React from "react";
import styles from "@/styles/Explore.module.css";
import { useState } from "react";
import Link from "next/link";
import { invokeAPI } from "@/utils/invokeAPI";
import UserLoader from "../loaders/UserLoader";

const Explore = () => {
  const [UsersSearch, setUsersSearch] = useState(true);
  const [usersList, setUsersList] = useState([]);
  const [groupsList, setGroupsList] = useState([]);
  const [searchValue, setSearchValue] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  //   async function fetchData() {
  //     setError("");
  //     try {
  //       const response = await invokeAPI("explore", null, "GET");
  //       if (response.code === 200) {
  //         setUsersList(response.data.users);
  //         setGroupsList(response.data.groups);
  //         setLoading(false);
  //       }
  //     } catch (error) {
  //       console.error("Failed to fetch users data:", error);
  //       setError("Failed to fetch...");
  //     }
  //   }

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Explore</h1>
      <div>
        <div className={styles.toggleContainer}>
          <button
            onClick={() => setUsersSearch(true)}
            className={UsersSearch ? styles.active : styles.inactive}
          >
            Users
          </button>
          <button
            onClick={() => setUsersSearch(false)}
            className={!UsersSearch ? styles.active : styles.inactive}
          >
            Groups
          </button>
        </div>

        {error && <p className={styles.error}>{error}</p>}

        <input
          type="text"
          placeholder={`Search ${UsersSearch ? "users" : "groups"}`}
          className={styles.searchInput}
          value={searchValue}
          onChange={(e) => setSearchValue(e.target.value)}
        />

        {loading && !error && (
          <div className={styles.loaderContainer}>
            <UserLoader />
            <UserLoader />
          </div>
        )}

        {!loading && error && <p>{error}</p>}

        {!loading && !error && (
          <div>
            <div className={styles.usersList}>
              <Link href="/profile">
                <div className={styles.userCard}>
                  <div>
                    <img
                      src="/imgs/defaultAvatar.jpg"
                      alt="user"
                      className={styles.userImage}
                    />
                  </div>

                  <div>
                    <h3 className={styles.userName}>User Name</h3>
                  </div>
                </div>
              </Link>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Explore;
