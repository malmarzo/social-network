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

 
    addMessageHandler("groupNotifier", (msg) => {

        if (msg.group_notifier.sender_id != userID){
                showInfo(`${msg.group_notifier.first_name} sent a message in a group called ${msg.group_notifier.group_name}`, {
                    duration: 3000,
                    position: "top-right",
                    //link: `/groupChat/${msg.group_notifier.id}`,
                });
           } 
    });
  
    addMessageHandler("inviteNotifier", (msg) => {
      showInfo(`${msg.content}`, {
        duration: 3000,
        position: "top-right",
        //link: `/chat`,
      });
      
    });

    addMessageHandler("requestNotifier", (msg) => {
      showInfo(`${msg.content}`, {
        duration: 3000,
        position: "top-right",
        //link: `/chat`,
      });
      
    });

    addMessageHandler("eventNotifier", (msg) => {
      showInfo(`${msg.content}`, {
        duration: 3000,
        position: "top-right",
        //link: `/chat`,
      });
      
    });

  }, [addMessageHandler]);

  return (
    <div>
     
    </div>

);
};

export default UserNotifier;

 