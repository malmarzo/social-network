import { invokeAPI } from "@/utils/invokeAPI";
import { useState, useEffect } from "react";
import { sendGroupMembersMessage } from "../groupChat/groupMessage";
import { useWebSocket } from "@/context/Websocket";
// import { useWebSocket } from "@/context/Websocket";

export default function DisplayRequestCard({ request,onRespond }) {
    const [showCard, setShowCard] = useState(false);
    const { sendMessage } = useWebSocket();

    useEffect(() => {
        if (request) {
            setShowCard(true);
        }
    }, [request]);

    
    const handleResponse = async (accepted) => {
        try {
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
             await  sendGroupMembersMessage(request.request.group_id, sendMessage); 
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


export const getRequests = async (sendMessage) => {
    const GetRequestMsg = {
        type: "getRequest",
    };
    sendMessage(GetRequestMsg);  

};


// CSS styles as JS object
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
