
// export const sendGroupMessage = async (groupID,senderID,users,message, sendMessage) => {
//     // const { sendMessage } = useWebSocket();
//     users.forEach((user) => {
//         console.log("the function is functioning");
//         const groupMsg = {
//             type: "groupMessage",
//            //invited_user: user, // Ensure it's a single recipient ID
//            // content: content,
//             group_message: {
//                 group_id: groupID,
//                 sender_id: senderID,
//                 recevier_id: user.id,  // The user being invited
//                 message:message,
//             }
//         };
//         //console.log(user);
//         sendMessage(groupMsg);  
//     });
// };



export const sendGroupMessage = async (groupID,senderID,message, sendMessage) => {
        console.log("the function is functioning");
        const groupMsg = {
            type: "groupMessage",
           //invited_user: user, // Ensure it's a single recipient ID
           // content: content,
            group_message: {
                group_id: groupID,
                sender_id: senderID,
               // recevier_id: user.id,  // The user being invited
                message:message,
            }
        };
        //console.log(user);
        sendMessage(groupMsg);  
    
};



export const sendUsersInvitationListMessage = async (groupID, sendMessage) => {
    console.log("the function is functioning");
    const usersInvitationListMsg = {
        type: "usersInvitationListMessage",
       //invited_user: user, // Ensure it's a single recipient ID
       // content: content,
       users_invitation_list_message: {
            group_id: parseInt(groupID, 10),
        }
    };
    //console.log(user);
    console.log("Sending WebSocket message:", usersInvitationListMsg);
    sendMessage(usersInvitationListMsg);  

};
