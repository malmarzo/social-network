export const sendEventMessage = async (groupID,senderID,title,description, dateTime, options,sendMessage) => {
    // 
    //console.log("the function event is functioning");
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
    //console.log(eventMsg);
    sendMessage(eventMsg);  

};