//import { useState, useEffect } from "react";
import { invokeAPI } from "@/utils/invokeAPI"; // Ensure this function exists
import AuthButton from "../components/Buttons/AuthButtons";
export default function GroupsPage() {
    
    // Create a new group
    const createGroup = async () => {
       
    };

    return (
        <div>
        <AuthButton text="create group" href="/createGroup" />
        <AuthButton text="request group" href="/requestGroup" />
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
