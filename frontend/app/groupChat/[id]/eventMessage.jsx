export const sendEventMessage = async (groupID,senderID,title,description, dateTime, options,sendMessage) => {
   
    const eventMsg = {
        type: "eventMessage",
        event_message: {
            group_id: groupID,
             title:title,
             description:description,
             date_time:dateTime,
             options:options,
             sender_id: senderID,
        }
    };
    
    sendMessage(eventMsg);  

};