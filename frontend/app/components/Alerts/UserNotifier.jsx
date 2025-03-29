import { useWebSocket } from "@/context/Websocket";
import  DisplayInvitationCard from "../../createGroup/invitationCard"
import { useState, useEffect } from "react";
import  DisplayRequestCard from "../../requestGroup/RequestCard"
import EventNotificationCard from "@/app/groupChat/[id]/eventNotificationCard";

//Used this component in the layout.js file to notify users
const UserNotifier = () => {
  const { addMessageHandler } = useWebSocket();
  const [invitations, setInvitations] = useState(""); // Array to hold invitations
  const [requests, setRequests] = useState(""); // Array to hold requests
  const [eventNotifications, setEventNotifications] = useState([]);
 

  

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
      setInvitations(msg);
      
     
    });

   addMessageHandler("request", (msg) => {
      //alert(msg.content);
      setRequests(msg);
      
     
    });

    addMessageHandler("eventNotificationMsg", (msg) => {
      //alert(msg.content);
      setEventNotifications((prevNotifications) => [...prevNotifications, msg]);
      
     
    });

    addMessageHandler("hello", (msg) => {
      alert(msg.content);
    });
  }, [addMessageHandler]);

  const handleDismissNotification = (index) => {
    setEventNotifications((prevNotifications) =>
      prevNotifications.filter((_, i) => i !== index) // Remove notification at the given index
    );
  };

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

        {eventNotifications.map((notification, index) => (
                <EventNotificationCard
                  key={index}
                  content={notification.content}
                  onDismiss={() => handleDismissNotification(index)} // Pass index for dismissal
                />
              ))}
        <>
       
        
        
        </>

    </div>

);
};

export default UserNotifier;
