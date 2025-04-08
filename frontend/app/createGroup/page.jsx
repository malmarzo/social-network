"use client";
import { useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";
import UsersList from "./userlist";
import { useRouter } from "next/navigation";  
import { sendInvitations } from "./sendInvitation";
import { fetchUsersData } from "./userlist";
import { useEffect } from "react";
import { useWebSocket } from "@/context/Websocket";


export default function CreateGroup() {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [groupID, setGroupID] = useState(null);
    const [groupCreatorID, setGroupCreatorID] = useState(null);
    const [selectedUsers, setSelectedUsers] = useState([]);  // Store selected users here
    const router = useRouter();
    const [users, setUsers] = useState([]);
    const { sendMessage } = useWebSocket();
    const [errors, setErrors] = useState({});
    useEffect(() => {
      const fetchUsers = async () => {
          const data = await fetchUsersData();  // Fetch users data
          setUsers(data);  // Set the fetched users data into state
      };

      fetchUsers();
  }, []);  // This will run only once when the component mounts
   
    const createGroup = async () => {
        // if (!title.trim() || !description.trim()) {
        //     alert("Title and Description are required.");
        //     return;
        // }
        const errors = {
            title: title.trim() ? "" : "Title is required",
            description: description.trim() ? "" : "Description is required",
          };

       
       
        const hasErrors = Object.values(errors).some((msg) => msg !== "");
          
        if (hasErrors) {
          setErrors(errors);
          return;
        }
        const body = { title, description };
        if (selectedUsers.length == 0){
            alert("you need at least to invite one person");
            return;
          }
        const response = await invokeAPI("groups", body, "POST");

      
        
        if (response.code === 200) {
            //setErrors({});
            console.log("Group created successfully:", response.group);
            setGroupID(response.group.id);
            setGroupCreatorID(response.group.creator_id);
            console.log(groupCreatorID);
            const getGroupsToRequest = () => {
                const GroupsToRequestMsg = { type: "groupsToRequest" };
                sendMessage(GroupsToRequestMsg);
            };
            getGroupsToRequest();
           

           
           router.push(`/groupChat/${response.group.id}`);

            // Invite users automatically after creating the group
            if (selectedUsers.length > 0) {
                // await sendInvitations(response.group.id, response.group.creator_id, selectedUsers);
                //test;
                console.log("Selected Users:", selectedUsers);

                console.log("test for send invitations");
                await sendInvitations(selectedUsers,sendMessage,response.group.id,response.group.creator_id);
                //end of test
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
                className="w-full p-4 mb-4 text-white bg-gray-800 rounded-lg border border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-300 ease-in-out shadow-md hover:shadow-lg"
            />
             {errors.title && (
                <p className="text-lg font-bold text-red-600 mt-2">{errors.title}</p>
                )}
                <br></br>

            <input 
                type="text" 
                placeholder="Description" 
                value={description} 
                onChange={(e) => setDescription(e.target.value)}
                className="w-full p-4 mb-4 text-white bg-gray-800 rounded-lg border border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-300 ease-in-out shadow-md hover:shadow-lg"
            />
            {errors.description && (
                 <p className="text-lg font-bold text-red-600 mt-2">{errors.description}</p>
                    )}
                    <br></br>
            {/* Pass selectedUsers state and setter to UsersList */}
            <UsersList 
                users = {users}
                selectedUsers={selectedUsers} 
                setSelectedUsers={setSelectedUsers} 
            />
            <br></br>

            <button 
                // onClick={createGroup}

                onClick={() => {
                    createGroup();
                    // getGroupsToRequest();
                }}
                className="w-full bg-blue-600 hover:bg-indigo-700 text-white font-bold py-3 px-4 rounded-md transition-all duration-200"
            >
                Create Group
            </button>
        </div>
    );
}
