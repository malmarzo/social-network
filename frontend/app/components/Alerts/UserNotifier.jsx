import { useEffect } from "react";
import { useWebSocket } from "@/context/Websocket";


//Used this component in the layout.js file to notify users
const UserNotifier = () => {
  const { addMessageHandler } = useWebSocket();

  useEffect(() => {
    //Adding msg Handlers (set the msg type and the function to handle it)

    addMessageHandler("newUser", (msg) => {
      alert("New user joined");
    });

    addMessageHandler("removeUser", (msg) => {
      alert("User left");
    });

    addMessageHandler("hello", (msg) => {
      alert(msg.content);
    });
  }, [addMessageHandler]);

  return null;
};

export default UserNotifier;
