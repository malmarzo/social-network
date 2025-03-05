// "use client";
// import { useEffect, useState } from "react";
// import { useRouter, useParams } from "next/navigation";
// import { invokeAPI } from "@/utils/invokeAPI";

// export default function GroupChat() {
//     const router = useRouter();
//     const { id } = useParams();
//     const [group, setGroup] = useState(null);

//     useEffect(() => {
//         const fetchGroup = async () => {
//             if (!id) {
//                 console.error("Group ID is undefined.");
//                 return;
//             }    
            
//             //const response = await fetch(`http://localhost:8080/groups/chat/${id}`);
//             const response = await invokeAPI(`groups/chat/${id}`,null , "GET");
//             console.log("Response status:", response.status);
//             if (response.code === 200) {
//                 //const data = await response.json();
//                 setGroup(response);
//                 //console.log(group);
//             } else {
//                 console.log("hellooooooo");
//                 router.push("/"); // Redirect if group not found
//             }
//         };

//         if (id) fetchGroup();
//     }, [id, router]);

//     if (!group) return <p>Loading group chat...</p>;

//     return (
       
//         <div className="max-w-3xl mx-auto p-6 bg-gradient-to-br from-gray-800 to-gray-900 rounded-2xl shadow-lg text-white border border-gray-700 backdrop-blur-lg">
//     <h2 className="text-3xl font-extrabold mb-3 text-blue-400">{group.group.title}</h2>
//     <p className="text-lg text-gray-300 italic">{group.group.description}</p>

//     <div className="mt-5 p-4 bg-gray-800 rounded-lg shadow-md border border-gray-700">
//         <p className="text-lg font-semibold text-gray-200">
//             👤 {group.group.firstname} {group.group.lastname}
//         </p>
//     </div>

//     <div className="mt-6 bg-gray-700 p-4 rounded-lg border border-gray-600">
//         <h3 className="text-xl font-bold text-white">💬 Chat messages</h3>
//         <p className="text-gray-300 italic">Chat messages will appear here...</p>
//     </div>
// </div>

//     );
// }



"use client";
import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { invokeAPI } from "@/utils/invokeAPI";
import UsersList from "../../groups/userlist";
import { sendInvitations } from "../../groups/sendInvitation";

export default function GroupChat() {
    const router = useRouter();
    const { id } = useParams();
    const [group, setGroup] = useState(null);
    const [selectedUsers, setSelectedUsers] = useState([]); // Track selected users
    const [users, setUsers] = useState(null);

    useEffect(() => {
        const fetchGroup = async () => {
            if (!id) {
                console.error("Group ID is undefined.");
                return;
            }    
            
            const response = await invokeAPI(`groups/chat/${id}`, null, "GET");
            console.log("Response status:", response.status);

            if (response.code === 200) {
                setGroup(response);
                setUsers(response);
                
               
            } else {
                console.log("hellooooooo");
                router.push("/"); // Redirect if group not found
            }
        };

        if (id) fetchGroup();
    }, [id, router]);

    // Function to handle sending invitations
    const handleInviteUsers = async () => {
        if (!group || !group.group || !group.group.creator_id) {
            console.error("Group data is missing.");
            return;
        }

        try {
            await sendInvitations(group.group.id, group.group.creator_id, selectedUsers);
            alert("Invitations sent successfully!");
            setSelectedUsers([]); // Reset selection after sending
        } catch (error) {
            console.error("Error sending invitations:", error);
        }
    };

    if (!group) return <p>Loading group chat...</p>;

    return (
        <div className="max-w-3xl mx-auto p-6 bg-gradient-to-br from-gray-800 to-gray-900 rounded-2xl shadow-lg text-white border border-gray-700 backdrop-blur-lg">
            <h2 className="text-3xl font-extrabold mb-3 text-blue-400">{group.group.title}</h2>
            <p className="text-lg text-gray-300 italic">{group.group.description}</p>

            <div className="mt-5 p-4 bg-gray-800 rounded-lg shadow-md border border-gray-700">
                <p className="text-lg font-semibold text-gray-200">
                    👤 {group.group.firstname} {group.group.lastname}
                </p>
            </div>

            <div className="mt-6 bg-gray-700 p-4 rounded-lg border border-gray-600">
                <h3 className="text-xl font-bold text-white">💬 Chat messages</h3>
                <p className="text-gray-300 italic">Chat messages will appear here...</p>
            </div>

            {/* Users List for Invitations */}
            <div className="mt-6 p-4 bg-gray-800 rounded-lg shadow-md border border-gray-700">
                <UsersList users={users.users} selectedUsers={selectedUsers} setSelectedUsers={setSelectedUsers} />
                <button 
                    onClick={handleInviteUsers}
                    className="mt-4 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    disabled={selectedUsers.length === 0}
                >
                    Invite Selected Users
                </button>
            </div>
        </div>
    );
}
