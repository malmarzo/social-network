
export const sendGroupMessage = async (groupID,senderID,message, sendMessage) => {
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
    const usersInvitationListMsg = {
        type: "usersInvitationListMessage",
       users_invitation_list_message: {
            group_id: parseInt(groupID, 10),
        }
    };
    sendMessage(usersInvitationListMsg);  

};



export const sendGroupMembersMessage = async (groupID, sendMessage) => {
    const groupMembersMsg = {
        type: "groupMembersMessage",
       group_members_message: {
            group_id: parseInt(groupID, 10),
        }
    };
    sendMessage( groupMembersMsg);  

};



export const sendActiveGroupMessage = async (status,groupID, sendMessage) => {
    const activeGroupMsg = {
        type: "activeGroupMessage",
       
        active_group_message: {
            status:status,
            group_id: parseInt(groupID),
        }
    };
    sendMessage(activeGroupMsg);  

};




export const sendResetCountMessage = async (groupID, sendMessage) => {
    const resetCountMsg = {
        type: "resetCountMessage",
        reset_count_message: {
            group_id: parseInt(groupID),
        }
    };
    sendMessage(resetCountMsg);  

};


export const handleRequestJoin = async ( groupID,groupCreator,currentUser,sendMessage) => {
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


export const getEvents = async (sendMessage) => {
    const GetEventsMsg = {
        type: "getEvents",
    };
    sendMessage( GetEventsMsg);  

};