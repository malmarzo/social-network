import React, { useEffect, useState, useCallback } from "react";
import styles from "@/styles/Explore.module.css";
import Link from "next/link";
import { invokeAPI } from "@/utils/invokeAPI";
import UserLoader from "../loaders/UserLoader";
import debounce from "lodash/debounce";

const Explore = () => {
  const [UsersSearch, setUsersSearch] = useState(true);
  const [usersList, setUsersList] = useState([]);
  const [groupsList, setGroupsList] = useState([]);
  const [searchValue, setSearchValue] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filteredUsers, setFilteredUsers] = useState([]);
  const [filteredGroups, setFilteredGroups] = useState([]);

  //Fetches the users and groups on component mount
  async function fetchData() {
    setError("");
    try {
      const response = await invokeAPI("explore", null, "GET");
      if (response.code === 200) {
        console.log(response.data);
        setUsersList(response.data.users_list);
        setGroupsList(response.data.groups_list);
        setFilteredUsers(response.data.users_list);
        setFilteredGroups(response.data.groups_list);
        setLoading(false);
      } else {
        setError("Failed to fetch...");
        setLoading(false);
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
      if (searchTerm.trim() === "") { //If the inout is empty then show all the list items
        setFilteredUsers(usersList);
        setFilteredGroups(groupsList);
        return;
      }

      const term = searchTerm.toLowerCase();

      if (UsersSearch) {
        const filtered = usersList.filter((user) =>
          user.nickname.toLowerCase().includes(term)
        );
        setFilteredUsers(filtered);
      } else {
        const filtered = groupsList.filter((group) =>
          group.name.toLowerCase().includes(term)
        );
        setFilteredGroups(filtered);
      }
    }, 300),
    [usersList, groupsList, UsersSearch]
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
      setFilteredUsers(usersList);
      setFilteredGroups(groupsList);
    } else {
      debouncedSearch(searchValue);
    }
  }, [UsersSearch]);

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

        <input
          type="text"
          placeholder={`Search ${UsersSearch ? "users" : "groups"}...`}
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
              {UsersSearch &&
                filteredUsers.length > 0 &&
                filteredUsers.map((user) => (
                  <Link key={user.id} href={`/profile/${user.id}`}>
                    <div className={styles.userCard}>
                      <div>
                        <img
                          src="/imgs/defaultAvatar.jpg"
                          alt="user"
                          className={styles.userImage}
                        />
                      </div>

                      <div>
                        <h3 className={styles.userName}>@{user.nickname}</h3>
                      </div>
                    </div>
                  </Link>
                ))}

              {UsersSearch && filteredUsers.length === 0 && (
                <p>No users found...</p>
              )}

              {!UsersSearch &&
                filteredGroups.length > 0 &&
                filteredGroups.map((group) => (
                  <Link key={group.id} href={`/group/${group.id}`}>
                    <div className={styles.userCard}>
                      <div>
                        <img
                          src="/imgs/defaultAvatar.jpg"
                          alt="user"
                          className={styles.userImage}
                        />
                      </div>

                      <div>
                        <h3 className={styles.userName}>{group.name}</h3>
                      </div>
                    </div>
                  </Link>
                ))}

              {!UsersSearch && filteredGroups.length === 0 && (
                <p>No groups found...</p>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Explore;
