"use client";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/context/AuthContext";
import Header from "./components/Headers/Header";
import WebsocketProvider from "@/context/Websocket";
import { NotificationProvider } from "@/context/NotificationContext";
import { ChatPageProvider } from "@/context/ChatPageContext";
import UserNotifier from "./components/Alerts/UserNotifier";
import ChatNotifier from "./components/chat/ChatNotifier";
import { usePathname } from 'next/navigation';
import {
  AlertProvider,
  ConfirmAction,
  Notification,
} from "./components/Alerts/PopUp";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function RootLayout({ children }) {
  // Use pathname to determine if we're on the chat page
  const pathname = usePathname();
  const isChatPage = pathname === '/chat';

  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <div className="main-container">
          <AuthProvider>
            <WebsocketProvider>
              <AlertProvider>
                <NotificationProvider>
                  <ChatPageProvider>
                    {/* Only render the header if we're not on the chat page */}
                    {!isChatPage && <Header />}
                    <UserNotifier />
                    <ChatNotifier />
                    <ConfirmAction />
                    <Notification />
                    {children}
                  </ChatPageProvider>
                </NotificationProvider>
              </AlertProvider>
            </WebsocketProvider>
          </AuthProvider>
        </div>
      </body>
    </html>
  );
}
