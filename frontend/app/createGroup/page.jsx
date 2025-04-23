"use client";
import { useState, useEffect } from "react";
import { invokeAPI } from "@/utils/invokeAPI";
import UsersList from "./userlist";
import { useRouter } from "next/navigation";
import { sendInvitations } from "./sendInvitation";
import { fetchUsersData } from "./userlist";
import { useWebSocket } from "@/context/Websocket";

export default function CreateGroup() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [groupID, setGroupID] = useState(null);
  const [groupCreatorID, setGroupCreatorID] = useState(null);
  const [selectedUsers, setSelectedUsers] = useState([]);
  const router = useRouter();
  const [users, setUsers] = useState([]);
  const { sendMessage } = useWebSocket();
  const [errors, setErrors] = useState({});
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const fetchUsers = async () => {
      const data = await fetchUsersData();
      setUsers(data);
    };

    fetchUsers();
  }, []);

  const createGroup = async () => {
    if (isLoading) return; // Early return if already loading
    setIsLoading(true);

    try {
      const errors = {
        title: title.trim() ? "" : "Title is required",
        description: description.trim() ? "" : "Description is required",
      };

      const hasErrors = Object.values(errors).some((msg) => msg !== "");

      if (hasErrors) {
        setErrors(errors);
        setIsLoading(false); // Reset loading state on error
        return;
      }

      if (selectedUsers.length === 0) {
        alert("You need to invite at least one person");
        setIsLoading(false); // Reset loading state
        return;
      }

      const body = { title, description };
      const response = await invokeAPI("groups", body, "POST");

      if (response.code === 200) {
        setGroupID(response.group.id);
        setGroupCreatorID(response.group.creator_id);
        setTitle("");
        setDescription("");
        setSelectedUsers([]);

        // Send groups request message
        const GroupsToRequestMsg = { type: "groupsToRequest" };
        sendMessage(GroupsToRequestMsg);

        // Invite users if any are selected
        if (selectedUsers.length > 0) {
          await sendInvitations(
            selectedUsers,
            sendMessage,
            response.group.id,
            response.group.creator_id
          );
        }

        // Navigate to the new group
        router.push(`/groupChat/${response.group.id}`);
      } else {
        console.error("Could not create the group");
        alert("Failed to create group. Please try again.");
      }
    } catch (error) {
      console.error("Error creating group:", error);
      alert("An error occurred while creating the group");
    } finally {
      setIsLoading(false); // Always reset loading state
    }
  };

  return (
    <div className="w-full min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-100 to-blue-300">
      <div className="w-full max-w-2xl p-8 bg-white rounded-2xl shadow-2xl border border-blue-200">
        <h2 className="text-2xl font-extrabold text-blue-800 mb-6 text-center">
          Create a Group
        </h2>

        <input
          type="text"
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="w-full p-4 mb-4 text-blue-900 placeholder-blue-400 bg-blue-100 rounded-lg border border-blue-300 focus:outline-none focus:ring-2 focus:ring-blue-500 shadow-sm hover:shadow-md transition"
        />
        {errors.title && (
          <p className="text-base font-semibold text-red-600 mt-1">
            {errors.title}
          </p>
        )}

        <input
          type="text"
          placeholder="Description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          className="w-full p-4 mb-4 text-blue-900 placeholder-blue-400 bg-blue-100 rounded-lg border border-blue-300 focus:outline-none focus:ring-2 focus:ring-blue-500 shadow-sm hover:shadow-md transition"
        />
        {errors.description && (
          <p className="text-base font-semibold text-red-600 mt-1">
            {errors.description}
          </p>
        )}

        <div className="flex justify-center mt-4">
          <UsersList
            users={users}
            selectedUsers={selectedUsers}
            setSelectedUsers={setSelectedUsers}
          />
        </div>

        <button
          disabled={isLoading}
          onClick={createGroup}
          className={`w-full font-bold py-3 px-4 rounded-lg mt-6 shadow-md transition-all duration-200 ${
            isLoading
              ? "bg-blue-400 cursor-not-allowed"
              : "bg-blue-600 hover:bg-blue-700 hover:shadow-xl text-white"
          }`}
        >
          {isLoading ? "Creating Group..." : "Create Group"}
        </button>
      </div>
    </div>
  );
}
