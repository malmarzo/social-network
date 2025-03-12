"use client";
import { useEffect, useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";

export default function GroupsToJoin() {
    console.log("Rendering GroupsToJoin component...");
    const [groups, setGroups] = useState([]);
    

    // Fetch groups when component loads
    useEffect(() => {
        async function fetchGroups() {
            try {
                const data = await invokeAPI("groups/request", null, "GET");
                if (Array.isArray(data)) {
                    //return data;
                    setGroups(data);
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

    // Handle Request to Join
    // const handleRequestJoin = async (groupId) => {
    //     try {
    //         const response = await fetch("/api/request-join", {
    //             method: "POST",
    //             headers: {
    //                 "Content-Type": "application/json",
    //             },
    //             body: JSON.stringify({ group_id: groupId }),
    //         });

    //         if (!response.ok) {
    //             throw new Error("Failed to request to join group");
    //         }

    //         alert("Request sent successfully!");
    //     } catch (err) {
    //         alert("Error: " + err.message);
    //     }
    // };

  

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
                            {/* <button
                                onClick={() => handleRequestJoin(group.id)}
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
                            </button> */}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}
