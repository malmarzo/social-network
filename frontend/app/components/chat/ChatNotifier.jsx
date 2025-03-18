'use client';

import { useEffect, useState, useRef } from 'react';
import { useWebSocket } from '@/context/Websocket';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';

// Component to handle chat notifications
const ChatNotifier = () => {
  const { addMessageHandler } = useWebSocket();
  const router = useRouter();
  const { isLoggedIn, currentUser } = useAuth();
  const [notificationSound] = useState(() => typeof Audio !== 'undefined' ? new Audio('/notification.mp3') : null);
  const currentPathRef = useRef('');
  
  // Keep track of current path
  useEffect(() => {
    if (typeof window !== 'undefined') {
      currentPathRef.current = window.location.pathname;
      
      const handleRouteChange = () => {
        currentPathRef.current = window.location.pathname;
      };
      
      window.addEventListener('popstate', handleRouteChange);
      return () => window.removeEventListener('popstate', handleRouteChange);
    }
  }, []);

  useEffect(() => {
    if (!isLoggedIn || !currentUser) return;

    // Handle chat message notifications
    const handleChatMessage = (msg) => {
      // Only show notification if the message is not from the current user
      if (msg.type === 'chat' && msg.userDetails.id !== String(currentUser.id)) {
        // Play notification sound if we're not already on the chat page
        const isOnChatPage = currentPathRef.current.includes('/chat');
        
        // If we're not on the chat page or the window is not focused, show a notification
        if (!isOnChatPage || document.hidden) {
          // Play sound notification
          if (notificationSound) {
            try {
              // Reset the audio to the beginning if it's already playing
              notificationSound.pause();
              notificationSound.currentTime = 0;
              notificationSound.play().catch(err => console.log('Error playing notification sound:', err));
            } catch (error) {
              console.error('Failed to play notification sound:', error);
            }
          }
          
          // Create a browser notification
          if (Notification.permission === 'granted') {
            const notification = new Notification('New Message from ' + msg.userDetails.nickname, {
              body: msg.content.length > 50 ? msg.content.substring(0, 47) + '...' : msg.content,
              icon: '/favicon.ico',
              badge: '/favicon.ico',
              tag: `chat-${msg.userDetails.id}`, // Group notifications from the same sender
              renotify: true // Notify each time even with the same tag
            });

            // Navigate to chat when notification is clicked
            notification.onclick = () => {
              router.push(`/chat?user=${msg.userDetails.id}`);
              window.focus();
            };
          } else if (Notification.permission !== 'denied') {
            Notification.requestPermission().then(permission => {
              if (permission === 'granted') {
                const notification = new Notification('New Message from ' + msg.userDetails.nickname, {
                  body: msg.content.length > 50 ? msg.content.substring(0, 47) + '...' : msg.content,
                  icon: '/favicon.ico',
                  badge: '/favicon.ico',
                  tag: `chat-${msg.userDetails.id}`,
                  renotify: true
                });

                notification.onclick = () => {
                  router.push(`/chat?user=${msg.userDetails.id}`);
                  window.focus();
                };
              }
            });
          }
        }
      }
    };

    // Request notification permission when component mounts
    if (typeof window !== 'undefined' && 'Notification' in window) {
      if (Notification.permission !== 'granted' && Notification.permission !== 'denied') {
        Notification.requestPermission();
      }
    }

    // Add message handler for chat messages
    addMessageHandler('chat', handleChatMessage);

    // Clean up
    return () => {
      // This is a placeholder for cleanup if needed
    };
  }, [addMessageHandler, router, isLoggedIn, currentUser, notificationSound]);

  return null; // This component doesn't render anything
};

export default ChatNotifier;
