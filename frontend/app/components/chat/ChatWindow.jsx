'use client';

import { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { useChat } from '@/context/ChatContext';
import { useNotification } from '@/context/NotificationContext';
import { useWebSocket } from '@/context/Websocket';
import { useAuth } from '@/context/AuthContext';
import Image from 'next/image';
import { formatDistanceToNow } from 'date-fns';

export default function ChatWindow({ toggleSidebar, isMobile }) {
  const [message, setMessage] = useState('');
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const messagesEndRef = useRef(null);
  const chatContainerRef = useRef(null);
  const messageInputRef = useRef(null);
  const { isConnected } = useWebSocket();
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
    onlineUsers
  } = useChat();
  
  // Group messages by date for better UI organization
  const groupedMessages = useMemo(() => {
    if (!chatHistory || chatHistory.length === 0) return {};
    
    const groups = {};
    chatHistory.forEach(msg => {
      const date = new Date(msg.created_at).toLocaleDateString();
      if (!groups[date]) groups[date] = [];
      groups[date].push(msg);
    });
    
    return groups;
  }, [chatHistory]);

  // Auto-scroll to bottom when chat history changes
  useEffect(() => {
    if (messagesEndRef.current && chatHistory.length > 0 && selectedUser) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
      
      // Check if the last message is from the other user and is new (within the last 2 seconds)
      const lastMessage = chatHistory[chatHistory.length - 1];
      if (lastMessage && 
          lastMessage.sender_id !== userID && 
          (new Date() - new Date(lastMessage.created_at)) < 2000) {
        
        // Play notification sound if available
        try {
          const audio = new Audio('/notification.mp3');
          audio.play().catch(e => console.log('Audio play failed:', e));
        } catch (e) {
          console.log('Audio not supported:', e);
        }
        
        // Try to show browser notification if permission is granted
        if ("Notification" in window) {
          if (Notification.permission === "granted") {
            new Notification(`${lastMessage.sender_name} sent you a message`, {
              body: lastMessage.message,
              icon: '/favicon.ico',
              tag: 'chat-message', // Replace previous notification
              renotify: true // Notify again even with same tag
            });
          } else if (Notification.permission !== "denied") {
            // Request permission if not already denied
            Notification.requestPermission();
          }
        }
      }
    }
  }, [chatHistory, selectedUser, userID]);

  // Save scroll position when loading more messages
  const saveScrollPosition = useCallback(() => {
    if (chatContainerRef.current) {
      const container = chatContainerRef.current;
      const scrollHeight = container.scrollHeight;
      const scrollTop = container.scrollTop;
      const clientHeight = container.clientHeight;
      
      // If user has scrolled up more than 200px, consider it a deliberate scroll up
      return { scrollHeight, scrollPosition: scrollHeight - scrollTop - clientHeight < 200 };
    }
    return { scrollHeight: 0, scrollPosition: true };
  }, []);

  // Load more messages when scrolling to top
  const handleScroll = useCallback(() => {
    if (!chatContainerRef.current || !hasMoreMessages || isLoading || !selectedUser) return;
    
    const { scrollTop } = chatContainerRef.current;
    
    // If scrolled to top (with a small buffer), load more messages
    if (scrollTop < 50) {
      const scrollInfo = saveScrollPosition();
      fetchChatHistory(currentPage + 1);
    }
  }, [fetchChatHistory, currentPage, hasMoreMessages, isLoading, saveScrollPosition, selectedUser]);

  // Add scroll event listener
  useEffect(() => {
    const container = chatContainerRef.current;
    if (container) {
      container.addEventListener('scroll', handleScroll);
      return () => container.removeEventListener('scroll', handleScroll);
    }
  }, [handleScroll]);

  // Handle message input change
  const handleMessageChange = (e) => {
    setMessage(e.target.value);
  };

  // Send a message
  const handleSendMessage = (e) => {
    e.preventDefault();
    if (!message.trim()) return;

    // Send the message via the context
    const success = sendChatMessage(message.trim());
    
    // Show notification based on success
    if (success) {
      showSuccess('Message sent', { duration: 2000, position: 'top-right' });
    } else if (!isConnected) {
      showError('You are offline. Message will be sent when connection is restored.', { 
        duration: 2000, 
        position: 'top-right' 
      });
    }
    
    // Request notification permission if not already granted or denied
    if ("Notification" in window && Notification.permission !== "granted" && Notification.permission !== "denied") {
      Notification.requestPermission();
    }
    
    // Clear the input field
    setMessage('');
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
      <div className="flex-1 flex items-center justify-center bg-gradient-to-b from-blue-50 to-white h-full overflow-hidden">
        <div className="text-center p-8 max-w-md bg-white rounded-xl shadow-lg border border-gray-100">
          <div className="bg-blue-100 rounded-full p-6 inline-block mb-6 shadow-inner">
            <svg className="h-16 w-16 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
            </svg>
          </div>
          <h2 className="text-2xl font-bold mb-3 text-blue-800">Your Messages</h2>
          <p className="text-gray-600 mb-6">
            Select a conversation from the sidebar to start chatting with friends and colleagues.
          </p>
          {isMobile && (
            <div className="flex justify-center">
              <button 
                className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-lg transition-colors duration-200 shadow-md font-medium flex items-center"
                onClick={toggleSidebar}
              >
                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16m-7 6h7"></path>
                </svg>
                Open Contacts
              </button>
            </div>
          )}
          {!isConnected && (
            <div className="mt-6 p-4 bg-red-50 text-red-700 rounded-lg border border-red-100 shadow-sm">
              <p className="font-medium flex items-center">
                <svg className="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
                </svg>
                You are currently offline
              </p>
              <p className="text-sm ml-7">Messages will be sent when your connection is restored.</p>
            </div>
          )}
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col bg-gradient-to-b from-gray-50 to-white h-full overflow-hidden">
      {/* Chat header */}
      {selectedUser && (
        <div className="p-4 border-b border-gray-200 flex items-center justify-between bg-white shadow-sm relative">
          <div className="flex items-center">
            {isMobile && (
              <button 
                className="md:hidden mr-3 text-blue-600 p-2 rounded-full hover:bg-blue-50 transition-colors duration-200"
                onClick={toggleSidebar}
                aria-label="Show contacts"
              >
                <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16m-7 6h7"></path>
                </svg>
              </button>
            )}
            <div className="relative">
              {selectedUser.avatar ? (
                <Image
                  src={selectedUser.avatar.startsWith('/uploads/') ? selectedUser.avatar : selectedUser.avatar.startsWith('uploads/') ? `/${selectedUser.avatar}` : `/uploads/${selectedUser.avatar}`}
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
                  onlineUsers.includes(String(selectedUser.user_id)) ? 'bg-green-500 animate-pulse' : 'bg-gray-400'
                }`}
                title={onlineUsers.includes(String(selectedUser.user_id)) ? 'Online' : 'Offline'}
              ></span>
            </div>
            <div className="ml-3">
              <h2 className="font-bold text-lg text-blue-800">{selectedUser.nickname}</h2>
              <div className="text-xs text-gray-500 flex items-center">
                {onlineUsers.includes(String(selectedUser.user_id)) ? (
                  <span className="text-green-600">Active now</span>
                ) : (
                  <span>Last seen {formatDistanceToNow(new Date(selectedUser.last_activity || Date.now()), { addSuffix: true })}</span>
                )}
              </div>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <button className="text-gray-500 hover:text-blue-600 p-2 rounded-full hover:bg-blue-50 transition-colors duration-200" title="Voice call">
              <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
              </svg>
            </button>
            <button className="text-gray-500 hover:text-blue-600 p-2 rounded-full hover:bg-blue-50 transition-colors duration-200" title="More options">
              <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path>
              </svg>
            </button>
          </div>
        </div>
      )}

      {/* Chat messages */}
      {selectedUser && (
        <div
          className="flex-1 overflow-y-auto p-4 max-h-[calc(100vh-140px)] scrollbar-thin scrollbar-thumb-gray-300 scrollbar-track-transparent"
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
                  <svg className="h-10 w-10 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"></path>
                  </svg>
                </div>
                <h3 className="text-xl font-bold mb-2 text-blue-800">Start a conversation</h3>
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
            // Display messages grouped by date
            Object.entries(groupedMessages).map(([date, messages]) => (
              <div key={date} className="mb-6">
                {/* Date separator */}
                <div className="flex justify-center mb-4">
                  <div className="bg-gray-200 text-gray-600 px-3 py-1 rounded-full text-xs font-medium">
                    {new Date(date).toLocaleDateString(undefined, { weekday: 'long', month: 'short', day: 'numeric' })}
                  </div>
                </div>
                
                {/* Messages for this date */}
                {messages.map((msg) => {
                  // Handle system messages
                  if (msg.is_system) {
                    return (
                      <div key={msg.id} className="flex justify-center my-2">
                        <div className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-xs max-w-[80%] shadow-sm">
                          {msg.message}
                        </div>
                      </div>
                    );
                  }
                  
                  const isCurrentUser = msg.sender_id !== String(selectedUser.user_id);
                  
                  return (
                    <div
                      key={msg.id}
                      className={`flex mb-3 ${isCurrentUser ? 'justify-end' : 'justify-start'}`}
                    >
                      {!isCurrentUser && (
                        <div className="flex-shrink-0 mr-2 self-end mb-1">
                          {selectedUser.avatar ? (
                            <Image
                              src={selectedUser.avatar.startsWith('/uploads/') ? selectedUser.avatar : selectedUser.avatar.startsWith('uploads/') ? `/${selectedUser.avatar}` : `/uploads/${selectedUser.avatar}`}
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
                            ? 'bg-blue-600 text-white rounded-2xl rounded-tr-sm shadow-sm'
                            : 'bg-white text-gray-800 rounded-2xl rounded-tl-sm shadow-sm border border-gray-200'
                        } px-4 py-2.5 break-words`}
                      >
                        <div className="whitespace-pre-wrap">{msg.message}</div>
                        <div
                          className={`text-xs mt-1 ${
                            isCurrentUser ? 'text-blue-200' : 'text-gray-500'
                          } flex items-center`}
                        >
                          <span>{formatTimestamp(msg.created_at)}</span>
                          {msg.pending && (
                            <span className="ml-1">
                              {msg.notSent ? '⚠️ Not sent' : '⏳ Sending...'}
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
      )}

      {/* Message input */}
      {selectedUser && (
        <div className="border-t border-gray-200 p-4 bg-white sticky bottom-0 z-10 shadow-sm">
          <form onSubmit={handleSendMessage} className="flex items-center">
            <button
              type="button"
              className="text-gray-500 hover:text-blue-600 p-2 rounded-full hover:bg-blue-50 transition-colors duration-200 mr-1"
              onClick={() => setShowEmojiPicker(!showEmojiPicker)}
              title="Add emoji"
            >
              <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </button>
            <button
              type="button"
              className="text-gray-500 hover:text-blue-600 p-2 rounded-full hover:bg-blue-50 transition-colors duration-200 mr-1"
              title="Attach file"
            >
              <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"></path>
              </svg>
            </button>
            <div className="relative flex-1">
              <input
                ref={messageInputRef}
                type="text"
                value={message}
                onChange={handleMessageChange}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    if (message.trim()) {
                      handleSendMessage(e);
                    }
                  }
                }}
                placeholder="Type a message..."
                className="w-full border border-gray-300 rounded-full px-4 py-2.5 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent pr-10 bg-white"
                disabled={!isConnected}
              />
              {message.trim() && isConnected && (
                <button
                  type="submit"
                  className="absolute right-1 top-1/2 transform -translate-y-1/2 p-1.5 bg-blue-600 text-white rounded-full hover:bg-blue-700 transition-colors duration-200"
                >
                  <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 12h14M12 5l7 7-7 7"></path>
                  </svg>
                </button>
              )}
            </div>
          </form>
          {!isConnected && (
            <div className="mt-3 text-sm text-red-500 flex items-center bg-red-50 p-2 rounded-lg border border-red-100">
              <svg className="h-4 w-4 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <span>You are currently offline. Messages will be sent when your connection is restored.</span>
            </div>
          )}
          {showEmojiPicker && (
            <div className="absolute bottom-20 right-4 bg-white rounded-lg shadow-lg border border-gray-200 p-2 z-20">
              {/* Emoji picker will be implemented here */}
              <div className="p-4 text-center text-gray-600">
                <svg className="h-8 w-8 mx-auto mb-2 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <p className="font-medium">Emoji picker coming soon!</p>
                <p className="text-xs mt-1 text-gray-500">This feature is under development</p>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Empty state when no user is selected */}
      {!selectedUser && (
        <div className="flex-1 flex items-center justify-center bg-white">
          <div className="text-center p-6 max-w-md">
            <div className="text-5xl mb-4">💬</div>
            <h2 className="text-2xl font-semibold mb-2">Select a conversation</h2>
            <p className="text-gray-600">
              Choose a user from the sidebar to start chatting or continue a conversation.
            </p>
            {!isConnected && (
              <div className="mt-4 p-3 bg-red-50 text-red-700 rounded-md">
                <p className="font-medium">You are currently offline</p>
                <p className="text-sm">Messages will be sent when your connection is restored.</p>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
