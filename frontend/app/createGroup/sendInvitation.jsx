// sendInvitations.js
import { invokeAPI } from "@/utils/invokeAPI";
import { useWebSocket } from "@/context/Websocket";



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

