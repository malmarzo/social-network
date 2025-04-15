import { invokeAPI } from "@/utils/invokeAPI";
import { useState, useEffect } from "react";
import { useWebSocket } from "@/context/Websocket";
import { sendGroupMembersMessage } from "../groupChat/groupMessage";

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

            
            onRespond(response.user_id, accepted);
            
            const getMyGroups = () => {
                const myGroupsMsg = { type: "myGroups" };
                sendMessage(myGroupsMsg);
            };
    
            getMyGroups();
            await  sendGroupMembersMessage(invitation.invite.group_id, sendMessage); 
    
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

const styles = {
    card: {
        background: "#fff",
        border: "1px solid #e2e8f0",
        padding: "15px",
        boxShadow: "0 2px 4px rgba(0, 0, 0, 0.1)",
        width: "100%", 
        borderRadius: "8px",
        marginBottom: "0.5rem", 
        fontFamily: "var(--font-geist-sans)",
    },
    acceptBtn: {
        background: "#16a34a", 
        color: "white",
        padding: "5px 10px",
        cursor: "pointer",
        marginRight: "5px",
        border: "none",
        borderRadius: "4px",
    },
    declineBtn: {
        background: "#dc2626", 
        color: "white",
        padding: "5px 10px",
        cursor: "pointer",
        border: "none",
        borderRadius: "4px",
    }
};
