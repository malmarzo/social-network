import { useWebSocket } from "@/context/Websocket";
import  DisplayInvitationCard from "../../createGroup/invitationCard"
import { useState, useEffect } from "react";
import  DisplayRequestCard from "../../requestGroup/RequestCard"
import Link from "next/link";
import MyGroups from "@/app/myGroups/page";



//Used this component in the layout.js file to notify users
const UserNotifier = () => {
  const { addMessageHandler } = useWebSocket();
  const [invitation, setInvitation] = useState(null);
  const [request, setRequest] = useState(null);
  
 

  

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

    addMessageHandler("request", (msg) => {
      //alert(msg.content);
      setRequest(msg);
      
     
    });

    addMessageHandler("hello", (msg) => {
      alert(msg.content);
    });
  }, [addMessageHandler]);

  return (
   
    <div>
      {/* this id for the invitation card */}
      <>
        {invitation && (
            <DisplayInvitationCard invitation={invitation} onRespond={(userId, accepted) => {
        console.log(`User ${userId} ${accepted ? "accepted" : "declined"} the invitation`);
        //console.log("hello", invitation);
        setInvitation(null); // Remove invitation after response
    }}  />
        )}
        </>
        
        <>
        {/* this for the request card  */}
        {request && (
          <DisplayRequestCard
            request={request}
            onRespond={(userId, accepted) => {
              console.log(
                `User ${userId} ${accepted ? "accepted" : "declined"} the request`
              );
              setRequest(null); // Remove request after response
            }}
          />
        )}
        </>
    </div>

);
};

export default UserNotifier;
