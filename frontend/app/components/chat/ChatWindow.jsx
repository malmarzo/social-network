"use client";

import { useState, useEffect, useRef, useCallback, useMemo } from "react";
import { useChat } from "@/context/ChatContext";
import { useNotification } from "@/context/NotificationContext";
import { useAuth } from "@/context/AuthContext";
import Image from "next/image";
import { formatDistanceToNow } from "date-fns";
import EmojiPicker from "emoji-picker-react";

export default function ChatWindow({ toggleSidebar, isMobile }) {
  const [message, setMessage] = useState("");
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const messagesEndRef = useRef(null);
  const chatContainerRef = useRef(null);
  const messageInputRef = useRef(null);
  const { showSuccess, showError, showInfo } = useNotification();
  const { userID } = useAuth();
  const {
    selectedUser,
    chatHistory,
    isLoading,
    hasMoreMessages,
    currentPage,
    fetchChatHistory,
    sendChatMessage,
    formatTimestamp,
    onlineUsers,
  } = useChat();

  // Group messages by date for better UI organization
  const groupedMessages = useMemo(() => {
    if (!chatHistory || chatHistory.length === 0) return {};

    const groups = {};
    chatHistory.forEach((msg) => {
      const date = new Date(msg.created_at).toLocaleDateString();
      if (!groups[date]) groups[date] = [];
      groups[date].push(msg);
    });

    return groups;
  }, [chatHistory]);

  // Auto-scroll to bottom when chat history changes
  useEffect(() => {
    if (messagesEndRef.current && chatHistory.length > 0 && selectedUser) {
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [chatHistory, selectedUser]);

  // Save scroll position when loading more messages
  const saveScrollPosition = useCallback(() => {
    if (chatContainerRef.current) {
      const container = chatContainerRef.current;
      const scrollHeight = container.scrollHeight;
      const scrollTop = container.scrollTop;
      const clientHeight = container.clientHeight;

      // If user has scrolled up more than 200px, consider it a deliberate scroll up
      return {
        scrollHeight,
        scrollPosition: scrollHeight - scrollTop - clientHeight < 200,
      };
    }
    return { scrollHeight: 0, scrollPosition: true };
  }, []);

  // Load more messages when scrolling to top
  const handleScroll = useCallback(() => {
    if (
      !chatContainerRef.current ||
      !hasMoreMessages ||
      isLoading ||
      !selectedUser
    )
      return;

    const { scrollTop } = chatContainerRef.current;

    // If scrolled to top (with a small buffer), load more messages
    if (scrollTop < 50) {
      const scrollInfo = saveScrollPosition();
      fetchChatHistory(currentPage + 1);
    }
  }, [
    fetchChatHistory,
    currentPage,
    hasMoreMessages,
    isLoading,
    saveScrollPosition,
    selectedUser,
  ]);

  // Add scroll event listener
  useEffect(() => {
    const container = chatContainerRef.current;
    if (container) {
      container.addEventListener("scroll", handleScroll);
      return () => container.removeEventListener("scroll", handleScroll);
    }
  }, [handleScroll]);

  // Handle message input change
  const handleMessageChange = (e) => {
    setMessage(e.target.value);
  };

  // Handle emoji selection
  const handleEmojiClick = (emojiData) => {
    setMessage((prev) => prev + emojiData.emoji);
    // Focus back on the input after selecting an emoji
    if (messageInputRef.current) {
      messageInputRef.current.focus();
    }
  };

  // Send a message
  const handleSendMessage = (e) => {
    e.preventDefault();
    if (!message.trim()) return;

    // Send the message via the context
    const success = sendChatMessage(message.trim());
    setMessage("");
  };

  // Fetch chat history when a user is selected
  useEffect(() => {
    if (selectedUser) {
      fetchChatHistory(0);
    }
  }, [selectedUser, fetchChatHistory]);

  // If no user is selected, show the empty state
  if (!selectedUser) {
    return (
      <div
        className="flex-1 flex items-center justify-center h-full overflow-hidden"
        style={{
          background: "linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)",
        }}
      >
        <div className="text-center p-8 max-w-md bg-white rounded-xl shadow-lg border border-gray-100">
          <div className="bg-blue-100 rounded-full p-6 inline-block mb-6 shadow-inner">
            <svg
              className="h-16 w-16 text-blue-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
              ></path>
            </svg>
          </div>
          <h2 className="text-2xl font-bold mb-3 text-blue-800">
            Your Messages
          </h2>
          <p className="text-gray-600 mb-6">
            Select a conversation from the sidebar to start chatting with
            friends and colleagues.
          </p>
          {isMobile && (
            <div className="flex justify-center">
              <button
                className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-lg transition-colors duration-200 shadow-md font-medium flex items-center"
                onClick={toggleSidebar}
              >
                <svg
                  className="w-5 h-5 mr-2"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M4 6h16M4 12h16m-7 6h7"
                  ></path>
                </svg>
                Open Contacts
              </button>
            </div>
          )}
        </div>
      </div>
    );
  } else {
    return (
      <div
        className="flex-1 flex flex-col h-[calc(100vh-64px)]"
        style={{
          background: "linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)",
        }}
      >
        {selectedUser && (
          <div className="flex-shrink-0 p-4 border-b border-gray-200 bg-white shadow-sm">
            <div className="flex items-center">
              {isMobile && (
                <button
                  className="md:hidden mr-3 text-blue-600 p-2 rounded-full hover:bg-blue-50 transition-colors duration-200"
                  onClick={toggleSidebar}
                  aria-label="Show contacts"
                >
                  <svg
                    className="h-5 w-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M4 6h16M4 12h16m-7 6h7"
                    ></path>
                  </svg>
                </button>
              )}
              <div className="relative">
                {selectedUser.avatar ? (
                  <Image
                    src={
                      selectedUser.avatar_mime_type
                        ? `data:${selectedUser.avatar_mime_type};base64,${selectedUser.avatar}`
                        : "/imgs/defaultAvatar.jpg"
                    }
                    alt={selectedUser.nickname}
                    width={44}
                    height={44}
                    className="rounded-full border-2 border-white shadow-sm"
                  />
                ) : (
                  <div className="w-11 h-11 rounded-full bg-gradient-to-br from-blue-500 to-blue-700 flex items-center justify-center text-white shadow-sm border-2 border-white">
                    {selectedUser.nickname.charAt(0).toUpperCase()}
                  </div>
                )}
                <span
                  className={`absolute bottom-0 right-0 inline-block w-3.5 h-3.5 rounded-full border-2 border-white ${
                    onlineUsers.includes(String(selectedUser.user_id))
                      ? "bg-green-500 animate-pulse"
                      : "bg-gray-400"
                  }`}
                  title={
                    onlineUsers.includes(String(selectedUser.user_id))
                      ? "Online"
                      : "Offline"
                  }
                ></span>
              </div>
              <div className="ml-3">
                <h2 className="font-bold text-lg text-blue-800">
                  {selectedUser.nickname}
                </h2>
                <div className="text-xs text-gray-500 flex items-center">
                  {onlineUsers.includes(String(selectedUser.user_id)) ? (
                    <span className="text-green-600">Active now</span>
                  ) : (
                    <span>
                     Inactive
                    </span>
                  )}
                </div>
              </div>
            </div>
          </div>
        )}

        <div
          className="flex-1 overflow-y-auto p-4 scrollbar-thin scrollbar-thumb-gray-300 scrollbar-track-transparent"
          ref={chatContainerRef}
          onScroll={handleScroll}
        >
          {isLoading && currentPage > 0 && (
            <div className="flex justify-center my-2">
              <div className="animate-spin rounded-full h-5 w-5 border-t-2 border-b-2 border-blue-600"></div>
            </div>
          )}

          {chatHistory.length === 0 && !isLoading ? (
            <div className="flex items-center justify-center h-full">
              <div className="text-center p-6 bg-white rounded-xl shadow-md border border-gray-100 max-w-md">
                <div className="bg-blue-100 rounded-full p-4 inline-block mb-4 shadow-inner">
                  <svg
                    className="h-10 w-10 text-blue-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"
                    ></path>
                  </svg>
                </div>
                <h3 className="text-xl font-bold mb-2 text-blue-800">
                  Start a conversation
                </h3>
                <p className="text-gray-600 mb-4">
                  Send your first message to {selectedUser.nickname}!
                </p>
                <button
                  onClick={() => messageInputRef.current?.focus()}
                  className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-lg transition-colors duration-200 shadow-md font-medium"
                >
                  Say Hello
                </button>
              </div>
            </div>
          ) : (
            Object.entries(groupedMessages).map(([date, messages]) => (
              <div key={date} className="mb-6">
                <div className="flex justify-center mb-4">
                  <div className="bg-gray-200 text-gray-600 px-3 py-1 rounded-full text-xs font-medium">
                    {new Date(date).toLocaleDateString(undefined, {
                      weekday: "long",
                      month: "short",
                      day: "numeric",
                    })}
                  </div>
                </div>

                {messages.map((msg) => {
                  if (msg.is_system) {
                    return (
                      <div key={msg.id} className="flex justify-center my-2">
                        <div className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-xs max-w-[80%] shadow-sm">
                          {msg.message}
                        </div>
                      </div>
                    );
                  }

                  const isCurrentUser =
                    msg.sender_id !== String(selectedUser.user_id);

                  return (
                    <div
                      key={msg.id}
                      className={`flex mb-3 ${
                        isCurrentUser ? "justify-end" : "justify-start"
                      }`}
                    >
                      {!isCurrentUser && (
                        <div className="flex-shrink-0 mr-2 self-end mb-1">
                          {selectedUser.avatar ? (
                            <Image
                              src={
                                selectedUser.avatar_mime_type
                                  ? `data:${selectedUser.avatar_mime_type};base64,${selectedUser.avatar}`
                                  : "/imgs/defaultAvatar.jpg"
                              }
                              alt={selectedUser.nickname}
                              width={28}
                              height={28}
                              className="rounded-full border border-gray-200"
                            />
                          ) : (
                            <div className="w-7 h-7 rounded-full bg-gradient-to-br from-blue-500 to-blue-700 flex items-center justify-center text-white text-xs border border-gray-200">
                              {selectedUser.nickname.charAt(0).toUpperCase()}
                            </div>
                          )}
                        </div>
                      )}
                      <div
                        className={`max-w-[75%] ${
                          isCurrentUser
                            ? "bg-blue-600 text-white rounded-2xl rounded-tr-sm shadow-sm"
                            : "bg-white text-gray-800 rounded-2xl rounded-tl-sm shadow-sm border border-gray-200"
                        } px-4 py-2.5 break-words`}
                      >
                        <div className="whitespace-pre-wrap">{msg.message}</div>
                        <div
                          className={`text-xs mt-1 ${
                            isCurrentUser ? "text-blue-200" : "text-gray-500"
                          } flex items-center`}
                        >
                          <span>{formatTimestamp(msg.created_at)}</span>
                          {msg.pending && (
                            <span className="ml-1">
                              {msg.notSent ? "⚠️ Not sent" : "⏳ Sending..."}
                            </span>
                          )}
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            ))
          )}
          <div ref={messagesEndRef} />
        </div>

        {selectedUser && (
          <div className="flex-shrink-0 border-t border-gray-200 p-4 bg-white">
            <form
              onSubmit={handleSendMessage}
              className="flex items-center gap-2"
            >
              <button
                type="button"
                className="flex-shrink-0 text-gray-500 hover:text-blue-600 p-2 rounded-full hover:bg-blue-50"
                onClick={() => setShowEmojiPicker(!showEmojiPicker)}
                title="Add emoji"
              >
                <svg
                  className="h-5 w-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  ></path>
                </svg>
              </button>

              <div className="relative flex-1">
                <input
                  ref={messageInputRef}
                  type="text"
                  value={message}
                  onChange={handleMessageChange}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" && !e.shiftKey) {
                      e.preventDefault();
                      if (message.trim()) {
                        handleSendMessage(e);
                      }
                    }
                  }}
                  placeholder="Type a message..."
                  className="w-full h-10 border border-gray-300 rounded-full px-4 pr-12 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                />
                {message.trim() && (
                  <button
                    type="submit"
                    className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 bg-blue-600 text-white rounded-full hover:bg-blue-700"
                  >
                    <svg
                      className="h-4 w-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M5 12h14M12 5l7 7-7 7"
                      ></path>
                    </svg>
                  </button>
                )}
              </div>
            </form>

            {showEmojiPicker && (
              <div className="absolute bottom-[80px] right-4 bg-white rounded-lg shadow-lg border border-gray-200">
                <div className="relative">
                  <button
                    onClick={() => setShowEmojiPicker(false)}
                    className="absolute top-2 right-2 text-gray-500 hover:text-gray-700 z-10 bg-white rounded-full p-1"
                    title="Close emoji picker"
                  >
                    <svg
                      className="h-4 w-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M6 18L18 6M6 6l12 12"
                      ></path>
                    </svg>
                  </button>
                  <EmojiPicker
                    onEmojiClick={handleEmojiClick}
                    searchDisabled={false}
                    width={300}
                    height={400}
                    previewConfig={{ showPreview: false }}
                    skinTonesDisabled
                  />
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    );
  }
}
