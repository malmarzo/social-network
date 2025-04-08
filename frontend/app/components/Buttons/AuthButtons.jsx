import React from "react";
import Link from "next/link";
import styles from "@/styles/AuthButtons.module.css";
import { ChatBubbleLeftRightIcon } from "@heroicons/react/24/outline";

const AuthButton = ({ text, href }) => {
  return (
    <button>
      <Link href={href} className={styles.button}>
        {href === "/chat" && (
          <ChatBubbleLeftRightIcon className="h-4 w-4 mr-1" />
        )}
        {text}
      </Link>
    </button>
  );
};

export default AuthButton;
