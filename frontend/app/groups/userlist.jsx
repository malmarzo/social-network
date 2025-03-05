
// import { useEffect, useState } from "react";
// import { invokeAPI } from "@/utils/invokeAPI"; 

// export default function UsersList({ selectedUsers, setSelectedUsers }) {
//     const [users, setUsers] = useState([]);

//     // Fetch available users
//     useEffect(() => {
//         const fetchUsers = async () => {
//             const data = await invokeAPI("groups/users", null, "GET");
//             if (Array.isArray(data)) {
//                 setUsers(data);
//                 console.log(users);
//             } else {
//                 console.error("Error fetching users:", data.error_msg);
//             }
//         };
//         fetchUsers();
//     }, []);

//     // Toggle user selection
//     const toggleUserSelection = (user) => {
//         setSelectedUsers((prev) =>
//             prev.includes(user.id) ? prev.filter((id) => id !== user.id) : [...prev, user.id]
//         );
//     };

//     return (
//         <div>
//             <h2>Select Users to Invite</h2>
//             <div style={{ maxHeight: "200px", overflowY: "scroll", border: "1px solid black", padding: "10px" }}>
//                 {users.map((user) => (
//                     <div
//                         key={user.id}
//                         onClick={() => toggleUserSelection(user)}
//                         style={{
//                             padding: "5px",
//                             cursor: "pointer",
//                             backgroundColor: selectedUsers.includes(user.id) ? "#bde0fe" : "black",
//                         }}
//                     >
//                         {user.nickname}
//                     </div>
//                 ))}
//             </div>
//         </div>
//     );
// }

import { invokeAPI } from "@/utils/invokeAPI";
export default function UsersList({ users, selectedUsers, setSelectedUsers }) {
    // Toggle user selection
    const toggleUserSelection = (user) => {
        setSelectedUsers((prev) =>
            prev.includes(user.id) ? prev.filter((id) => id !== user.id) : [...prev, user.id]
        );
    };

    return (
        <div>
            <h2>Select Users to Invite</h2>
            <div style={{ maxHeight: "200px", overflowY: "scroll", border: "1px solid black", padding: "10px" }}>
                {users.map((user) => (
                    <div
                        key={user.id}
                        onClick={() => toggleUserSelection(user)}
                        style={{
                            padding: "5px",
                            cursor: "pointer",
                            backgroundColor: selectedUsers.includes(user.id) ? "#bde0fe" : "black",
                        }}
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

