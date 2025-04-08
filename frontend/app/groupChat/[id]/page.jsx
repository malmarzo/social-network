
"use client";
import { useEffect, useState, useRef } from "react";
import { useRouter, useParams } from "next/navigation";
import { invokeAPI } from "@/utils/invokeAPI";
import UsersList from "../../createGroup/userlist";
import { sendInvitations } from "../../createGroup/sendInvitation";
import { useWebSocket } from "@/context/Websocket";
import { sendGroupMessage } from "../groupMessage";
import EmojiPicker from "emoji-picker-react";
import { sendTypingMessage } from "./typingMessage"
import { sendEventMessage } from "./eventMessage";
import { withOptions } from "tailwindcss/plugin";
import { sendEventResponseMessage } from "./eventResponseMessage";
import  PostsFeed  from "./postFeed"
import { sendUsersInvitationListMessage } from "../groupMessage";
import { sendGroupMembersMessage } from "../groupMessage";
export default function GroupChat() {
    const router = useRouter();
    const { id } = useParams();
    const [group, setGroup] = useState(null);
    const [selectedUsers, setSelectedUsers] = useState([]); 
    const [users, setUsers] = useState(null);
    const [members, setMembers] = useState(null);
     const { sendMessage } = useWebSocket();
     const [message, setMessage] = useState(""); 
     const { addMessageHandler } = useWebSocket();
     const [messages, setMessages] = useState([]);
     const [showEmojiPicker, setShowEmojiPicker] = useState(false);
     const [typingStatus, setTypingStatus] = useState(""); 
     const messagesEndRef = useRef(null);
     const [events, setEvents] = useState([]);
     const [title, setTitle] = useState("");
     const [description, setDescription] = useState("");
     const [dateTime, setDateTime] = useState("");
     const [options, setOptions] = useState([""]);
     const [day, setDay] = useState("");
     const [selectedOption, setSelectedOption] = useState({});
     const [isUsersListVisible, setIsUsersListVisible] = useState(false); 
     const [isMembersListVisible, setIsMembersListVisible] = useState(false);
     const [activeSection, setActiveSection] = useState("events");
     const [showEventForm, setShowEventForm] = useState(false);
      const [errors, setErrors] = useState({});


    

     const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    useEffect(() => {
        scrollToBottom();
      
    }, [messages]);

  
      
    useEffect(() => {

     
        const fetchGroup = async () => {
            if (!id) {
                console.error("Group ID is undefined.");
                return;
            }    
            
            const response = await invokeAPI(`groups/chat/${id}`, null, "GET");
            console.log("Response status:", response.status);

            if (response.code === 200) {
                
                setGroup(response);
                
                
                if (response.group && response.group.chat_history) {
                    setMessages(response.group.chat_history.map((msg) => ({ group_message: msg })));
                }

                if (response.group && response.group.event_history) {
                    setEvents(response.group.event_history.map((msg) => ({ event_message: msg })));
                }
                
                if (response.group && response.group.event_responses_history) {
                    console.log("hello");
                    const eventResponses = response.group.event_responses_history.map((msg) => ({
                        option_id: msg.option_id,
                        sender_id: msg.sender_id,
                        first_name: msg.first_name,
                    }));
                
                    console.log("Event Responses:", eventResponses); // Debugging the mapped event responses
                
                    // Now set the eventResponses
                    setSelectedOption((prev) => {
                        // Initialize an accumulator object to update the selections
                        const updatedSelections = eventResponses.reduce((acc, { option_id, sender_id, first_name }) => {
                            // If the senderId already exists for this option_id, update the entry; otherwise, add it
                            const existingUsers = acc[option_id] || [];
                            const updatedOption = existingUsers.filter((user) => user.senderId !== sender_id);
                            updatedOption.push({ senderId: sender_id, firstName: first_name });
                
                            // Update the accumulator object
                            acc[option_id] = updatedOption;
                            return acc;
                        }, {});
                
                        console.log("Updated Selections:", updatedSelections); // Debugging the updated selections
                
                        return updatedSelections;
                    });
                
                }
                
               
            } else {
                router.push("/"); // Redirect if group not found
            }
        };
       
          addMessageHandler("groupMessage", (msg) => {
            console.log("Received message:", msg); // Debug log
            setMessages((prev) => [...prev, msg]); // Append new messages
           
        });

        addMessageHandler("typingMessage", (msg) => {
            setTypingStatus(msg.typing_message.content);
            // Clear typing status after 3 seconds
            setTimeout(() => setTypingStatus(""), 1500);
        });

        addMessageHandler("eventMessage", (msg) => {
            setEvents((prev) => [...prev,msg]);
        });
       
        addMessageHandler("eventResponseMessage", (msg) => {
            const { option_id, sender_id, first_name } = msg.event_response_message;
        
            setSelectedOption((prev) => {
                // Create a new object where we first remove the user's previous selection
                const updatedSelections = Object.keys(prev).reduce((acc, key) => {
                    acc[key] = prev[key].filter((user) => user.senderId !== sender_id);  // Remove old selection for this user
                    return acc;
                }, {});
        
                // Add the user's new selection to the correct option
                return {
                    ...updatedSelections,
                    [option_id]: [...(updatedSelections[option_id] || []), { senderId: sender_id, firstName: first_name }],
                };
            });
        });

      
     
        addMessageHandler("usersInvitationList", (msg) => {
            setUsers(msg.users_invitation_list_message); // Assuming this is the structure of the response
          });
      
          addMessageHandler("groupMembers", (msg) => {
            setMembers(msg.users_invitation_list_message.users); // Assuming this is the structure of the response
          });


        if (id) fetchGroup();
    }, [id,router,sendMessage,addMessageHandler]);



   console.log("yasseer",users);

    const handleEmojiClick = (emojiObject) => {
        setMessage((prevMessage) => prevMessage + emojiObject.emoji); // Add emoji to message
    };


    const handleButtonClick = async () => {
        // Toggle visibility of the users list
         setIsUsersListVisible((prevState) => !prevState);
    
        // Only send the message if it's the first time showing the list
        if (!isUsersListVisible) {
           try {
            await  sendUsersInvitationListMessage(id, sendMessage);
           
        } catch (error) {
            console.error("Error sending usersinvitationlist:", error);
        } 
        }
      };

      const handleButtonClick2 = async () => {
        // Toggle visibility of the users list
         setIsMembersListVisible((prevState) => !prevState);
    
        // Only send the message if it's the first time showing the list
        if (!isMembersListVisible) {
           try {
            await  sendGroupMembersMessage(id, sendMessage);
           
        } catch (error) {
            console.error("Error sending groupMemberslist:", error);
        } 
        }
      };



      

    //console.log(users);
    // Function to handle sending invitations
    const handleInviteUsers = async () => {
        if (!group || !group.group || !group.group.creator_id) {
            console.error("Group data is missing.");
            return;
        }

        try {
            await sendInvitations(selectedUsers,sendMessage,group.group.id, group.group.creator_id);
            await  sendUsersInvitationListMessage(id, sendMessage);
            alert("Invitations sent successfully!");
            setSelectedUsers([]); // Reset selection after sending
        } catch (error) {
            console.error("Error sending invitations:", error);
        }
    };

    if (!group) return <p>Loading group chat...</p>;

    const handleSendMessage = async () => {
        if (!message.trim()) return; // Don't send empty messages
        // Simulate sending the message (this can be replaced with a backend API call)
       sendGroupMessage(group.group.id,group.group.current_user,message,sendMessage)
    //    ,group.group.group_members
   

        setMessage(""); // Clear input after sending
    };
    
    const handleTyping = () => {
        sendTypingMessage(group.group.id, group.group.current_user, sendMessage);
    };
      

  
    const handleAddOption = () => {
        setOptions([...options, '']); // Add new empty option
    };
   
    const handleOptionChange = (index, value) => {
        const newOptions = [...options];
        newOptions[index] = value;
        setOptions(newOptions);
    };

      // this code to only allow future dates for the events
      const now = new Date();
      const formattedDateTime = now.toISOString().slice(0, 16);
     
    const handleDateChange = (e) => {
        const selectedDateTime = e.target.value;
        setDateTime(selectedDateTime); 
    };
    //console.log(day);
    const handleSendEvent = () => {
        // Filter out empty options and map them to Option objects with an ID and Text field
        const formattedOptions = options
            .filter(option => option.trim() !== '')  // Remove empty strings
            .map((option, index) => ({ id: index + 1, text: option }));  // Create Option objects
    
        const selectedDateTime = new Date(dateTime);
        const now = new Date();
        const errors = {
            title: title.trim() ? "" : "Title is required",
            description: description.trim() ? "" : "Description is required",
            dateTime: dateTime.trim() ? "" : "Date & Time are required",
            option: formattedOptions.length >= 2 ? "" : "At least two valid options are required",
          };

          if (selectedDateTime < now) {
            setErrors((prev) => ({
              ...prev,
              dateTime: "You cannot select a past time",
            }));
            return;
          }
          const hasErrors = Object.values(errors).some((msg) => msg !== "");
          
          if (hasErrors) {
            setErrors(errors);
            return;
          }
        sendEventMessage(group.group.id, group.group.current_user, title, description, dateTime, formattedOptions, sendMessage);
        console.log("Event is sent");
        
        // Reset form fields after sending the event
        setErrors({});
        setTitle('');
        setDescription('');
        setDateTime('');
        setOptions(['']);
        setShowEventForm(false); 

    }

    const handleResponseEvent = (eventId, optionId, sendMessage) => {
        // sendEventResponseMessage(groupId, eventId, eventId, userId, optionId, sendMessage);
        sendEventResponseMessage(group.group.id, eventId, group.group.current_user, optionId, sendMessage);
    };
    
    


    
    
   
    return (
        <div className="max-w-3xl mx-auto p-6 bg-gradient-to-br from-gray-800 to-gray-900 rounded-2xl shadow-lg text-white border border-gray-700 backdrop-blur-lg">
            <h2 className="text-3xl font-extrabold mb-3 text-blue-400">{group.group.title}</h2>
            <p className="text-lg text-gray-300 italic">Description: {group.group.description}</p>

            <div className="mt-5 p-4 bg-gray-800 rounded-lg shadow-md border border-gray-700">
                <p className="text-lg font-semibold text-gray-200">
                    👤 {group.group.firstname} {group.group.lastname} (Admin)
                </p>
            </div>


            {/* displaying group memebers */}
            <div>
            {/* Button to toggle visibility of members list */}
            <button
                onClick={handleButtonClick2}
                className="mt-4 bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
            >
                {isMembersListVisible ? "Hide Members List" : "Show Members List"}
            </button>

            {/* Conditionally render the members list if visible */}
            {isMembersListVisible && (
                <div className="mt-4 p-4 bg-gray-800 rounded-lg shadow-md border border-gray-700">
                {members && members.length > 0 ? (
                    members.map((member, index) => (
                    <div key={index} className="member text-white mb-2">
                        <p>{member.nickname}</p>
                    </div>
                    ))
                ) : (
                    <p className="text-gray-400">No members found</p>
                )}
                </div>
            )}
            </div>
            {/* end of displaying group members */}

          
            {/* displaying Chat messages section */}
            <div className="mt-6 bg-gray-700 p-4 rounded-lg border border-gray-600">
                <h3 className="text-xl font-bold text-white">💬 Chat messages</h3>

                {/* Messages display */}
                     <div className="h-40 overflow-y-auto bg-gray-800 p-3 rounded-lg border border-gray-700 mt-2">
                    
                    {messages.length > 0 ? (
                    messages.map((msg) => (
                        msg.group_message ? (  // Add this check
                            <div key={msg.group_message.id} className={`mb-3 ${msg.group_message.sender_id === group.group.current_user ? "text-right" : ""}`}>
                                <p className={`text-sm font-semibold ${msg.group_message.sender_id === group.group.current_user ? "text-green-400" : "text-blue-400"}`}>
                                    {msg.group_message.sender_id === group.group.current_user ? "You" : msg.group_message.first_name}
                                </p>
                                <p className="text-sm text-white-300">{msg.group_message.message}</p>
                                <p className="text-xs text-white-500">{msg.group_message.date_time}</p>
                            </div>
                        ) : null  
                    ))
                ) : (
                    <p className="text-gray-400 italic">No messages yet...</p>
                )}
                    <div ref={messagesEndRef}></div>
                </div> 
                {/* end of message displays */}


                 {/* Display typing status */}
                 {typingStatus && <p className="mt-2 text-green-400 italic">{typingStatus}</p>}
                        {/* end of typing display */}


                {/* Input for typing a message */}
                <div className="mt-4 flex items-center space-x-3">
                    <input
                        type="text"
                        className="flex-1 p-2 rounded-lg bg-gray-800 text-white border border-gray-700"
                        placeholder="Type your message..."
                        value={message}
                        onChange={(e) => {
                            setMessage(e.target.value);
                            handleTyping();
                        }}
                        onKeyUp={(e) => e.key === "Enter" && handleSendMessage()} // Send message on Enter
                    />
                     <button onClick={() => setShowEmojiPicker(!showEmojiPicker)} className="bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded">
                        😊
                    </button>
                   
                    <button
                       onClick={handleSendMessage}
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                       disabled={!message.trim()}
                    >
                        Send
                    </button>
                </div>
                 {/* Show emoji picker when toggled */}
                 {showEmojiPicker && (
                    <div className="absolute mt-2">
                        <EmojiPicker onEmojiClick={handleEmojiClick} theme="dark" />
                    </div>
                )}
            </div>
            {/* end for input for typing messages */}

           


            {/* the section for the event creation  */}
            {/* Event Creation Section */}
           
                {showEventForm && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-70 backdrop-blur-sm">
                    <div className="bg-gray-800 p-6 rounded-lg w-full max-w-xl mx-auto">
                    <h3 className="text-xl font-bold text-white mb-4">🎉 Create Event</h3>

                    <div className="space-y-3">
                        <input
                        type="text"
                        placeholder="Event Title"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        className="w-full p-2 rounded-lg bg-gray-900 text-white border border-gray-700"
                        />
                         {errors.title && (
                            <p className="text-red-500 text-sm mt-1">{errors.title}</p>
                        )}
                        <textarea
                        placeholder="Event Description"
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        className="w-full p-2 rounded-lg bg-gray-900 text-white border border-gray-700"
                        />
                         {errors.description && (
                            <p className="text-red-500 text-sm mt-1">{errors.description}</p>
                        )}
                        <input
                        type="datetime-local"
                        value={dateTime}
                        min={formattedDateTime}
                        onChange={handleDateChange}
                        className="w-full p-2 rounded-lg bg-gray-900 text-white border border-gray-700"
                        />
                         {errors.dateTime && (
                            <p className="text-red-500 text-sm mt-1">{errors.dateTime}</p>
                        )}
                        <div>
                        <h4 className="text-white font-semibold">Options</h4>
                        {options.map((option, index) => (
                            <div key={index} className="flex items-center space-x-2 mt-2">
                            <input
                                type="text"
                                placeholder={`Option ${index + 1}`}
                                value={option}
                                onChange={(e) => handleOptionChange(index, e.target.value)}
                                className="flex-1 p-2 rounded-lg bg-gray-900 text-white border border-gray-700"
                            />
                            {errors.option && (
                                <p className="text-red-500 text-sm mt-1">{errors.option}</p>
                            )}
                            {index === options.length - 1 && (
                                <button
                                onClick={handleAddOption}
                                className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
                                >
                                Add
                                </button>
                            )}
                            </div>
                        ))}
                        </div>

                        <div className="flex justify-between mt-4">
                        <button
                            // onClick={handleSendEvent}
                            onClick={() => {
                                handleSendEvent();
                               // setShowEventForm(false);
                               
                              }}
                            
                            className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                        >
                            Create Event
                        </button>
                        <button
                            onClick={() => { setShowEventForm(false);
                                setTitle('');
                                setDescription('');
                                setDateTime('');
                                setOptions(['']);
                                setErrors({});
                            }}
                            className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded"
                        >
                            Cancel
                        </button>
                        </div>
                    </div>
                    </div>
                </div>
                )}


                {/* end of event creation */}

                        <br></br>
            {/* here i will put buttons to show sections either events or posts */}
            <div className="flex justify-center space-x-4 mb-6">
                <button
                    onClick={() => setActiveSection("events")}
                    className={`px-4 py-2 rounded-lg font-bold transition ${
                        activeSection === "events"
                            ? "bg-blue-600 text-white"
                            : "bg-gray-700 text-gray-300 hover:bg-gray-600"
                    }`}
                >
                    Events
                </button>
                <button
                    onClick={() => setActiveSection("posts")}
                    className={`px-4 py-2 rounded-lg font-bold transition ${
                        activeSection === "posts"
                            ? "bg-blue-600 text-white"
                            : "bg-gray-700 text-gray-300 hover:bg-gray-600"
                    }`}
                >
                    Posts
                </button>
            </div>
            {/* end of buttons to show sections */}


                {/* here i will display the event  */}
                {activeSection === "events" && (
                <div className="event-container p-4">
                     <h2 className="text-2xl font-bold mb-4  text-center">Events</h2>
            {/* toggele button to show event creation form  */}
            <button
            onClick={() => setShowEventForm(true)}
            className="bg-blue-600 text-white px-4 py-2 rounded-lg font-bold hover:bg-blue-700 transition"
            >
                Create New Event
            </button>
            {/* end of toggele button for event creation */}

           <br></br>
           <br></br>

            <div className="event-list space-y-3 h-[600px] overflow-y-auto border border-gray-200 rounded-lg bg-gray-50 shadow-sm p-3">
                {events.length > 0 ? (
                    events.map((event) => (
                        <div key={event.event_message.event_id} className="event-card p-3 border border-gray-300 rounded-lg shadow bg-white">
                            <h3 className="text-lg font-bold text-blue-500 truncate">{event.event_message.title}</h3>
                            <p className="text-sm text-gray-600 truncate">{event.event_message.description}</p>
                            <p className="text-sm text-gray-600 truncate">{event.event_message.first_name}</p>
                            <p className="text-xs text-gray-500 mt-1">
                                <span className="font-semibold">Date & Time:</span> {event.event_message.date_time}
                                <br />
                                <span className="font-semibold">Day:</span> {event.event_message.day}
                            </p>

<ul className="mt-2 text-sm text-gray-700 space-y-2 bg-gray-100 p-4 rounded-lg shadow-md">
            <span className="font-bold text-gray-800">Options:</span>
            {event.event_message.options.map((option, index) => (
                <li
                    key={index}
                    className="list-disc list-inside flex justify-between items-center p-2 bg-white rounded-md shadow-sm hover:bg-gray-50 transition"
                > 
                 <div>
                 <span className="text-gray-900">{option.text}</span>
                    {selectedOption[option.id]?.map((user) => (
                        <span key={user.senderId} className="text-blue-600 font-medium">
                             <br></br> {user.firstName}
                        </span>
                    ))}
                </div> 

                    <button
                        onClick={() => handleResponseEvent(event.event_message.event_id, option.id,sendMessage)}
                        className={`ml-4 px-3 py-1 rounded-lg text-white ${
                            selectedOption === option.id ? "bg-blue-500" : "bg-gray-500 hover:bg-blue-400"
                        }`}
                    >
                        {selectedOption === option.id ? "Selected" : "Choose"}
                    </button>
                </li>
            ))}
        </ul>
                        </div>
                    ))
                ) : (
                    <p className="text-gray-500">No upcoming events yet...</p>
                )}
            </div>
        </div>
        )}
                {/* end of displaying the event */}


                {/* here i will display the posts & comments */}
                {activeSection === "posts" && (
                <div>
        <PostsFeed isGroup={true} groupID={id}/>
      </div>
      )}
                {/* end of displaying posts and comments */}
                       
            {/* Users List for Invitations */}

            <div>
                {/* Button to toggle the visibility of the users list */}
                <button
                    onClick={handleButtonClick}
                    className="mt-4 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                >
                    {isUsersListVisible ? "Hide Users List" : "Show Users List"}
                </button>

                {/* Conditionally render the users list div */}
                 {isUsersListVisible && users && users.users && (
                    <div className="mt-6 p-4 bg-gray-800 rounded-lg shadow-md border border-gray-700">
                    <UsersList users={users.users} selectedUsers={selectedUsers} setSelectedUsers={setSelectedUsers} />
                    <button
                        onClick={handleInviteUsers}
                        className="mt-4 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                        disabled={selectedUsers.length === 0}
                    >
                        Invite Selected Users
                    </button>
                    </div>
                )} 
                </div>
            {/* end of users invitation list */}
        </div>
    );
}
