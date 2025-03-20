"use client";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/context/AuthContext";
import Header from "./components/Headers/Header";
import WebsocketProvider from "@/context/Websocket";
import UserNotifier from "./components/Alerts/UserNotifier";
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
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <div className="main-container">
          <AuthProvider>
            <WebsocketProvider>
              <AlertProvider>
                <UserNotifier />
                <Header />
                <ConfirmAction />
                <Notification />
                {children}
              </AlertProvider>
            </WebsocketProvider>
          </AuthProvider>
        </div>
      </body>
    </html>
  );
}
