"use client";
import { useEffect, useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";
import { useWebSocket } from "@/context/Websocket";

export default function GroupsToJoin() {
    console.log("Rendering GroupsToJoin component...");
    const [groups, setGroups] = useState([]);
    const [currentUser, setCurrentUser] = useState([]);
     const { sendMessage } = useWebSocket();
    

    // Fetch groups when component loads
    useEffect(() => {
        async function fetchGroups() {
            try {
                const data = await invokeAPI("groups/list", null, "GET");
                if (Array.isArray(data.groups)) {
                    //return data;
                    setGroups(data.groups);
                    setCurrentUser(data.current_user)
                } else {
                    console.error("Error fetching groups:", data.error_msg);
                    return [];
                }
            } catch (error) {
                console.error("Error fetching error:", error);
                return [];
            }
        }

        fetchGroups();
    }, []);

   

    const handleRequestJoin = async ( groupID,groupCreator,currentUser) => {
        // const { sendMessage } = useWebSocket();
            console.log("the function is functioning");
            const requestMsg = {
                type: "request",
               //invited_user: user, // Ensure it's a single recipient ID
                content: "a user request to join a group",
                request: {
                    group_id: groupID,
                    group_creator: groupCreator,  // The user who is sending the invite
                    user_id:currentUser,
                },
            };
            sendMessage(requestMsg);  // Send each invitation
    };
    
    
  

    return (
        <div style={{  padding: "20px", color: "white" }}>
            
            <h2>Groups You Can Request to Join</h2>
            {groups.length === 0 ? (
                <p>No available groups to request.</p>
            ) : (
                <ul>
                    {groups.map((group) => (
                        <li key={group.id} style={{ marginBottom: "10px" }}>
                            <strong>{group.title}</strong>
                            <button
                                onClick={() => handleRequestJoin(group.id, group.creator_id, currentUser)}
                                style={{
                                    marginLeft: "10px",
                                    padding: "5px 10px",
                                    backgroundColor: "blue",
                                    color: "white",
                                    border: "none",
                                    cursor: "pointer",
                                }}
                            >
                                Request to Join
                            </button>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}
