// sendInvitations.js
import { invokeAPI } from "@/utils/invokeAPI";
import { useWebSocket } from "@/context/Websocket";


// export const sendInvitations = async (groupId, invitedBy, users) => {
//     await Promise.all(
//         users.map((userId) =>
//             invokeAPI("groups/invite", {
//                 group_id: groupId,
//                 user_id: userId,
//                 invited_by: invitedBy,
//             }, "POST")
//         )
//     );
// };

export const sendInvitations = async (users, sendMessage, groupID,invitedBy) => {
    // const { sendMessage } = useWebSocket();
    users.forEach((user) => {
        console.log("the function is functioning");
        const inviteMsg = {
            type: "invite",
           //invited_user: user, // Ensure it's a single recipient ID
            content: "You are invited to join a group",
            invite: {
                group_id: groupID,
                user_id: user,  // The user being invited
                invited_by: invitedBy  // The user who is sending the invite
            }
        };
        console.log(user);
        sendMessage(inviteMsg);  // Send each invitation
    });
};


// function sendInvitation(recipientId) {
//     // let inviteMsg = {
//     //     type: "invite",
//     //     userDetails: { id: currentUserId, nickname: currentUserNickname },
//     //     recipientId: recipientId,
//     //     content: "You have been invited to join the group."
//     // };

//     const inviteMsg = {};
//     inviteMsg.type = "invite";
//     inviteMsg.recipientIds = recipientId;
//     inviteMsg.content = "inviting someone";
//     sendMessage(inviteMsg);
    
// }
