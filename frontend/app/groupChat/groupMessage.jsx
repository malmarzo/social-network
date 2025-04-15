
export const sendGroupMessage = async (groupID,senderID,message, sendMessage) => {
        console.log("the function is functioning");
        const groupMsg = {
            type: "groupMessage",
            group_message: {
                group_id: groupID,
                sender_id: senderID,
                message:message,
            }
        };
        sendMessage(groupMsg);  
    
};



export const sendUsersInvitationListMessage = async (groupID, sendMessage) => {
    console.log("the function is functioning");
    const usersInvitationListMsg = {
        type: "usersInvitationListMessage",
       users_invitation_list_message: {
            group_id: parseInt(groupID, 10),
        }
    };
    console.log("Sending WebSocket message:", usersInvitationListMsg);
    sendMessage(usersInvitationListMsg);  

};



export const sendGroupMembersMessage = async (groupID, sendMessage) => {
    console.log("the function is functioning");
    const groupMembersMsg = {
        type: "groupMembersMessage",
       group_members_message: {
            group_id: parseInt(groupID, 10),
        }
    };
    console.log("Sending WebSocket message:",  groupMembersMsg);
    sendMessage( groupMembersMsg);  

};



export const sendActiveGroupMessage = async (status,groupID, sendMessage) => {
    console.log("the function is functioning");
    const activeGroupMsg = {
        type: "activeGroupMessage",
       
        active_group_message: {
            status:status,
            group_id: parseInt(groupID),
        }
    };
    console.log("Sending message:", JSON.stringify(activeGroupMsg));
    sendMessage(activeGroupMsg);  

};




export const sendResetCountMessage = async (groupID, sendMessage) => {
    console.log("the function is functioning");
    const resetCountMsg = {
        type: "resetCountMessage",
        reset_count_message: {
            group_id: parseInt(groupID),
        }
    };
    sendMessage(resetCountMsg);  

};


export const handleRequestJoin = async ( groupID,groupCreator,currentUser,sendMessage) => {
        console.log("the function is functioning");
        const requestMsg = {
            type: "request",
            content: "",
            request: {
                group_id: groupID,
                group_creator: groupCreator,
                user_id:currentUser,
            },
        };
        sendMessage(requestMsg);  
};
