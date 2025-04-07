import { invokeAPI } from "@/utils/invokeAPI";
import { useState, useEffect } from "react";
import { sendGroupMembersMessage } from "../groupChat/groupMessage";
// import { useWebSocket } from "@/context/Websocket";

export default function DisplayRequestCard({ request,onRespond }) {
    const [showCard, setShowCard] = useState(false);
    // const { sendMessage } = useWebSocket();

    useEffect(() => {
        if (request) {
            setShowCard(true);
        }
    }, [request]);

    
    const handleResponse = async (accepted) => {

       
        try {
            
            // Call the Golang API
           // console.log(invitation.invite.invited_by);
           console.log( request.request.group_id);
           console.log(request.request.user_id);
          const response = await invokeAPI("groups/request", {
            "type": "request",
                    "request": {
                     "group_id": request.request.group_id,
                     "group_creator":request.request.group_creator,
                     "user_id": request.request.user_id,
                     "accepted": accepted
                    }
                }
          
            , "POST");

            // Call parent function to remove the invitation from the list
            onRespond(response.user_id, accepted);
            // not tested 
             await  sendGroupMembersMessage(request.request.group_id, sendMessage); 
            // if (!accepted) {
            //     // Request updated group list
            //     const getGroupsToRequest = () => {
            //         const GroupsToRequestMsg = { type: "groupsToRequest" };
            //         sendMessage(GroupsToRequestMsg);
            //     };
            //     getGroupsToRequest();
            // }

            // Hide the card after responding
            setShowCard(false);
        } catch (error) {
            console.error("Failed to update invitation response:", error);
        }
    };
    
    if (!showCard || !request) return null;

    return (
        <div style={styles.card}>
            <p>{request.content}</p>
            <button style={styles.acceptBtn} onClick={() => handleResponse(true)}>Accept</button>
            <button style={styles.declineBtn} onClick={() => handleResponse(false)}>Decline</button>
        </div>
    );
}

// CSS styles as JS object
const styles = {
    card: {
        background: "#fff",
        border: "1px solid black",
        padding: "15px",
        boxShadow: "2px 2px 10px rgba(0, 0, 0, 0.2)",
        position: "fixed",
        bottom: "20px",
        right: "20px",
        width: "250px",
        borderRadius: "8px",
    },
    acceptBtn: {
        background: "green",
        color: "white",
        padding: "5px 10px",
        cursor: "pointer",
        marginRight: "5px",
    },
    declineBtn: {
        background: "red",
        color: "white",
        padding: "5px 10px",
        cursor: "pointer",
    }
};
