"use client";
import { useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";
import UsersList from "./userlist";
import { useRouter } from "next/navigation";  
import { sendInvitations } from "./sendInvitation";
import { fetchUsersData } from "./userlist";
import { useEffect } from "react";


export default function CreateGroup() {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [groupID, setGroupID] = useState(null);
    const [groupCreatorID, setGroupCreatorID] = useState(null);
    const [selectedUsers, setSelectedUsers] = useState([]);  // Store selected users here
    const router = useRouter();
    const [users, setUsers] = useState([]);
    useEffect(() => {
      const fetchUsers = async () => {
          const data = await fetchUsersData();  // Fetch users data
          setUsers(data);  // Set the fetched users data into state
      };

      fetchUsers();
  }, []);  // This will run only once when the component mounts
   
    const createGroup = async () => {
        if (!title.trim() || !description.trim()) {
            alert("Title and Description are required.");
            return;
        }
        if (selectedUsers.length == 0){
          alert("you need at least to invite one person");
          return;
        }
       
        const body = { title, description };
        const response = await invokeAPI("groups", body, "POST");
        
        if (response.code === 200) {
            console.log("Group created successfully:", response.group);
            setGroupID(response.group.id);
            setGroupCreatorID(response.group.creator_id);
            console.log(groupCreatorID);

           
           router.push(`/groupChat/${response.group.id}`);

            // Invite users automatically after creating the group
            if (selectedUsers.length > 0) {
                await sendInvitations(response.group.id, response.group.creator_id, selectedUsers);
            }
            
            alert("Group created and users invited successfully!");
        } else {
            console.log("Could not create the group");
        }
    };

    return (
        <div className="max-w-md mx-auto p-6 bg-gray-900 rounded-lg shadow-md text-white">
            <h2 className="text-xl font-bold mb-4 text-center">Create a Group</h2>

            <input 
                type="text" 
                placeholder="Title" 
                value={title} 
                onChange={(e) => setTitle(e.target.value)}
                className="w-full p-3 mb-3 text-white rounded-md border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />

            <input 
                type="text" 
                placeholder="Description" 
                value={description} 
                onChange={(e) => setDescription(e.target.value)}
                className="w-full p-3 mb-3 text-white rounded-md border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />

            {/* Pass selectedUsers state and setter to UsersList */}
            <UsersList 
                users = {users}
                selectedUsers={selectedUsers} 
                setSelectedUsers={setSelectedUsers} 
            />

            <button 
                onClick={createGroup}
                className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-3 px-4 rounded-md transition-all duration-200"
            >
                Create Group
            </button>
        </div>
    );
}
