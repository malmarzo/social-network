"use client";
import { useEffect, useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI"; // API helper function
import Link from "next/link";
import { useWebSocket } from "@/context/Websocket";

export default function MyGroups() {
    console.log("Rendering MyGroups component...");
    const { addMessageHandler } = useWebSocket();
    const [ myGroups, setMyGroups] = useState(null);
    const { sendMessage } = useWebSocket(); 

    useEffect(() => {
        // Request my groups once when the component mounts
        const getMyGroups = () => {
            const myGroupsMsg = { type: "myGroups" };
            sendMessage(myGroupsMsg);
        };

        getMyGroups(); 

        // Adding message handler
        addMessageHandler("myGroups", (msg) => {
            setMyGroups(msg);
        });

        // Cleanup function (optional but good practice)
        return () => {
            // Remove the message handler if your WebSocket context supports it
        };
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
                        <Link href={`/groupChat/${group.id}`} style={{ color: "#1e90ff", textDecoration: "underline", cursor: "pointer" }}>
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
