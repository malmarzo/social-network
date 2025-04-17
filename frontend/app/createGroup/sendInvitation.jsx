// sendInvitations.js
import { invokeAPI } from "@/utils/invokeAPI";
import { useWebSocket } from "@/context/Websocket";



export const sendInvitations = async (users, sendMessage, groupID,invitedBy) => {
    users.forEach((user) => {
        const inviteMsg = {
            type: "invite",
            content: "",
            invite: {
                group_id: groupID,
                user_id: user,  
                invited_by: invitedBy 
            }
        };
        sendMessage(inviteMsg);  
    });
};


export const getInvitations = async (sendMessage) => {
        const GetInviteMsg = {
            type: "getInvite",
        };
        sendMessage(GetInviteMsg);  
   
};

