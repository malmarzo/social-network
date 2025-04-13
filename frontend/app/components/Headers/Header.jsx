import React from "react";
import style from "./Header.module.css";
import AuthButton from "../Buttons/AuthButtons";
import LogoutButton from "@/app/logout/page";
import { useAuth } from "@/context/AuthContext";
import { sendActiveGroupMessage } from "@/app/groupChat/groupMessage";
import { useWebSocket } from "@/context/Websocket";
//import { usePathname } from "next/navigation";
import { useRouter, usePathname } from "next/navigation";
import Link from "next/link";
  

const Header = () => {
  const { isLoggedIn, loading } = useAuth();
  const router = useRouter(); // ✅ Required to navigate programmatically
  const pathname = usePathname(); // ✅ Required to know which page you're on
  const groupId = sessionStorage.getItem("navigatedForwardToGroup");
  const { sendMessage } = useWebSocket();

  console.log(isLoggedIn);
  return (
    <header className={style.header}>
      <div className={style.logoCont}>
        <div className={style.logo}></div>
      </div>
      <div className={style.buttons}>

       {pathname.startsWith("/groupChat") && (
    <button
      onClick={() => {
        router.replace("/myGroups");
        sendActiveGroupMessage("false",groupId,sendMessage);
        sessionStorage.removeItem("navigatedForwardToGroup");
      }}
    >
      ← Back
    </button>
    )} 


        {!isLoggedIn && !loading && (
          <>
            <AuthButton text="Login" href="/login" />
            <AuthButton text="Signup" href="/signup" />
          </>
        )}

        {isLoggedIn && !loading && <LogoutButton />}
      </div>
      
    </header>
    
  );
};

export default Header;
