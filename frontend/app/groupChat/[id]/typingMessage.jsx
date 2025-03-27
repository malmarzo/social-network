export const sendTypingMessage = async (groupID,senderID, sendMessage) => {
    console.log("the function typing is functioning");
    const typingMsg = {
        type: "typingMessage",
        typing_message: {
            group_id: groupID,
            sender_id: senderID,
        }
    };
    //console.log(user);
    sendMessage(typingMsg);  

};