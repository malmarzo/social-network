
"use client"; 
import AuthButton from "../components/Buttons/AuthButtons";
import Link from "next/link";
export default function GroupsPage() {
   
    return (
        <div>
        <AuthButton text="create group" href="/createGroup" />
        {/* <AuthButton text="request group" href="/requestGroup" /> */}
        <Link href="/myGroups">
        <button
        //  onClick={getMyGroups} 
        style={{ padding: "10px", backgroundColor: "#1e90ff", color: "white", border: "none", cursor: "pointer" }}>
            My Groups
        </button>
        </Link>

        <Link href="/requestGroup">
        <button
        //  onClick={getMyGroups} 
        style={{ padding: "10px", backgroundColor: "#1e90ff", color: "white", border: "none", cursor: "pointer" }}>
            Request Groups
        </button>
        </Link>
        </div>
       
    );
}

// Basic Styles
const styles = {
    container: { padding: "20px", textAlign: "center" },
    createSection: { marginBottom: "20px" },
    input: { padding: "8px", marginRight: "10px" },
    buttonContainer: { marginBottom: "20px" },
    button: { padding: "10px 15px", margin: "5px", cursor: "pointer" },
    listContainer: { marginTop: "20px", border: "1px solid #ddd", padding: "10px" }
};
