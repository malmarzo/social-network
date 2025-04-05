'use client';

import { useState, useEffect, useRef } from 'react';
import { useChat } from '@/context/ChatContext';
import { useRouter, useSearchParams } from 'next/navigation';
import Image from 'next/image';
import { formatDistanceToNow } from 'date-fns';
import { useWebSocket } from '@/context/Websocket';
import { useNotification } from '@/context/NotificationContext';
import { useAuth } from '@/context/AuthContext';

export default function ChatSidebar() {
  const { isConnected } = useWebSocket();
  const { showInfo } = useNotification();
  const { userID } = useAuth();
  const {
    chatUsers,
    allUsers,
    onlineUsers,
    selectedUser,
    isLoading,
    showAllUsers,
    unreadMessages,
    toggleUserView,
    selectUser,
    formatTimestamp,
    truncateMessage
  } = useChat();
  const [searchTerm, setSearchTerm] = useState('');
  const searchInputRef = useRef(null);

  // Get eligible chat users with their unread status and ensure last_message_time is properly set
  const usersWithUnreadStatus = allUsers
    .map(user => ({
      ...user,
      unreadCount: unreadMessages[user.user_id] || 0,
      // Ensure last_message_time is set for sorting
      last_message_time: user.last_message_time || user.last_activity || null
    }));
  
  // Sort users by unread messages and recency
  const filteredUsers = usersWithUnreadStatus
    .filter(user => {
      if (!searchTerm) return true;
      return user.nickname.toLowerCase().includes(searchTerm.toLowerCase());
    })
    .sort((a, b) => {
      // First sort by unread messages (highest priority)
      if (a.unreadCount !== b.unreadCount) {
        return b.unreadCount - a.unreadCount; // Users with unread messages first
      }
      
      // Then sort by most recent message
      const aTime = a.last_message_time ? new Date(a.last_message_time).getTime() : 0;
      const bTime = b.last_message_time ? new Date(b.last_message_time).getTime() : 0;
      return bTime - aTime; // Most recent conversations first
      
      // Check if users have message history
      const aHasMessage = Boolean(a.last_message);
      const bHasMessage = Boolean(b.last_message);
      
      // If one has a message and the other doesn't, prioritize the one with a message
      if (aHasMessage !== bHasMessage) {
        return aHasMessage ? -1 : 1; // User with message comes first
      }
      
      // If both have messages, sort by most recent message time
      if (aHasMessage && bHasMessage) {
        // Get timestamps, defaulting to 0 if not present
        const aTime = a.last_message_time ? new Date(a.last_message_time).getTime() : 0;
        const bTime = b.last_message_time ? new Date(b.last_message_time).getTime() : 0;
        return bTime - aTime; // Most recent conversations first
      }
      
      // If neither has messages, keep the original order
      return 0;
    });

  // Focus search input when pressing '/' key
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === '/' && document.activeElement.tagName !== 'INPUT') {
        e.preventDefault();
        searchInputRef.current?.focus();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, []);

  // Clear search when Escape key is pressed
  const handleSearchKeyDown = (e) => {
    if (e.key === 'Escape') {
      setSearchTerm('');
      searchInputRef.current?.blur();
    }
  };

  return (
    <div className="w-full h-full flex flex-col overflow-hidden bg-gradient-to-b from-gray-50 to-white border-r border-gray-200 shadow-lg">
      {/* Header */}
      <div className="p-4 border-b border-gray-200 sticky top-0 z-10 bg-white shadow-sm">
        <div className="flex justify-between items-center mb-3">
          <h2 className="text-xl font-bold text-blue-800 ml-20">Messages</h2>
          {!isConnected ? (
            <div className="text-xs text-red-500 flex items-center bg-red-50 px-2 py-1 rounded-full">
              <span className="inline-block w-2 h-2 rounded-full bg-red-500 mr-1 animate-pulse"></span>
              Offline
            </div>
          ) : (
            <div className="text-xs text-green-600 flex items-center bg-green-50 px-2 py-1 rounded-full">
              <span className="inline-block w-2 h-2 rounded-full bg-green-500 mr-1 animate-pulse"></span>
              Online
            </div>
          )}
        </div>
        
        {/* Search bar */}
        <div className="relative ml-0">
          <input
            ref={searchInputRef}
            type="text"
            placeholder="Search conversations... (Press '/' to focus)"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onKeyDown={handleSearchKeyDown}
            className="w-full pl-10 pr-4 py-2 rounded-full border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-gray-50 text-sm"
          />
          <div className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>
          {searchTerm && (
            <button
              onClick={() => setSearchTerm('')}
              className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          )}
        </div>
      </div>

      {/* User list */}
      <div className="overflow-y-auto flex-1 scrollbar-thin scrollbar-thumb-gray-300 scrollbar-track-transparent">
        {isLoading ? (
          <div className="flex justify-center items-center h-32">
            <div className="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-blue-500 shadow-md"></div>
          </div>
        ) : filteredUsers.length === 0 ? (
          <div className="p-8 text-center text-gray-500">
            {searchTerm ? (
              <>
                <div className="mb-3 text-4xl">🔍</div>
                <p className="font-medium mb-1">No matches found</p>
                <p className="text-sm text-gray-400">Try a different search term</p>
              </>
            ) : (
              <>
                <div className="mb-3 text-4xl">💬</div>
                <p className="font-medium mb-1">No conversations yet</p>
                <p className="text-sm text-gray-400">Start chatting with someone!</p>
              </>
            )}
          </div>
        ) : (
          <div className="divide-y divide-gray-100">
            {filteredUsers.map((user) => (
              <div
                key={user.user_id}
                className={`p-4 cursor-pointer transition-all duration-200 hover:bg-blue-50 ${selectedUser && selectedUser.user_id === user.user_id ? 'bg-blue-50 border-l-4 border-l-blue-500' : 'border-l-4 border-l-transparent'}`}
                onClick={() => {
                  selectUser(user);
                  
                  // Show notification when switching to a user with unread messages
                  const unreadCount = unreadMessages[user.user_id] || 0;
                  if (unreadCount > 0) {
                    showInfo(`${unreadCount} unread message${unreadCount > 1 ? 's' : ''} from ${user.nickname}`, {
                      duration: 3000,
                      position: 'top-right'
                    });
                  }
                }}
              >
                <div className="flex items-center">
                  <div className="relative flex-shrink-0">
                    {user.avatar ? (
                      <Image
                        src={user.avatar_mime_type ? 
                          `data:${user.avatar_mime_type};base64,${user.avatar}` : 
                          "/imgs/defaultAvatar.jpg"}
                        alt={user.nickname}
                        width={48}
                        height={48}
                        className="rounded-full border-2 border-white shadow-sm"
                      />
                    ) : (
                      <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-500 to-blue-700 flex items-center justify-center text-white shadow-sm border-2 border-white">
                        {user.nickname.charAt(0).toUpperCase()}
                      </div>
                    )}
                    
                    {/* Online status indicator */}
                    <div 
                      className={`absolute bottom-0 right-0 w-3.5 h-3.5 rounded-full border-2 border-white ${user.isOnline || onlineUsers.includes(String(user.user_id)) ? 'bg-green-500 animate-pulse' : 'bg-gray-400'}`} 
                      title={user.isOnline || onlineUsers.includes(String(user.user_id)) ? 'Online' : user.lastSeen ? `Last seen ${formatTimestamp(new Date(user.lastSeen * 1000).toISOString())}` : 'Offline'}
                    ></div>
                    
                    {/* Unread message badge */}
                    {unreadMessages[user.user_id] > 0 && (
                      <div 
                        className="absolute -top-1 -right-1 bg-blue-600 text-white text-xs rounded-full w-6 h-6 flex items-center justify-center shadow-md border border-white"
                        title={`${unreadMessages[user.user_id]} unread message${unreadMessages[user.user_id] > 1 ? 's' : ''}`}
                      >
                        {unreadMessages[user.user_id] > 9 ? '9+' : unreadMessages[user.user_id]}
                      </div>
                    )}
                  </div>
                  <div className="ml-3 flex-1 min-w-0"> {/* min-width-0 helps with text truncation */}
                    <div className="flex justify-between items-center">
                      <h3 className={`font-medium truncate ${unreadMessages[user.user_id] > 0 ? 'text-black font-semibold' : 'text-gray-800'}`}>
                        {user.nickname}
                      </h3>
                      <span className="text-xs text-gray-500 whitespace-nowrap ml-2 flex-shrink-0">
                        {formatTimestamp(user.last_message_time)}
                      </span>
                    </div>
                    <p className={`text-sm mt-1 truncate ${unreadMessages[user.user_id] > 0 ? 'text-black font-medium' : user.last_message ? 'text-gray-600' : 'text-gray-400 italic'}`}>
                      {user.last_message ? truncateMessage(user.last_message) : 'Start a conversation...'}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Status footer */}
      <div className="p-3 border-t border-gray-200 bg-gray-50 text-xs text-gray-500 flex justify-between items-center">
        <div className="flex items-center">
          <span className={`inline-block w-2 h-2 rounded-full mr-1.5 ${isConnected ? 'bg-green-500' : 'bg-red-500'}`}></span>
          <span>{isConnected ? 'Connected' : 'Disconnected'}</span>
        </div>
        <div>
          {filteredUsers.length} {filteredUsers.length === 1 ? 'conversation' : 'conversations'}
        </div>
      </div>
    </div>
  );
}
