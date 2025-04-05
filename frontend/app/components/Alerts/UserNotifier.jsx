import { useEffect } from "react";
import { useWebSocket } from "@/context/Websocket";

//Used this component in the layout.js file to notify users
const UserNotifier = () => {
  const { addMessageHandler } = useWebSocket();

  useEffect(() => {
    //Adding msg Handlers (set the msg type and the function to handle it)

    addMessageHandler("new_follow_request", (msg) => {
      alert(
        "You have a new follow request from " + msg.followRequest.senderNickname
      );
    });
  }, [addMessageHandler]);

  return null;
};

export default UserNotifier;
