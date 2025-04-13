// import React from "react";
// import style from "@/styles/Header.module.css";
// import AuthButton from "../Buttons/AuthButtons";
// import LogoutButton from "@/app/logout/page";
// import { useAuth } from "@/context/AuthContext";
// <<<<<<< HEAD
// import { sendActiveGroupMessage } from "@/app/groupChat/groupMessage";
// import { useWebSocket } from "@/context/Websocket";
// //import { usePathname } from "next/navigation";
// import { useRouter, usePathname } from "next/navigation";
// import Link from "next/link";
  

// const Header = () => {
//   const { isLoggedIn, loading } = useAuth();
//   const router = useRouter(); // ✅ Required to navigate programmatically
//   const pathname = usePathname(); // ✅ Required to know which page you're on
//   const groupId = sessionStorage.getItem("navigatedForwardToGroup");
//   const { sendMessage } = useWebSocket();

//   console.log(isLoggedIn);
// =======
// import Link from "next/link";
// import { HomeIcon } from "@heroicons/react/24/outline";
// import NotificationButton from "./NotificationButton";

// const Header = () => {
//   const { isLoggedIn, loading } = useAuth();
// >>>>>>> origin/master
//   return (
//     <header className={style.header}>
//       <div className={style.logoCont}>
//         <Link href="/" className={style.homeLink}>
//           <HomeIcon className={style.homeIcon} />
//           <span className={style.homeText}>Home</span>
//         </Link>
//       </div>
// <<<<<<< HEAD
//       <div className={style.buttons}>

//        {pathname.startsWith("/groupChat") && (
//     <button
//       onClick={() => {
//         router.replace("/myGroups");
//         sendActiveGroupMessage("false",groupId,sendMessage);
//         sessionStorage.removeItem("navigatedForwardToGroup");
//       }}
//     >
//       ← Back
//     </button>
//     )} 


// =======
//       <nav className={style.buttons}>
// >>>>>>> origin/master
//         {!isLoggedIn && !loading && (
//           <>
//             <AuthButton text="Login" href="/login" />
//             <AuthButton text="Sign Up" href="/signup" />
//           </>
//         )}
// <<<<<<< HEAD

//         {isLoggedIn && !loading && <LogoutButton />}
//       </div>
      
// =======
//         {isLoggedIn && !loading && (
//           <>
//             <NotificationButton />
//             <AuthButton text="Chat" href="/chat" />
//           <LogoutButton />
//           </>
//         )}
//       </nav>
// >>>>>>> origin/master
//     </header>
    
//   );
// };

// export default Header;


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
  //const groupId = sessionStorage.getItem("navigatedForwardToGroup");
  const { sendMessage } = useWebSocket();
  const [groupId, setGroupId] = useState(null);

  useEffect(() => {
    if (typeof window !== "undefined") {
      const storedGroupId = sessionStorage.getItem("navigatedForwardToGroup");
      setGroupId(storedGroupId);
    }
  }, []);

  return (
    <header className={style.header}>
      <div className={style.logoCont}>
        <Link href="/" className={style.homeLink}>
          <HomeIcon className={style.homeIcon} />
          <span className={style.homeText}>Home</span>
        </Link>
      </div>
      <nav className={style.buttons}>

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
      )};

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