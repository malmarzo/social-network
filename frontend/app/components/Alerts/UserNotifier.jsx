import { useWebSocket } from "@/context/Websocket";
import  DisplayInvitationCard from "../../groups/invitationCard"
import { useState, useEffect } from "react";


//Used this component in the layout.js file to notify users
const UserNotifier = () => {
  const { addMessageHandler } = useWebSocket();
  const [invitation, setInvitation] = useState(null);

  

  useEffect(() => {
    //Adding msg Handlers (set the msg type and the function to handle it)

    addMessageHandler("newUser", (msg) => {
      alert("New user joined");
    });

    addMessageHandler("removeUser", (msg) => {
      alert("User left");
    });

    addMessageHandler("invite", (msg) => {
      //alert(msg.content);
      setInvitation(msg);
      
     
    });

    addMessageHandler("hello", (msg) => {
      alert(msg.content);
    });
  }, [addMessageHandler]);

  return (
    <div>
        {invitation && (
            <DisplayInvitationCard invitation={invitation} onRespond={(userId, accepted) => {
        console.log(`User ${userId} ${accepted ? "accepted" : "declined"} the invitation`);
        //console.log("hello", invitation);
        setInvitation(null); // Remove invitation after response
    }}  />
        )}
    </div>
);
};

export default UserNotifier;
