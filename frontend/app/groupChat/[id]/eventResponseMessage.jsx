export const sendEventResponseMessage = async (groupID,eventID, userID, optionID,sendMessage) => {
    const eventResponseMsg = {
        type: "eventResponseMessage",
        event_Response_message: {
            group_id:groupID,
            event_id:eventID,
            sender_id:userID,
            option_id:optionID,
        }
    };
    sendMessage(eventResponseMsg); 
};