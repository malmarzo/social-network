import React from "react";
import style from "@/styles/Header.module.css";
import AuthButton from "../Buttons/AuthButtons";
import LogoutButton from "@/app/logout/page";
import { useAuth } from "@/context/AuthContext";
import Link from "next/link";

const Header = () => {
  const { isLoggedIn, loading } = useAuth();
  return (
    <header className={style.header}>
      <div className={style.logoCont}>
        <div className={style.logo}></div>
      </div>
      <div className={style.nav}>
        {isLoggedIn && !loading && (
          <nav className={style.navLinks}>
            <Link href="/" className={style.navLink}>Home</Link>
            <Link href="/chat" className={style.navLink}>Chat</Link>
          </nav>
        )}
      </div>
      <div className={style.buttons}>
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
