"use client";
import { useState } from "react";
import { invokeAPI } from "@/utils/invokeAPI";

export default function CreateGroup() {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    //const [creator_id, setCreator_id] = useState("");

    const createGroup = async () => {
        const body = { 
            title, 
            description 
        };
        const response = await invokeAPI("groups", body, "POST");
      //   if (!response) {
      //     throw new Error("No response received from API.");
      // }

      // if (!response.ok) {
      //     throw new Error("API request failed with status");
      // }

      // const data = await response.json();
      // console.log("API response data:", data);
    };

    return (
        <div className="max-w-md mx-auto p-6 bg-gray-900 rounded-lg shadow-md text-white">
        <h2 className="text-xl font-bold mb-4 text-center">Create a Group</h2>
        
        <input 
          type="text" 
          placeholder="Title" 
          value={title} 
          onChange={(e) => setTitle(e.target.value)}
          className="w-full p-3 mb-3 text-white rounded-md border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      
        <input 
          type="text" 
          placeholder="Description" 
          value={description} 
          onChange={(e) => setDescription(e.target.value)}
          className="w-full p-3 mb-3 text-white rounded-md border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      
        <button 
          onClick={createGroup}
          className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-3 px-4 rounded-md transition-all duration-200"
        >
          Create
        </button>
      </div>
      
    );
}
