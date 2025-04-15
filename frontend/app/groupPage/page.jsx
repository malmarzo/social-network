
// "use client"; 
// import AuthButton from "../components/Buttons/AuthButtons";
// import Link from "next/link";
// import { useState, useEffect } from "react";
// export default function GroupsPage() {
//     // const [showGroupsNotifier, setShowGroupsNotifier] = useState(false);
//     // const toggleNotifier = () => {
//     //     setShowGroupsNotifier((prev) => !prev);
//     // };
   
//     return (
//         <div>
//         <AuthButton text="create group" href="/createGroup" />
//         {/* <AuthButton text="request group" href="/requestGroup" /> */}
//         <Link href="/myGroups">
//         <button 
//         style={{ padding: "10px", backgroundColor: "#1e90ff", color: "white", border: "none", cursor: "pointer" }}>
//             My Groups
//         </button>
//         </Link>

//         <Link href="/requestGroup">
//         <button
//         //  onClick={getMyGroups} 
//         style={{ padding: "10px", backgroundColor: "#1e90ff", color: "white", border: "none", cursor: "pointer" }}>
//             Request Groups
//         </button>
//         </Link>

//         {/* Button to toggle the UserNotifier */}
//         {/* <button
//                 onClick={toggleNotifier}
//                 style={styles.button}>
//                 {showGroupsNotifier ? "Hide Notifications" : "Show Notifications"}
//             </button> */}

           
//         </div>
       
//     );
// }

// // Basic Styles
// const styles = {
//     container: { padding: "20px", textAlign: "center" },
//     button: { padding: "10px 15px", margin: "10px", backgroundColor: "#1e90ff", color: "white", border: "none", cursor: "pointer" },
// };