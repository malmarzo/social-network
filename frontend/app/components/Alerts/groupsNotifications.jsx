import { useWebSocket } from "@/context/Websocket";
import DisplayInvitationCard from "../../createGroup/invitationCard";
import { useState, useEffect } from "react";
import DisplayRequestCard from "../../requestGroup/RequestCard";

const GroupsNotifications = () => {
  const { addMessageHandler } = useWebSocket();
  const [invitations, setInvitations] = useState(""); // Array to hold invitations
  const [requests, setRequests] = useState(""); // Array to hold requests

  useEffect(() => {
    // Adding WebSocket message handlers for different types
    addMessageHandler("invite", (msg) => {
      //alert(msg.content);
      setInvitations(msg);
      
     
    });

   addMessageHandler("request", (msg) => {
      //alert(msg.content);
      setRequests(msg);
      
     
    });

  }, [addMessageHandler]);

  
  return (
   <div>
      {/* this id for the invitation card */}
      <>
        {invitations && (
            <DisplayInvitationCard invitation={invitations} onRespond={(userId, accepted) => {
        console.log(`User ${userId} ${accepted ? "accepted" : "declined"} the invitation`);
        //console.log("hello", invitation);
        setInvitations(null); // Remove invitation after response
    }}  />
        )}
        </>
        
        <>
        {/* this for the request card  */}
        {requests && (
          <DisplayRequestCard
            request={requests}
            onRespond={(userId, accepted) => {
              console.log(
                `User ${userId} ${accepted ? "accepted" : "declined"} the request`
              );
              setRequests(null); // Remove request after response
            }}
          />
        )}
        </>
    </div>
  );
};

export default GroupsNotifications;
