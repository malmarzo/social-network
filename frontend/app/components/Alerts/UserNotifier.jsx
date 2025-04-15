import { useWebSocket } from "@/context/Websocket";
import { useState, useEffect } from "react";
import { useNotification } from "@/context/NotificationContext";
import { useAuth } from "@/context/AuthContext";


//Used this component in the layout.js file to notify users
const UserNotifier = () => {
  const { addMessageHandler } = useWebSocket();
  const { showInfo } = useNotification();
  const { userID } = useAuth();

  useEffect(() => {
    //Adding msg Handlers (set the msg type and the function to handle it)

    addMessageHandler("newUser", (msg) => {
      alert("New user joined");
    });

    addMessageHandler("removeUser", (msg) => {
      alert("User left");
    });


    addMessageHandler("new_follow_request", (msg) => {
      showInfo("New Follow Request", {
        duration: 3000,
        position: "top-right",
        link: `/profile/${userID}`,
      });
    });

    addMessageHandler("new_chat_message", (msg) => {
      showInfo(`${msg.userDetails.nickname} sent you a message`, {
        duration: 3000,
        position: "top-right",
        link: `/chat`,
      });
    });

  }, [addMessageHandler]);

  return (
    <div>
     
    </div>

);
};

export default UserNotifier;

 