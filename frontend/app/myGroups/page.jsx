// "use client";
// import { useEffect, useState } from "react";
// import Link from "next/link";
// import { useWebSocket } from "@/context/Websocket";
// import { sendUsersInvitationListMessage } from "../groupChat/groupMessage";

// export default function MyGroups() {
//     console.log("Rendering MyGroups component...");
//     const { addMessageHandler } = useWebSocket();
//     const [ myGroups, setMyGroups] = useState(null);
//     const { sendMessage } = useWebSocket();
//     const [activeGroup, setActiveGroup] = useState(null);
    
//     useEffect(() => {
//         // Request my groups once when the component mounts
        
//         const getMyGroups = () => {
//             const myGroupsMsg = { type: "myGroups" };
//             sendMessage(myGroupsMsg);
//         };

//         getMyGroups(); 

       
//         // Adding message handler
//         addMessageHandler("myGroups", (msg) => {
//                 setMyGroups(msg.my_groups);
//             //setMyGroups(msg);
//         });
//         addMessageHandler("groupMessage", () => {
//             console.log("New message received, refreshing groups...");
//             getMyGroups(); // Re-fetch the groups to update the list
//         });
//         // Cleanup function (optional but good practice)
       
//     }, [addMessageHandler, sendMessage]); 

//     return (
        
//         <div style={{ padding: "20px", color: "white" }}>
//     <h2>My Groups</h2>
//     {myGroups == "" ? (
//         <p>You are not a member of any groups.</p>
//     ) : (
//         myGroups && (
//             <ul>
//                 {myGroups.map((group) => (
//                     <li key={group.id} style={{ marginBottom: "10px" }}>
//                          <Link href={`/groupChat/${group.id}`} style={{ color: "#1e90ff", textDecoration: "underline", cursor: "pointer" }}
//                           onClick={() => {
//                            // sendUsersInvitationListMessage(group.id, sendMessage);
//                         }}
//                          >
//                             <strong>{group.title}</strong>
//                         </Link> 
                      
//                     </li>
//                 ))}
//             </ul> 

//         )
        
//     )}
// </div>

//     );
// }


"use client";
import { useEffect, useState, useRef } from "react";
import Link from "next/link";
import { useWebSocket } from "@/context/Websocket";
import { sendUsersInvitationListMessage } from "../groupChat/groupMessage";
import { sendActiveGroupMessage } from "../groupChat/groupMessage";
export default function MyGroups() {
    console.log("Rendering MyGroups component...");
    const { addMessageHandler } = useWebSocket();
    const { sendMessage } = useWebSocket();
    const [myGroups, setMyGroups] = useState(null);
    //const [activeGroup, setActiveGroup] = useState(null);
    const [groupCounts, setGroupCounts] = useState({});

    //const historyIndex = useRef((window.history.state && window.history.state.idx) || 0);

    

    useEffect(() => {
       
        const getMyGroups = () => {
            const myGroupsMsg = { type: "myGroups" };
            sendMessage(myGroupsMsg);
        };

        getMyGroups();
       
        addMessageHandler("myGroups", (msg) => {
            console.log("Received myGroups message:", msg);
            //setMyGroups(msg.my_groups); // Now it's an array
            setMyGroups(Array.isArray(msg.my_groups) ? msg.my_groups : []);
        });

        addMessageHandler("groupMessage", (msg) => {
            console.log("New message received, refreshing groups...");
            getMyGroups();
            //test
            // const { group_id, count } = msg.group_message;
            // if (!group_id) return;
        
            // setGroupCounts(prevCounts => ({
            //     ...prevCounts,
            //     [group_id]: count
            // }));
        

            


        });

    }, [addMessageHandler, sendMessage]);

    return (
        <div style={{ padding: "20px", color: "white" }}>
            {/* <h2>My Groups</h2> */}

            {myGroups?.length === 0 ? (
                <p style={{ color: "white" }}>You are not a member of any groups.</p>
            ) : (
                myGroups && (
                    <ul>
                        {myGroups.map((group) => (
                            <li key={group.id} style={{ marginBottom: "10px" }}>
                                <Link
                                    href={`/groupChat/${group.id}`}
                                    style={{ color: "#1e90ff", textDecoration: "underline", cursor: "pointer" }}
                                    onClick={() => {
                                        sendActiveGroupMessage("true",group.id,sendMessage);
                                        sessionStorage.setItem("navigatedForwardToGroup",group.id);
                                        sendResetCountMessage(group.id,sendMessage);
                                        // sendUsersInvitationListMessage(group.id, sendMessage);
                                    }}
                                >
                                    <strong>{group.title}</strong>
                                    {group.count >0 && (

                                    <span style={{ marginLeft: "8px", color: "yellow" }}>
                                    {/* ({groupCounts[group.id]}) */}
                                            {group.count}
                                    </span>
                                    )}
                                    {/* {groupCounts[group.id] > 0 && ( */}
                                        <span style={{ marginLeft: "8px", color: "yellow" }}>
                                            {/* ({groupCounts[group.id]}) */}
                                            
                                        </span>
                                    {/* )} */}

                                </Link>
                            </li>
                        ))}
                    </ul>
                )
            )}
        </div>
    );
}




export const sendResetCountMessage = async (groupID, sendMessage) => {
    console.log("the function is functioning");
    const resetCountMsg = {
        type: "resetCountMessage",
       
        reset_count_message: {
            group_id: parseInt(groupID),
        }
    };
    //console.log("Sending message:", JSON.stringify(activeGroupMsg));
    //console.log(user);
    sendMessage(resetCountMsg);  

};
