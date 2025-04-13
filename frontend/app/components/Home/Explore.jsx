import React, { useEffect, useState, useCallback } from "react";
import styles from "@/styles/Explore.module.css";
import Link from "next/link";
import { invokeAPI } from "@/utils/invokeAPI";
import UserLoader from "../loaders/UserLoader";
import debounce from "lodash/debounce";
import { PlusCircle } from "lucide-react";

const Explore = () => {
  const [UsersSearch, setUsersSearch] = useState(true);
  const [usersList, setUsersList] = useState([]);
  const [allGroupList, setAllGroupList] = useState([]);
  const [myGroupList, setMyGroupList] = useState([]);
  const [notMyGroupList, setNotMyGroupList] = useState([]);
  const [searchValue, setSearchValue] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filteredUsers, setFilteredUsers] = useState([]);
  const [filteredGroups, setFilteredGroups] = useState([]);
  const [groupsType, setGroupsType] = useState("all");

  //Fetches the users and groups on component mount
  async function fetchData() {
    setError("");
    try {
      const response = await invokeAPI("explore", null, "GET");
      if (response.code === 200) {
        console.log(response.data);
        setUsersList(response.data.users_list || []);
        setAllGroupList(response.data.all_groups_list || []);
        setMyGroupList(response.data.my_groups_list || []);
        setNotMyGroupList(response.data.not_my_groups_list || []);

        // Set the initial filtered lists to the full lists
        setFilteredUsers(response.data.users_list || []);
        setFilteredGroups(response.data.all_groups_list || []);
        setLoading(false);
      } else {
        setError("Failed to fetch...");
        setLoading(false);
      }
    } catch (error) {
      console.error("Failed to fetch users data:", error);
      setError("Failed to fetch...");
      setLoading(false);
      setFilteredUsers([]);
      setFilteredGroups([]);
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
        setFilteredUsers(usersList);
        setFilteredGroups(allGroupList);
        return;
      }

      const term = searchTerm.toLowerCase();

      if (UsersSearch) {
        const filtered = usersList.filter((user) =>
          user.nickname.toLowerCase().includes(term)
        );
        setFilteredUsers(filtered);
      } else {
        let filtered = [];
        if (groupsType === "all") {
          filtered = allGroupList.filter((group) =>
            group.name.toLowerCase().includes(term)
          );
        } else if (groupsType === "myGroups") {
          filtered = myGroupList.filter((group) =>
            group.name.toLowerCase().includes(term)
          );
        } else {
          filtered = notMyGroupList.filter((group) =>
            group.name.toLowerCase().includes(term)
          );
        }
        setFilteredGroups(filtered);
      }
    }, 300),
    [
      usersList,
      UsersSearch,
      groupsType,
      allGroupList,
      myGroupList,
      notMyGroupList,
    ]
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
      setFilteredGroups(allGroupList);
    } else {
      debouncedSearch(searchValue);
    }
  }, [UsersSearch]);

  return (
    <div className={styles.wrapper}>
      <h1 className={styles.title}>Explore</h1>
      <div className={styles.container}>
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

        <div className={styles.searchAndCreateContainer}>
          <input
            type="text"
            placeholder={`Search ${UsersSearch ? "users" : "groups"}...`}
            className={styles.searchInput}
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
          />

          {!UsersSearch && (
            <Link
              href={"/createGroup"}
            >
              <div className={styles.addButtonContainer}>
                <button className={styles.addButton}>
                  <PlusCircle size={20} />
                </button>
                <span className={styles.addButtonText}>New</span>
              </div>
            </Link>
          )}
        </div>

        {!UsersSearch && (
          <div>
            <div className={styles.groupTypeButtonsContainer}>
              <div className={styles.groupTypeButtons}>
                <button
                  onClick={() => {
                    setGroupsType("all");
                    setFilteredGroups(allGroupList);
                  }}
                  className={
                    groupsType === "all" ? styles.active : styles.inactive
                  }
                >
                  All
                </button>
                <button
                  onClick={() => {
                    setGroupsType("myGroups");
                    setFilteredGroups(myGroupList);
                  }}
                  className={
                    groupsType === "myGroups" ? styles.active : styles.inactive
                  }
                >
                  My Groups
                </button>
                <button
                  onClick={() => {
                    setGroupsType("notMyGroups");
                    setFilteredGroups(notMyGroupList);
                  }}
                  className={
                    groupsType === "notMyGroups"
                      ? styles.active
                      : styles.inactive
                  }
                >
                  Not My Groups
                </button>
              </div>
            </div>
          </div>
        )}

        {loading && !error && (
          <div className={styles.loaderContainer}>
            <UserLoader />
            <UserLoader />
          </div>
        )}

        {!loading && error && <p className={styles.noResults}>{error}</p>}

        {!loading && !error && (
          <div className={styles.usersList}>
            {UsersSearch &&
              filteredUsers.map((user) => (
                <Link key={user.id} href={`/profile/${user.id}`}>
                  <div className={styles.userCard}>
                    <img
                      src="/imgs/defaultAvatar.jpg"
                      alt={user.nickname}
                      className={styles.userImage}
                    />
                    <span className={styles.userName}>@{user.nickname}</span>
                  </div>
                </Link>
              ))}

            {UsersSearch && filteredUsers.length === 0 && (
              <p className={styles.noResults}>No users found</p>
            )}

            {!UsersSearch &&
              filteredGroups &&
              filteredGroups.map((group) => (
                <Link key={group.id} href={`/group/${group.id}`}>
                  <div className={styles.userCard}>
                    <img
                      src="/imgs/defaultAvatar.jpg"
                      alt={group.name}
                      className={styles.userImage}
                    />
                    <span className={styles.userName}>{group.name}</span>
                  </div>
                </Link>
              ))}

            {!UsersSearch &&
              (!filteredGroups || filteredGroups.length === 0) && (
                <p className={styles.noResults}>No groups found</p>
              )}
          </div>
        )}
      </div>
    </div>
  );
};

export default Explore;
