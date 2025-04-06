import { invokeAPI } from "@/utils/invokeAPI";
import { useState, useEffect } from "react";
import { useWebSocket } from "@/context/Websocket";

export default function DisplayInvitationCard({ invitation,onRespond }) {
    const [showCard, setShowCard] = useState(false);
    const { sendMessage } = useWebSocket(); 

    useEffect(() => {
        if (invitation) {
            setShowCard(true);
        }
    }, [invitation]);

    const handleResponse = async (accepted) => {

       
        try {
            
            // Call the Golang API
            console.log(invitation.invite.invited_by);
          const response = await invokeAPI("groups/invitation", {
            "type": "invite",
                 "userDetails": {
                 "username": "john_doe"
                     },
                    "content": "You have a new invitation!",
                    "invite": {
                     "group_id": invitation.invite.group_id,
                     "user_id": invitation.invite.user_id,
                     "invited_by": invitation.invite.invited_by,
                     "accepted": accepted
                    }
                }
          
            , "POST");

            // Call parent function to remove the invitation from the list
            onRespond(response.user_id, accepted);
            // if (accepted) {
            //     // Request updated group list
            //     sendMessage({ type: "myGroups" });
            // }
            // Hide the card after responding
            const getMyGroups = () => {
                const myGroupsMsg = { type: "myGroups" };
                sendMessage(myGroupsMsg);
            };
    
            getMyGroups(); 
    
            setShowCard(false);
        } catch (error) {
            console.error("Failed to update invitation response:", error);
        }
    };
    
    if (!showCard || !invitation) return null;

    return (
        <div style={styles.card}>
            <p>{invitation.content}</p>
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
