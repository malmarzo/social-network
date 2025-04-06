"use client";
import { useEffect, useState } from "react";
import Link from "next/link";
import { useWebSocket } from "@/context/Websocket";
import { sendUsersInvitationListMessage } from "../groupChat/groupMessage";

export default function MyGroups() {
    console.log("Rendering MyGroups component...");
    const { addMessageHandler } = useWebSocket();
    const [ myGroups, setMyGroups] = useState(null);
    const { sendMessage } = useWebSocket(); 
    const [activeGroup, setActiveGroup] = useState(null);

    useEffect(() => {
        // Request my groups once when the component mounts
        const getMyGroups = () => {
            const myGroupsMsg = { type: "myGroups" };
            sendMessage(myGroupsMsg);
        };

        getMyGroups(); 

        // Adding message handler
        addMessageHandler("myGroups", (msg) => {
            if (!msg || msg.length === 0) {
                setMyGroups([]); // Set groups as empty
            } else {
                setMyGroups(msg);
            }
            //setMyGroups(msg);
        });
        addMessageHandler("groupMessage", () => {
            console.log("New message received, refreshing groups...");
            getMyGroups(); // Re-fetch the groups to update the list
        });
        // Cleanup function (optional but good practice)
       
    }, [addMessageHandler, sendMessage]); 

    return (
        
        <div style={{ padding: "20px", color: "white" }}>
    <h2>My Groups</h2>
    {myGroups && myGroups.my_groups.length === 0 ? (
        <p>You are not a member of any groups.</p>
    ) : (
        myGroups && (
            <ul>
                {myGroups.my_groups.map((group) => (
                    <li key={group.id} style={{ marginBottom: "10px" }}>
                         <Link href={`/groupChat/${group.id}`} style={{ color: "#1e90ff", textDecoration: "underline", cursor: "pointer" }}
                          onClick={() => {
                           // sendUsersInvitationListMessage(group.id, sendMessage);
                        }}
                         >
                            <strong>{group.title}</strong>
                        </Link> 
                      
                    </li>
                ))}
            </ul> 

        )
        
    )}
</div>

    );
}
