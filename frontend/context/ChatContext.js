"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
} from "react";
import { useAuth } from "./AuthContext";
import { useWebSocket } from "./Websocket";
import { useNotification } from "./NotificationContext";
import { useChatPage } from "./ChatPageContext";
import { invokeAPI } from "@/utils/invokeAPI";

const ChatContext = createContext();

export function ChatProvider({ children, initialUrlUserId = null }) {
  const { isLoggedIn, loading, userID, userNickname } = useAuth();
  const { addMessageHandler, sendMessage } = useWebSocket();
  const { showSuccess, showError, showInfo, showWarning } = useNotification();
  const { isOnChatPage } = useChatPage();

  // Chat state
  const [chatUsers, setChatUsers] = useState([]);
  const [allUsers, setAllUsers] = useState([]);
  const [onlineUsers, setOnlineUsers] = useState([]);
  const [selectedUser, setSelectedUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [showAllUsers, setShowAllUsers] = useState(true);
  // Initialize unreadMessages from localStorage if available
  const [unreadMessages, setUnreadMessages] = useState(() => {
    if (typeof window !== "undefined") {
      try {
        const saved = localStorage.getItem("unreadMessages");
        return saved ? JSON.parse(saved) : {};
      } catch (e) {
        console.error("Failed to load unread messages from localStorage:", e);
        return {};
      }
    }
    return {};
  });
  const [chatHistory, setChatHistory] = useState([]);
  const [hasMoreMessages, setHasMoreMessages] = useState(true);
  const [currentPage, setCurrentPage] = useState(0);

  // Helper function to save unread messages to localStorage
  const saveUnreadMessages = useCallback((messages) => {
    if (typeof window !== "undefined") {
      try {
        localStorage.setItem("unreadMessages", JSON.stringify(messages));
      } catch (e) {
        console.error("Failed to save unread messages to localStorage:", e);
      }
    }
  }, []);

  // Fetch chat users
  const fetchChatUsers = useCallback(async () => {
    if (!isLoggedIn) return;

    try {
      setIsLoading(true);

      // Get users with their last message timestamps in a single API call
      const response = await invokeAPI("chat/users", null, "GET");

      if (response && response.status === "Success") {
        let users = Array.isArray(response.data) ? response.data : [];

        // Filter out the logged-in user
        users = users.filter((user) => String(user.user_id) !== String(userID));

        // Sort users by recent activity (based on last_activity if available)
        const sortedUsers = [...users].sort((a, b) => {
          // If both users have last_activity, sort by most recent
          if (a.last_activity && b.last_activity) {
            return new Date(b.last_activity) - new Date(a.last_activity);
          }

          // If only one has last_activity, prioritize that one
          if (a.last_activity && !b.last_activity) return -1;
          if (!a.last_activity && b.last_activity) return 1;

          // If neither has activity, keep original order
          return 0;
        });

        // Add hasHistory flag based on last_activity and ensure last_message_time is set
        const usersWithFlags = sortedUsers.map((user) => ({
          ...user,
          hasHistory: !!user.last_activity,
          lastMessageTime: user.last_activity
            ? new Date(user.last_activity).getTime()
            : 0,
          last_message_time:
            user.last_message_time || user.last_activity || null,
        }));

        setChatUsers(usersWithFlags);
      } else {
        const errorMsg = response?.error_msg || "Unknown error";
        console.error("Failed to fetch chat users:", errorMsg);
        showError(`Couldn't load chat users: ${errorMsg}`);
        setChatUsers([]);
      }
    } catch (error) {
      console.error("Error fetching chat users:", error);
      showError("Network error while loading chat users");
      setChatUsers([]);
    } finally {
      setIsLoading(false);
    }
  }, [isLoggedIn]);

  // Fetch eligible chat users (users who follow you or whom you follow)
  const fetchEligibleChatUsers = useCallback(async () => {
    if (!isLoggedIn) return;

    try {
      setIsLoading(true);
      const response = await invokeAPI("chat/eligible-users", null, "GET");

      if (response && response.status === "Success") {
        const users = Array.isArray(response.data) ? response.data : [];
        // Filter out the logged-in user
        const filteredUsers = users.filter(
          (user) => String(user.id) !== String(userID)
        );
        // Format users to match the expected structure
        const formattedUsers = filteredUsers.map((user) => ({
          user_id: user.id,
          nickname: user.nickname,
          avatar: user.avatar || null,
          hasHistory: false,
        }));
        setAllUsers(formattedUsers);
      } else {
        console.error(
          "Failed to fetch eligible chat users:",
          response?.error_msg || "Unknown error"
        );
        setAllUsers([]);
      }
    } catch (error) {
      console.error("Error fetching eligible chat users:", error);
      setAllUsers([]);
    } finally {
      setIsLoading(false);
    }
  }, [isLoggedIn, userID]);

  // Fetch online users
  const fetchOnlineUsers = useCallback(async () => {
    if (!isLoggedIn) return;

    try {
      const response = await invokeAPI("chat/online", null, "GET");

      if (response && response.status === "Success") {
        setOnlineUsers(Array.isArray(response.data) ? response.data : []);
      } else {
        console.error(
          "Failed to fetch online users:",
          response?.error_msg || "Unknown error"
        );
        setOnlineUsers([]);
      }
    } catch (error) {
      console.error("Error fetching online users:", error);
      setOnlineUsers([]);
    }
  }, [isLoggedIn]);

  // Fetch all user statuses with last seen timestamps
  const fetchAllUserStatus = useCallback(async () => {
    if (!isLoggedIn) return;

    try {
      const response = await invokeAPI("chat/all-status", null, "GET");

      if (response && response.status === "Success") {
        // Create a map of user IDs to their status information
        const userStatusMap = {};
        if (Array.isArray(response.data)) {
          response.data.forEach((status) => {
            userStatusMap[status.user_id] = {
              isOnline: status.is_online,
              lastSeen: status.last_seen,
            };
          });
        }

        // Update the user list with status information
        setAllUsers((prevUsers) => {
          return prevUsers.map((user) => ({
            ...user,
            isOnline: userStatusMap[user.user_id]?.isOnline || false,
            lastSeen: userStatusMap[user.user_id]?.lastSeen || 0,
          }));
        });

        // Also update chat users with the same information
        setChatUsers((prevUsers) => {
          return prevUsers.map((user) => ({
            ...user,
            isOnline: userStatusMap[user.user_id]?.isOnline || false,
            lastSeen: userStatusMap[user.user_id]?.lastSeen || 0,
          }));
        });
      } else {
        console.error(
          "Failed to fetch user statuses:",
          response?.error_msg || "Unknown error"
        );
      }
    } catch (error) {
      console.error("Error fetching user statuses:", error);
    }
  }, [isLoggedIn]);

  // Fetch chat history
  const fetchChatHistory = useCallback(
    async (pageNum = 0) => {
      if (!selectedUser || !isLoggedIn) return;

      try {
        setIsLoading(true);
        const limit = 20;
        const offset = pageNum * limit;

        // Make sure we're passing the user ID as a string
        const otherUserId = String(selectedUser.user_id);

        const response = await invokeAPI(
          `chat/history?otherUserId=${otherUserId}&limit=${limit}&offset=${offset}`,
          null,
          "GET"
        );

        if (response && response.code === 200 && response.data) {
          // Ensure we always have an array, even if data is null or undefined
          const messages = Array.isArray(response.data.messages)
            ? response.data.messages
            : [];

          if (messages.length < limit) {
            setHasMoreMessages(false);
          }

          // Process messages to ensure they have the correct format
          const processedMessages = messages.map((msg, index) => ({
            ...msg,
            // Ensure unique IDs by adding page number and index if ID is missing
            id: msg.id
              ? String(msg.id)
              : `msg_${Date.now()}_${pageNum}_${index}_${Math.random()
                  .toString(36)
                  .substr(2, 9)}`,
            sender_id: msg.sender_id || "0", // Keep as string
            receiver_id: msg.receiver_id || "0", // Keep as string
            created_at: msg.created_at || new Date().toISOString(),
            message: msg.message || "",
            sender_name: msg.sender_name || "",
            pending: false,
          }));

          if (pageNum === 0) {
            setChatHistory(processedMessages);
          } else {
            // Ensure no duplicate IDs when merging arrays
            setChatHistory((prev) => {
              // Create a Set of existing IDs for fast lookup
              const existingIds = new Set(prev.map((msg) => msg.id));

              // Filter out any messages that already exist in the current chat history
              const uniqueNewMessages = processedMessages.filter(
                (msg) => !existingIds.has(msg.id)
              );

              return [...uniqueNewMessages, ...prev];
            });
          }

          setCurrentPage(pageNum);
        } else {
          const errorMsg = response?.error_msg || "Unknown error";
          console.error("Failed to fetch chat history:", errorMsg);
          showError(`Couldn't load chat history: ${errorMsg}`);

          // If there's no chat history yet, just set an empty array
          if (pageNum === 0) {
            setChatHistory([]);
          }
        }
      } catch (error) {
        console.error("Error fetching chat history:", error);
      } finally {
        setIsLoading(false);
      }
    },
    [selectedUser, isLoggedIn]
  );

  // Send a message
  const sendChatMessage = useCallback(
    (content) => {
      if (!selectedUser || !content.trim()) return false;

      const trimmedMessage = content.trim();

      // Create message object
      const messageId = `msg_${Date.now()}_${Math.random()
        .toString(36)
        .substr(2, 9)}`;
      const msg = {
        type: "chat",
        userDetails: {
          id: String(userID),
          nickname: userNickname,
        },
        receiverId: String(selectedUser.user_id),
        content: trimmedMessage,
        messageId: messageId,
      };

      // Create a temporary message object for optimistic UI update
      const tempMessage = {
        id: msg.messageId,
        sender_id: String(userID),
        receiver_id: String(selectedUser.user_id),
        message: trimmedMessage,
        created_at: new Date().toISOString(),
        sender_name: userNickname,
        pending: true,
      };

      // Add message to chat history immediately for better UX
      setChatHistory((prev) => [...prev, tempMessage]);

      // Update the selected user's last_message_time and last_message
      setChatUsers((prev) => {
        const now = new Date().toISOString();
        const updatedUsers = prev.map((user) => {
          if (user.user_id === selectedUser.user_id) {
            return {
              ...user,
              last_message_time: now,
              last_message: trimmedMessage,
            };
          }
          return user;
        });

        // Re-sort users by recency
        return updatedUsers.sort((a, b) => {
          // First sort by unread messages
          const aUnread = unreadMessages[a.user_id] || 0;
          const bUnread = unreadMessages[b.user_id] || 0;

          if (aUnread !== bUnread) {
            return bUnread - aUnread; // Users with unread messages first
          }

          // Check if users have message history
          const aHasMessage = Boolean(a.last_message);
          const bHasMessage = Boolean(b.last_message);

          // If one has a message and the other doesn't, prioritize the one with a message
          if (aHasMessage !== bHasMessage) {
            return aHasMessage ? -1 : 1; // User with message comes first
          }

          // Then by recency
          const aTime = a.last_message_time
            ? new Date(a.last_message_time).getTime()
            : 0;
          const bTime = b.last_message_time
            ? new Date(b.last_message_time).getTime()
            : 0;
          return bTime - aTime; // Most recent first
        });
      });

      // Also update the selected user to have the latest message
      setSelectedUser((prev) => ({
        ...prev,
        last_message_time: new Date().toISOString(),
        last_message: trimmedMessage,
      }));

      // Send message via WebSocket
      const sent = sendMessage(msg);

      // Show error notification if message couldn't be sent
      // if (!sent && !isConnected) {
      //   showWarning('Message will be sent when connection is restored', {
      //     duration: 3000,
      //     position: 'bottom-right'
      //   });
      // }

      return sent;
    },
    [selectedUser, userID, userNickname, sendMessage]
  );

  // Toggle between all users and chat users
  const toggleUserView = useCallback(() => {
    setShowAllUsers((prev) => !prev);
  }, []);

  // Select a user to chat with
  const selectUser = useCallback(
    (user) => {
      // Check if we're selecting the same user (refresh case)
      const isSameUser = selectedUser && selectedUser.user_id === user.user_id;

      // Create an updated user object, preserving the original last_message_time
      const updatedUser = {
        ...user,
        // Don't update the timestamp when selecting a user to maintain proper order
        hasHistory: true, // Mark as having history
        unread: 0, // Always reset unread count when selecting a user
      };

      // Always update the selected user to trigger a refresh
      setSelectedUser(updatedUser);

      // Clear unread messages for this user
      if (unreadMessages[user.user_id] && unreadMessages[user.user_id] > 0) {
        console.log(
          `Clearing unread messages for ${user.nickname} (${user.user_id})`
        );

        setUnreadMessages((prev) => {
          const updated = { ...prev };
          delete updated[user.user_id]; // Remove this user's unread count completely

          // Save to localStorage using our helper function
          saveUnreadMessages(updated);

          return updated;
        });

        // Update chat users list to reflect the change in unread count
        setChatUsers((prevUsers) => {
          return prevUsers.map((u) => {
            if (u.user_id === user.user_id) {
              return { ...u, unread: 0 };
            }
            return u;
          });
        });
      }

      // Reset chat history state
      setChatHistory([]);
      setCurrentPage(0);
      setHasMoreMessages(true);

      // Update the user in the chatUsers list to maintain sorting
      setChatUsers((prev) => {
        // First update the selected user
        let updatedUsers = prev.map((u) =>
          u.user_id === user.user_id
            ? updatedUser
            : u.user_id !== user.user_id
            ? {
                ...u,
                unread: unreadMessages[u.user_id] || 0, // Make sure all users have their unread count
              }
            : u
        );

        // Then resort the list to ensure users with unread messages stay at the top
        updatedUsers = updatedUsers.sort((a, b) => {
          // First sort by unread messages (using the unread property we added)
          const aUnread = a.unread || unreadMessages[a.user_id] || 0;
          const bUnread = b.unread || unreadMessages[b.user_id] || 0;

          if (bUnread !== aUnread) {
            return bUnread - aUnread; // Users with unread messages first
          }

          // Check if users have message history
          const aHasMessage = Boolean(a.last_message);
          const bHasMessage = Boolean(b.last_message);

          // If one has a message and the other doesn't, prioritize the one with a message
          if (aHasMessage !== bHasMessage) {
            return aHasMessage ? -1 : 1; // User with message comes first
          }

          // Then by recency
          const aTime = a.last_message_time
            ? new Date(a.last_message_time).getTime()
            : 0;
          const bTime = b.last_message_time
            ? new Date(b.last_message_time).getTime()
            : 0;
          return bTime - aTime; // Most recent first
        });

        return updatedUsers;
      });

      // If it's the same user, we still want to refresh the chat history
      if (isSameUser) {
        // Small delay to ensure state updates have processed
        setTimeout(() => {
          fetchChatHistory(0);
        }, 50);
      }
    },
    [selectedUser, fetchChatHistory]
  );

  // Monitor WebSocket connection status
  useEffect(() => {
    if (!isLoggedIn) return;

    // When connection status changes, update UI accordingly
    // console.log('WebSocket connection status:', isConnected ? 'Connected' : 'Disconnected');

    // Show notification when connection status changes
    // if (isConnected) {
    //   showSuccess('Connected to chat server', { duration: 3000 });
    // } else {
    //   showWarning('Disconnected from chat server. Trying to reconnect...', {
    //     duration: 0, // Persistent until reconnected
    //     position: 'top-center'
    //   });
    // }
  }, [isLoggedIn, showSuccess, showWarning]);

  // Initialize data when user logs in
  useEffect(() => {
    if (isLoggedIn && !loading) {
      // Fetch data in sequence to ensure we have the most up-to-date information
      const initializeData = async () => {
        await fetchChatUsers();
        await fetchEligibleChatUsers();
        await fetchOnlineUsers();
        await fetchAllUserStatus();

        // After fetching data, ensure users with unread messages are sorted to the top
        // and add unread count to each user object for easier sorting
        setChatUsers((prevUsers) => {
          return [...prevUsers]
            .map((user) => ({
              ...user,
              unread: unreadMessages[user.user_id] || 0,
            }))
            .sort((a, b) => {
              // First sort by unread messages
              if (a.unread !== b.unread) {
                return b.unread - a.unread; // Users with unread messages first
              }

              // Then by recency of last message
              const aTime = a.last_message_time
                ? new Date(a.last_message_time).getTime()
                : 0;
              const bTime = b.last_message_time
                ? new Date(b.last_message_time).getTime()
                : 0;
              return bTime - aTime; // Most recent first
            });
        });
      };

      initializeData();

      // Set up intervals to refresh online users and user statuses
      const onlineUsersInterval = setInterval(fetchOnlineUsers, 10000);
      const userStatusInterval = setInterval(fetchAllUserStatus, 15000);

      return () => {
        clearInterval(onlineUsersInterval);
        clearInterval(userStatusInterval);
      };
    }
  }, [
    isLoggedIn,
    loading,
    fetchChatUsers,
    fetchEligibleChatUsers,
    fetchOnlineUsers,
    fetchAllUserStatus,
    unreadMessages,
  ]);

  // Handle URL parameters for direct navigation to a chat
  useEffect(() => {
    if (!isLoggedIn || !chatUsers.length || !initialUrlUserId) return;

    // Find the user in the chat users list
    const userFromURL = chatUsers.find(
      (user) => String(user.user_id) === initialUrlUserId
    );
    if (
      userFromURL &&
      (!selectedUser || selectedUser.user_id !== userFromURL.user_id)
    ) {
      setSelectedUser(userFromURL);
    }
  }, [initialUrlUserId, chatUsers, isLoggedIn, selectedUser, setChatUsers]);

  // Helper function to truncate message for preview (moved to top level)

  // Handle incoming chat messages
  useEffect(() => {
     const handleChatMessage = (msg) => {
      // Only process chat messages from other users
      if (
        msg.type === "chat" &&
        msg.userDetails &&
        msg.userDetails.id !== String(userID)
      ) {
        const senderName = msg.userDetails.nickname;
        const senderId = msg.userDetails.id;

        // Always increment unread count if not viewing that user's messages or not on chat page
        // This ensures notifications work even when the user is not on the chat page
        const shouldIncrementUnread =
          !isOnChatPage ||
          !selectedUser ||
          senderId !== String(selectedUser.user_id);

        // Store the unread message count in localStorage to persist across page refreshes
        const currentCount = unreadMessages[senderId] || 0;
        const newUnreadCount = shouldIncrementUnread
          ? currentCount + 1
          : currentCount;

        // Always update the user in the chat list, even if we're not incrementing unread count
        // This ensures the user list is always up to date with the latest messages

        // Update state for unread messages
        setUnreadMessages((prev) => {
          const updated = {
            ...prev,
            [senderId]: newUnreadCount,
          };

          // Save to localStorage using our helper function
          saveUnreadMessages(updated);

          // Log for debugging
          console.log(
            `Unread messages for ${senderName} (${senderId}): ${newUnreadCount}`
          );

          return updated;
        });

        // Check if this user is already in our chat users list
        const userExists = chatUsers.some(
          (u) => String(u.user_id) === senderId
        );

        // If the user doesn't exist in our chat list, we need to add them
        // This ensures users who haven't chatted before still show up with unread messages
        if (!userExists) {
          // Try to find the user in allUsers list
          const userFromAllUsers = allUsers.find(
            (u) => String(u.user_id) === senderId
          );

          if (userFromAllUsers) {
            // Add this user to chatUsers with the new message
            setChatUsers((prevUsers) => [
              {
                ...userFromAllUsers,
                last_message_time: new Date().toISOString(),
                hasHistory: true,
                last_message: msg.content,
                unread: newUnreadCount,
              },
              ...prevUsers,
            ]);
          }
        } else {
          // Move this user to the top of the chat list and update all users with their unread counts
          setChatUsers((prevUsers) => {
            // Find the user
            const userIndex = prevUsers.findIndex(
              (u) => String(u.user_id) === senderId
            );
            if (userIndex === -1) return prevUsers; // User not found (shouldn't happen at this point)

            // Create a new array with the user moved to the top
            const newUsers = [...prevUsers];
            const [user] = newUsers.splice(userIndex, 1);

            // Update the user's lastMessageTime to now and ensure hasHistory is true
            const updatedUser = {
              ...user,
              last_message_time: new Date().toISOString(), // Use consistent property name
              hasHistory: true,
              last_message: msg.content, // Store the last message content
              unread: newUnreadCount, // Add unread count directly to user object for sorting
            };

            // Update all other users with their current unread counts
            const updatedUsers = newUsers.map((u) => ({
              ...u,
              unread: unreadMessages[u.user_id] || 0,
            }));

            // Put the user with the new message at the top, then sort the rest by unread count and recency
            return [updatedUser, ...updatedUsers].sort((a, b) => {
              // Skip the first user (the one with the new message)
              if (a.user_id === updatedUser.user_id) return -1;
              if (b.user_id === updatedUser.user_id) return 1;

              // Sort by unread count first
              if (a.unread !== b.unread) {
                return b.unread - a.unread;
              }

              // Then by recency
              const aTime = a.last_message_time
                ? new Date(a.last_message_time).getTime()
                : 0;
              const bTime = b.last_message_time
                ? new Date(b.last_message_time).getTime()
                : 0;
              return bTime - aTime;
            });
          });
        }

        // // Always show notification for new messages, regardless of page
        // showInfo(`${senderName} sent you a message`, {
        //   duration: 3000,
        //   position: "top-right",
        //   link: "/chat",
        // });

        // Always play notification sound
        // try {
        //   const audio = new Audio("/notification.mp3");
        //   audio.play().catch((e) => console.log("Audio play failed:", e));
        // } catch (e) {
        //   console.log("Audio not supported:", e);
        // }
      }

      // Update chat history if message is from/to the selected user
      if (
        selectedUser &&
        msg.type === "chat" &&
        ((msg.userDetails &&
          msg.userDetails.id === String(selectedUser.user_id)) ||
          msg.receiverId === String(selectedUser.user_id))
      ) {
        const newMessage = {
          id: msg.messageId || `msg_${Date.now()}`,
          sender_id: msg.userDetails ? msg.userDetails.id : userID,
          receiver_id: msg.receiverId,
          message: msg.content,
          created_at: msg.timestamp
            ? new Date(msg.timestamp).toISOString()
            : new Date().toISOString(),
          sender_name: msg.userDetails
            ? msg.userDetails.nickname
            : userNickname,
        };

        setChatHistory((prev) => {
          // Check if this is a confirmation of a pending message
          const pendingMsgIndex = prev.findIndex(
            (m) =>
              m.pending &&
              m.sender_id === newMessage.sender_id &&
              m.receiver_id === newMessage.receiver_id &&
              m.message === newMessage.message
          );

          if (pendingMsgIndex >= 0) {
            // Replace the pending message with the confirmed one
            const updatedHistory = [...prev];
            updatedHistory[pendingMsgIndex] = {
              ...updatedHistory[pendingMsgIndex],
              id: newMessage.id,
              pending: false,
            };
            return updatedHistory;
          } else {
            // Check if we already have this message (avoid duplicates)
            const existingMsg = prev.find(
              (m) =>
                m.id === newMessage.id ||
                (m.message === newMessage.message &&
                  m.sender_id === newMessage.sender_id &&
                  Math.abs(
                    new Date(m.created_at) - new Date(newMessage.created_at)
                  ) < 1000)
            );

            if (existingMsg) {
              return prev; // Don't add duplicate messages
            }

            return [...prev, newMessage];
          }
        });
      }

      // Refresh chat users list to update last message info
      fetchChatUsers();
    };

    // Add message handler for chat messages
    addMessageHandler("chat", handleChatMessage);

    return () => {
      // Cleanup if needed
    };
  }, [addMessageHandler, selectedUser, userID, fetchChatUsers]);

  // Load chat history when selected user changes
  useEffect(() => {
    if (selectedUser) {
      fetchChatHistory(0);
    }
  }, [selectedUser, fetchChatHistory]);

  // Helper functions
  const formatTimestamp = (timestamp) => {
    if (!timestamp) return "";

    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);
    const isToday = date.toDateString() === now.toDateString();
    const isYesterday =
      new Date(now - 86400000).toDateString() === date.toDateString();

    // Just now (less than 1 minute ago)
    if (diffMins < 1) {
      return "Just now";
    }
    // Minutes ago (less than 1 hour)
    else if (diffMins < 60) {
      return `${diffMins}m ago`;
    }
    // Hours ago (less than 24 hours)
    else if (diffHours < 24 && isToday) {
      return `${diffHours}h ago`;
    }
    // Yesterday
    else if (isYesterday) {
      return "Yesterday";
    }
    // This week (less than 7 days)
    else if (diffDays < 7) {
      return date.toLocaleDateString([], { weekday: "short" });
    }
    // Older than a week
    else {
      return date.toLocaleDateString([], { month: "short", day: "numeric" });
    }
  };

  const truncateMessage = (message, maxLength = 30) => {
    if (!message) return "";
    return message.length > maxLength
      ? message.substring(0, maxLength) + "..."
      : message;
  };

  // Provide the context value
  const value = {
    // State
    chatUsers,
    allUsers,
    onlineUsers,
    selectedUser,
    isLoading,
    showAllUsers,
    unreadMessages,
    chatHistory,
    hasMoreMessages,
    currentPage,

    // Actions
    fetchChatUsers,
    fetchEligibleChatUsers,
    fetchOnlineUsers,
    fetchAllUserStatus,
    fetchChatHistory,
    sendChatMessage,
    toggleUserView,
    selectUser,

    // Helpers
    formatTimestamp,
    truncateMessage,
  };

  return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>;
}

export function useChat() {
  const context = useContext(ChatContext);
  if (context === undefined) {
    throw new Error("useChat must be used within a ChatProvider");
  }
  return context;
}
