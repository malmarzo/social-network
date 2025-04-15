import React, { useEffect, useState, useCallback } from "react";
import styles from "@/styles/Explore.module.css";
import Link from "next/link";
import { invokeAPI } from "@/utils/invokeAPI";
import UserLoader from "../loaders/UserLoader";
import debounce from "lodash/debounce";
import { PlusCircle } from "lucide-react";
import { useWebSocket } from "@/context/Websocket";
import { sendActiveGroupMessage } from "../../groupChat/groupMessage";
import { sendResetCountMessage } from "../../groupChat/groupMessage";
import { handleRequestJoin } from "../../groupChat/groupMessage";
import  DisplayInvitationCard from "../../createGroup/invitationCard"
import  DisplayRequestCard from "../../requestGroup/RequestCard"
import EventNotificationCard from "@/app/groupChat/[id]/eventNotificationCard";

const Explore = () => {
  const [UsersSearch, setUsersSearch] = useState(false);
  const [usersList, setUsersList] = useState([]);
  const [allGroupList, setAllGroupList] = useState([]);
  const [myGroupList, setMyGroupList] = useState([]);
  const [notMyGroupList, setNotMyGroupList] = useState([]);
  const [searchValue, setSearchValue] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filteredUsers, setFilteredUsers] = useState([]);
  const [filteredGroups, setFilteredGroups] = useState([]);
  const [groupsType, setGroupsType] = useState("myGroups");
  const { addMessageHandler } = useWebSocket();
  const { sendMessage } = useWebSocket();
  const [currentUser, setCurrentUser] = useState([]);
  const [invitations, setInvitations] = useState([]); 
  const [requests, setRequests] = useState([]);
  const [eventNotifications, setEventNotifications] = useState([]);
    

    
  const getMyGroups = async () => {
    setLoading(true);
    const myGroupsMsg = { type: "myGroups" };
    sendMessage(myGroupsMsg);
    
};

const getGroupsToRequest = () => {
  setLoading(true);
  const GroupsToRequestMsg = { type: "groupsToRequest" };
  sendMessage(GroupsToRequestMsg);
};

  //Fetches the users and groups on component mount
  async function fetchData() {
    setError("");
    try {
      const response = await invokeAPI("explore", null, "GET");
      if (response.code === 200) {
        console.log(response.data);
        setUsersList(response.data.users_list || []);
        setAllGroupList(response.data.all_groups_list || []);
        

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
    // test if somthing went wrong
    getMyGroups();
    // end of test
  }

  useEffect(() => {
    fetchData();

  addMessageHandler("myGroups", (msg) => {
    setLoading(false);
    setGroupsType("myGroups");
    console.log("Received myGroups message:", msg);
    const updatedMyGroups = Array.isArray(msg.my_groups) ? msg.my_groups : [];
    setMyGroupList(updatedMyGroups);
  
    if (!UsersSearch && groupsType === "myGroups") {
      setLoading(false);
      setFilteredGroups(updatedMyGroups); 
    }
  });

  addMessageHandler("groupMessage", (msg) => {
      console.log("New message received, refreshing groups...");
      getMyGroups();  
  });

  addMessageHandler("groupsToRequest", (msg) => {
    setLoading(false);
    //setGroupsType("notMyGroups");
    if (!msg.my_groups || msg.my_groups.length === 0) {
      setNotMyGroupList([]); 
    } else {
      setNotMyGroupList(msg.my_groups);
    }
    setCurrentUser(msg.userDetails.id)

    if (!UsersSearch && groupsType === "notMyGroups") {
      setLoading(false);
      setFilteredGroups(notMyGroupList); // 
    }
});

    addMessageHandler("invite", (msg) => {
      setInvitations((prevNotifications) => [...prevNotifications, msg]);
    });

    addMessageHandler("request", (msg) => {
      setRequests((prevNotifications) => [...prevNotifications, msg]);
    });

    addMessageHandler("eventNotificationMsg", (msg) => {
      setEventNotifications((prevNotifications) => [...prevNotifications, msg]);
    });

  }, [addMessageHandler, sendMessage]);

  // Debounced search function will run every 300ms after the user changed the input
  const debouncedSearch = useCallback(
    debounce((searchTerm) => {
      if (searchTerm.trim() === "") {
        //If the inout is empty then show all the list items
        setFilteredUsers(usersList);


       if (groupsType === "all") {
      } else if (groupsType === "myGroups") {
        setFilteredGroups(myGroupList);
      } else {
        setFilteredGroups(notMyGroupList);
      }
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
         
        } else if (groupsType === "myGroups") {
          filtered = myGroupList.filter((group) =>
            group.title.toLowerCase().includes(term)
          );
        } else {
          filtered = notMyGroupList.filter((group) =>
            group.title.toLowerCase().includes(term)
          );
        }
        setFilteredGroups(filtered);
      }
    }, 300),
    [
      usersList,
      UsersSearch,
      groupsType,
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
    } else {
      debouncedSearch(searchValue);
    }
  }, [UsersSearch]);

  useEffect(() => {
    if (!UsersSearch && groupsType === "myGroups") {
      setFilteredGroups(myGroupList);
    }else if(!UsersSearch && groupsType === "notMyGroups"){
      setFilteredGroups(notMyGroupList);

    }
  }, [myGroupList,notMyGroupList, UsersSearch, groupsType]);

  const handleDismissNotification = (index) => {
    setEventNotifications((prevNotifications) =>
      prevNotifications.filter((_, i) => i !== index) 
    );
  };
  


  return (
    <div className={styles.wrapper}>
      <h1 className={styles.title}>Explore</h1>
      <div className={styles.container}>
        <div className={styles.toggleContainer}>
         
          <button
            onClick={() => setUsersSearch(false)}
            className={!UsersSearch ? styles.active : styles.inactive}
          >
            Groups
          </button>

          <button
            onClick={() => setUsersSearch(true)}
            className={UsersSearch ? styles.active : styles.inactive}
          >
            Users
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
                    getMyGroups();
                    
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
                    getGroupsToRequest();
                  }}
                  className={
                    groupsType === "notMyGroups"
                      ? styles.active
                      : styles.inactive
                  }
                >
                  Not My Groups
                </button>
                
                <button
                  onClick={() => {
                     setGroupsType("all");
                  }}
                  className={
                    groupsType === "all" ? styles.active : styles.inactive
                  }
                >
                   🔔
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


              {!UsersSearch && filteredGroups && (
                <>
                  {groupsType === "myGroups" && filteredGroups.map((group) => (
                    <Link
                      key={group.id}
                      href={`/groupChat/${group.id}`}
                      onClick={() => {
                        sendActiveGroupMessage("true", group.id, sendMessage);
                        sessionStorage.setItem("navigatedForwardToGroup", group.id);
                        sendResetCountMessage(group.id, sendMessage);
                      }}
                    >
                      <div className={styles.userCard}>
                        <img
                          src="/imgs/defaultAvatar.jpg"
                          alt={group.title}
                          className={styles.userImage}
                        />
                        <span className={styles.userName}>{group.title}</span>
                        {group.count > 0 && (
                          <span className={styles.groupCount}>{group.count}</span>
                        )}
                      </div>
                    </Link>
                  ))}

                  {groupsType === "notMyGroups" && filteredGroups.map((group) => (
                    <div key={group.id}>
                      <div className={styles.userCard}>
                        <img
                          src="/imgs/defaultAvatar.jpg"
                          alt={group.title}
                          className={styles.userImage}
                        />
                        <span className={styles.userName}>{group.title}</span>
                        <button
                          className={styles.joinButton}
                          onClick={() => {
                            handleRequestJoin(group.id, group.creator_id, currentUser, sendMessage);
                            getGroupsToRequest();
                          }}
                        >
                          Join
                        </button>
                      </div>
                    </div>
                  ))}

                    {groupsType === "all" ? (
                      <div className={styles.groupNotification}>
                      {/* this id for the invitation card */}
                      <>

                    {invitations && invitations.map((invitation, index) => (
                      <DisplayInvitationCard
                        key={index}
                        invitation={invitation}
                        onRespond={(userId, accepted) => {
                          console.log(`User ${userId} ${accepted ? "accepted" : "declined"} the invitation`);

                          // Remove only the responded invitation
                          setInvitations((prev) => prev.filter((_, i) => i !== index));
                        }}
                      />
                    ))} 
                        </>
                        
                        <>
                        {/* this for the request card  */}
                        {requests && requests.map((request, index) => (
                        <DisplayRequestCard
                          key={index}
                          request={request}
                          onRespond={(userId, accepted) => {
                            console.log(`User ${userId} ${accepted ? "accepted" : "declined"} the request`);

                            // Remove only the responded request
                            setRequests((prev) => prev.filter((_, i) => i !== index));
                          }}
                        />
                      ))}
                        </>

                       {/* this for the event notification  */}
                        {eventNotifications.map((notification, index) => (
                                <EventNotificationCard
                                  key={index}
                                  content={notification.content}
                                  onDismiss={() => handleDismissNotification(index)} // Pass index for dismissal
                                />
                              ))}
                        <>
                       
                        </>
                
                    </div>
                
                    ) : null}
                </>
              )}




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
