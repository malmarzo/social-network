

import { invokeAPI } from "@/utils/invokeAPI";
export default function UsersList({ users, selectedUsers, setSelectedUsers }) {
    // Toggle user selection
    const toggleUserSelection = (user) => {
        setSelectedUsers((prev) =>
            prev.includes(user.id) ? prev.filter((id) => id !== user.id) : [...prev, user.id]
        );
    };

    return (
        <div className="w-full max-w-md bg-gray-900 p-6 rounded-xl shadow-lg border border-gray-700">
    <h2 className="text-white text-xl font-semibold mb-4">Select Users to Invite</h2>
    
    <div className="max-h-64 overflow-y-auto divide-y divide-gray-700 rounded-md border border-gray-700">
        {users.map((user) => (
        <div
            key={user.id}
            onClick={() => toggleUserSelection(user)}
            className={`p-3 cursor-pointer transition duration-200 ease-in-out ${
            selectedUsers.includes(user.id)
                ? "bg-blue-600 text-white"
                : "bg-gray-800 text-gray-300 hover:bg-gray-700"
            }`}
        >
            {user.nickname}
        </div>
        ))}
    </div>
    </div>

    );
};


export const fetchUsersData = async () => {
    try {
        const data = await invokeAPI("groups/users", null, "GET");
        if (Array.isArray(data)) {
            return data;
        } else {
            console.error("Error fetching users:", data.error_msg);
            return [];
        }
    } catch (error) {
        console.error("Error fetching users:", error);
        return [];
    }
};

