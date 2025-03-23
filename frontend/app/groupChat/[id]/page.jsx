
"use client";
import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { invokeAPI } from "@/utils/invokeAPI";
import UsersList from "../../createGroup/userlist";
import { sendInvitations } from "../../createGroup/sendInvitation";
import { useWebSocket } from "@/context/Websocket";
import { sendGroupMessage } from "../groupMessage";

export default function GroupChat() {
    const router = useRouter();
    const { id } = useParams();
    const [group, setGroup] = useState(null);
    const [selectedUsers, setSelectedUsers] = useState([]); // Track selected users
    const [users, setUsers] = useState(null);
     const { sendMessage } = useWebSocket();
     const [message, setMessage] = useState(""); // State for the input message
     const { addMessageHandler } = useWebSocket();
     const [messages, setMessages] = useState([]);
     const [messages1, setMessages1] = useState([]);
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
                
                if (response.group && response.group.chat_history) {
                    // console.log(response.group.chat_history); 
                    //setMessages1(response.group.chat_history); // Populate initial chat history
                    // response.group.chat_history.forEach((msg) => {
                    //     setMessages1((prev) => [...prev, msg]); // Append each message one by one
                    // });
                    //console.log("hellloooo",messages1);
                    setMessages(response.group.chat_history.map((msg) => ({ group_message: msg })));

                }
               
            } else {
                console.log("hellooooooo");
                router.push("/"); // Redirect if group not found
            }
        };
       // setMessages(group.group.chat_history)
          addMessageHandler("groupMessage", (msg) => {
            console.log("Received message:", msg); // Debug log
            setMessages((prev) => [...prev, msg]); // Append new messages
        });

        if (id) fetchGroup();
    }, [id, router,addMessageHandler]);
   //console.log("hellloooo",messages1);
    // Function to handle sending invitations
    const handleInviteUsers = async () => {
        if (!group || !group.group || !group.group.creator_id) {
            console.error("Group data is missing.");
            return;
        }

        try {
            await sendInvitations(selectedUsers,sendMessage,group.group.id, group.group.creator_id);
            alert("Invitations sent successfully!");
            setSelectedUsers([]); // Reset selection after sending
        } catch (error) {
            console.error("Error sending invitations:", error);
        }
    };

    if (!group) return <p>Loading group chat...</p>;

    const handleSendMessage = async () => {
        if (!message.trim()) return; // Don't send empty messages
        // Simulate sending the message (this can be replaced with a backend API call)
       sendGroupMessage(group.group.id,group.group.current_user,message,sendMessage)
    //    ,group.group.group_members

        setMessage(""); // Clear input after sending
    };
    
  
      
   
    
   
    return (
        <div className="max-w-3xl mx-auto p-6 bg-gradient-to-br from-gray-800 to-gray-900 rounded-2xl shadow-lg text-white border border-gray-700 backdrop-blur-lg">
            <h2 className="text-3xl font-extrabold mb-3 text-blue-400">{group.group.title}</h2>
            <p className="text-lg text-gray-300 italic">{group.group.description}</p>

            <div className="mt-5 p-4 bg-gray-800 rounded-lg shadow-md border border-gray-700">
                <p className="text-lg font-semibold text-gray-200">
                    👤 {group.group.firstname} {group.group.lastname}
                </p>
            </div>

            {/* Chat messages section */}
            {/* this section wher i will work */}
            <div className="mt-6 bg-gray-700 p-4 rounded-lg border border-gray-600">
                <h3 className="text-xl font-bold text-white">💬 Chat messages</h3>

                {/* Messages display */}
                    <div className="h-40 overflow-y-auto bg-gray-800 p-3 rounded-lg border border-gray-700 mt-2">
                    {messages.length > 0 ? (
                    messages.map((msg) => (
                        <div key={msg.group_message.id} className={`mb-3 ${msg.group_message.sender_id === group.group.current_user ? "text-right" : ""}`}>
                            <p className={`text-sm font-semibold ${msg.group_message.sender_id === group.group.current_user ? "text-green-400" : "text-blue-400"}`}>
                                {msg.group_message.sender_id === group.group.current_user ? "You" : msg.group_message.first_name}
                            </p>
                            <p className="text-sm text-white-300">{msg.group_message.message}</p>
                            <p className="text-xs text-white-500">{msg.group_message.date_time}</p>
                        </div>
    ))
) : (
    <p className="text-gray-400 italic">No messages yet...</p>
                    )}
                    {/* {messages.length > 0 ? (
                        messages.map((msg, index) => (
                            <div 
                                key={msg?.group_message?.id || `msg-${index}`} // Fallback key if id is missing
                                className={`mb-3 ${msg?.group_message?.sender_id === group.group.current_user ? "text-right" : ""}`}
                            >
                                <p className={`text-sm font-semibold ${msg?.group_message?.sender_id === group.group.current_user ? "text-green-400" : "text-blue-400"}`}>
                                    {msg?.group_message?.sender_id === group.group.current_user ? "You" : msg?.group_message?.first_name}
                                </p>
                                <p className="text-sm text-white-300">{msg?.group_message?.message}</p>
                                <p className="text-xs text-white-500">{msg?.group_message?.date_time}</p>
                            </div>
                        ))
                    ) : (
                        <p className="text-gray-400 italic">No messages yet...</p>
                    )} */}

                </div>


                {/* Input for typing a message */}
                <div className="mt-4 flex items-center space-x-3">
                    <input
                        type="text"
                        className="flex-1 p-2 rounded-lg bg-gray-800 text-white border border-gray-700"
                        placeholder="Type your message..."
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        onKeyUp={(e) => e.key === "Enter" && handleSendMessage()} // Send message on Enter
                    />
                    <button
                       onClick={handleSendMessage}
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                       disabled={!message.trim()}
                    >
                        Send
                    </button>
                </div>
            </div>


            {/* end of section where i will work  */}




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
