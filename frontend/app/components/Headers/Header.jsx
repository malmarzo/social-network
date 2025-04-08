import React from "react";
import style from "@/styles/Header.module.css";
import AuthButton from "../Buttons/AuthButtons";
import LogoutButton from "@/app/logout/page";
import { useAuth } from "@/context/AuthContext";
import Link from "next/link";
import { HomeIcon } from "@heroicons/react/24/outline";

const Header = () => {
  const { isLoggedIn, loading } = useAuth();
  return (
    <header className={style.header}>
      <div className={style.logoCont}>
        <Link href="/" className={style.homeLink}>
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
            <AuthButton text="Chat" href="/chat" />
          <LogoutButton />
          </>
        )}
      </nav>
    </header>
  );
};

export default Header;
