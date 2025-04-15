// "use client";
// import { useEffect, useState } from "react";
// import { invokeAPI } from "@/utils/invokeAPI";
// import { useWebSocket } from "@/context/Websocket";



// export default function GroupsToJoin() {
//     console.log("Rendering GroupsToJoin component...");
//     const [groups, setGroups] = useState([]);
//     const [currentUser, setCurrentUser] = useState([]);
//      const { sendMessage } = useWebSocket();
//      const { addMessageHandler } = useWebSocket();
    
//     const getGroupsToRequest = () => {
//         const GroupsToRequestMsg = { type: "groupsToRequest" };
//         sendMessage(GroupsToRequestMsg);
//     };

//     useEffect(() => {
//         getGroupsToRequest(); 

//         // Adding message handler
//         addMessageHandler("groupsToRequest", (msg) => {
//             if (!msg.my_groups || msg.my_groups.length === 0) {
//                 setGroups([]); // Set groups as empty
//             } else {
//                 setGroups(msg.my_groups);
//             }
//             setCurrentUser(msg.userDetails.id)
//         });

//         // Cleanup function (optional but good practice)
//         return () => {
//             // Remove the message handler if your WebSocket context supports it
//         };
//     }, [addMessageHandler, sendMessage]); 
   

//     const handleRequestJoin = async ( groupID,groupCreator,currentUser) => {
//         // const { sendMessage } = useWebSocket();
//             console.log("the function is functioning");
//             const requestMsg = {
//                 type: "request",
//                //invited_user: user, // Ensure it's a single recipient ID
//                 content: "",
//                 request: {
//                     group_id: groupID,
//                     group_creator: groupCreator,  // The user who is sending the invite
//                     user_id:currentUser,
//                 },
//             };
//             sendMessage(requestMsg);  // Send each invitation
//     };
    
    
  

//     return (
//         <div style={{  padding: "20px", color: "white" }}>
            
//             <h2>Groups You Can Request to Join</h2>
//             {groups.length === 0 ? (
//                 <p>No available groups to request.</p>
//             ) : (
//                 <ul>
//                     {groups.map((group) => (
//                         <li key={group.id} style={{ marginBottom: "10px" }}>
//                             <strong>{group.title}</strong>
//                             <button
//                                 onClick={() => {
//                                     handleRequestJoin(group.id, group.creator_id, currentUser);
//                                     getGroupsToRequest();
//                                 }}
//                                 style={{
//                                     marginLeft: "10px",
//                                     padding: "5px 10px",
//                                     backgroundColor: "blue",
//                                     color: "white",
//                                     border: "none",
//                                     cursor: "pointer",
//                                 }}
//                             >
//                                 Request to Join
//                             </button>
//                         </li>
//                     ))}
//                 </ul>
//             )}
//         </div>
//     );
// }
