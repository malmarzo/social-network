'use client';

import { createContext, useContext, useState, useEffect } from 'react';

const ChatPageContext = createContext({
  isOnChatPage: false,
  setIsOnChatPage: () => {}
});

export function ChatPageProvider({ children }) {
  const [isOnChatPage, setIsOnChatPage] = useState(false);
  
  // Update chat page status when route changes
  useEffect(() => {
    // Check if we're on the chat page
    const checkIfChatPage = () => {
      const isChatRoute = window.location.pathname.includes('/chat');
      // Only update if the value has changed to avoid unnecessary renders
      if (isChatRoute !== isOnChatPage) {
        setIsOnChatPage(isChatRoute);
      }
    };
    
    // Check on initial load
    checkIfChatPage();
    
    // Add event listener for browser back/forward navigation
    const handlePopState = () => {
      // Use setTimeout to ensure this runs after React's lifecycle methods
      setTimeout(checkIfChatPage, 0);
    };
    window.addEventListener('popstate', handlePopState);
    
    // For Next.js route changes
    if (typeof window !== 'undefined') {
      // Use Next.js router events instead of overriding history methods
      const handleRouteChange = () => {
        // Use setTimeout to ensure this runs after React's lifecycle methods
        setTimeout(checkIfChatPage, 0);
      };
      
      // Create a MutationObserver to watch for URL changes
      // This is a safer approach than overriding history methods
      const observer = new MutationObserver((mutations) => {
        mutations.forEach(() => {
          setTimeout(checkIfChatPage, 0);
        });
      });
      
      // Start observing the document with the configured parameters
      observer.observe(document, { subtree: true, childList: true });
      
      return () => {
        // Clean up all event listeners and observers
        window.removeEventListener('popstate', handlePopState);
        observer.disconnect();
      };
    }
    
    return () => {
      // Clean up event listener
      window.removeEventListener('popstate', handlePopState);
    };
  }, [isOnChatPage]);
  
  return (
    <ChatPageContext.Provider value={{ isOnChatPage, setIsOnChatPage }}>
      {children}
    </ChatPageContext.Provider>
  );
}

export function useChatPage() {
  return useContext(ChatPageContext);
}
