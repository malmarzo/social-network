
import React from "react";
import style from "@/styles/Header.module.css";
import AuthButton from "../Buttons/AuthButtons";
import LogoutButton from "@/app/logout/page";
import { useAuth } from "@/context/AuthContext";
import Link from "next/link";
import { HomeIcon } from "@heroicons/react/24/outline";
import NotificationButton from "./NotificationButton";
import { useWebSocket } from "@/context/Websocket";
import { useRouter, usePathname } from "next/navigation";
import { sendActiveGroupMessage } from "@/app/groupChat/groupMessage";
import { useEffect, useState } from "react";

const Header = () => {
  const { isLoggedIn, loading } = useAuth();
  const router = useRouter(); // 
  const pathname = usePathname(); // 
  const { sendMessage } = useWebSocket();
 

  return (
    <header className={style.header}>
      <div className={style.logoCont}>
        <Link href="/" className={style.homeLink}
        onClick={() => {
          if (pathname.startsWith("/groupChat")) {
            const groupId = sessionStorage.getItem("navigatedForwardToGroup");
            sendActiveGroupMessage("false", groupId, sendMessage);
            sessionStorage.removeItem("navigatedForwardToGroup");
          }
        }}
        
        >
          <HomeIcon className={style.homeIcon} />
          <span className={style.homeText}>Home</span>

    
        </Link>
      </div>
      <nav className={style.buttons}>


        {!isLoggedIn && !loading && (
          <>
            <AuthButton text="Login" href="/login" />
            <AuthButton text="Sign Up" href="/signup" />
          </>
        )}
        {isLoggedIn && !loading && (
          <>
            <NotificationButton />
            <AuthButton text="Chat" href="/chat" />
          <LogoutButton />
          </>
        )}
      </nav>
    </header>
  );
};

export default Header;