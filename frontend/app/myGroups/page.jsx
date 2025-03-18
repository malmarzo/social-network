"use client";
import { useEffect, useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI"; // API helper function
import Link from "next/link";

export default function MyGroups() {
    console.log("Rendering MyGroups component...");
    const [groups, setGroups] = useState([]);

    // Fetch user groups when the component loads
    useEffect(() => {
        async function fetchGroups() {
            try {
                const data = await invokeAPI("groups/mygroups", null, "GET");
                if (Array.isArray(data)) {
                    setGroups(data);
                } else {
                    console.error("Error fetching groups:", data.error_msg);
                }
            } catch (error) {
                console.error("Error fetching groups:", error);
            }
        }
        fetchGroups();
    }, []);

    return (
        <div style={{ padding: "20px", color: "white" }}>
            <h2>My Groups</h2>
            {groups.length === 0 ? (
                <p>You are not a member of any groups.</p>
            ) : (
                <ul>
                    {groups.map((group) => (
                        <li key={group.id} style={{ marginBottom: "10px" }}>
                            {/* <strong>{group.title}</strong> */}
                            <Link href={`/groupChat/${group.id}`} style={{ color: "#1e90ff", textDecoration: "underline", cursor: "pointer" }}>
                                <strong>{group.title}</strong>
                            </Link>
                            {/* <p>{group.description}</p> */}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}
