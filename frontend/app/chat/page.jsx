'use client';

import { useEffect, useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { ChatProvider } from '@/context/ChatContext';
import ChatSidebar from '../components/chat/ChatSidebar';
import ChatWindow from '../components/chat/ChatWindow';

export default function ChatPage() {
  const { isLoggedIn, loading } = useAuth();
  const searchParams = useSearchParams();
  const router = useRouter();
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [isMobile, setIsMobile] = useState(false);
  
  // Toggle sidebar visibility
  const toggleSidebar = () => setSidebarOpen(!sidebarOpen);

  // Add an effect to disable body scrolling when the chat page is mounted
  useEffect(() => {
    // Save the original overflow style
    const originalStyle = document.body.style.overflow;
    
    // Disable scrolling on the body
    document.body.style.overflow = 'hidden';
    
    // Check if we're on mobile
    const checkMobile = () => {
      const isMobileView = window.innerWidth < 768;
      setIsMobile(isMobileView);
      // Auto-close sidebar on mobile
      if (isMobileView && sidebarOpen) {
        setSidebarOpen(false);
      } else if (!isMobileView && !sidebarOpen) {
        setSidebarOpen(true);
      }
    };
    
    // Check on initial load
    checkMobile();
    
    // Add resize listener
    window.addEventListener('resize', checkMobile);
    
    // Restore original style when component unmounts
    return () => {
      document.body.style.overflow = originalStyle;
      window.removeEventListener('resize', checkMobile);
    };
  }, [sidebarOpen]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!isLoggedIn) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Please log in to access chat</h1>
          <a href="/login" className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded">
            Go to Login
          </a>
        </div>
      </div>
    );
  }



  return (
    <ChatProvider initialUrlUserId={searchParams.get('userId')}>
      <div className="relative flex h-screen overflow-hidden bg-gray-100">
        {/* Mobile menu button */}
        {isMobile && (
          <button 
            onClick={toggleSidebar}
            className="md:hidden fixed top-4 left-4 z-50 bg-blue-500 text-white p-2 rounded-full shadow-lg"
            aria-label={sidebarOpen ? 'Close menu' : 'Open menu'}
          >
            {sidebarOpen ? (
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            ) : (
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            )}
          </button>
        )}
        
        {/* Sidebar with responsive behavior */}
        <div 
          className={`${sidebarOpen ? 'translate-x-0' : '-translate-x-full'} 
            md:translate-x-0 transition-transform duration-300 ease-in-out 
            fixed md:relative z-40 md:z-auto w-80 md:w-1/3 lg:w-1/4 h-full`}
        >
          <ChatSidebar />
        </div>
        
        {/* Overlay for mobile */}
        {sidebarOpen && isMobile && (
          <div 
            className="fixed inset-0 bg-black bg-opacity-50 z-30 md:hidden"
            onClick={toggleSidebar}
          />
        )}
        
        {/* Main chat window */}
        <div className="flex-1 relative">
          <ChatWindow toggleSidebar={toggleSidebar} isMobile={isMobile} />
        </div>
      </div>
    </ChatProvider>
  );
}
