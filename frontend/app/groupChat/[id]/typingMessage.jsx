export const sendTypingMessage = async (groupID,senderID, sendMessage) => {
    const typingMsg = {
        type: "typingMessage",
        typing_message: {
            group_id: groupID,
            sender_id: senderID,
        }
    };
    sendMessage(typingMsg);  

};