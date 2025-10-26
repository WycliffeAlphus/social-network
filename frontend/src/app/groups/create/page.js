"use client";

import { useState } from "react";
import Link from "next/link"; 
import { useRouter } from "next/navigation"; 
export default function CreateGroupPage() {
  const [groupName, setGroupName] = useState("");
  const [groupDescription, setGroupDescription] = useState("");
  const [error, setError] = useState(null); 
  const [loading, setLoading] = useState(false); 
  const router = useRouter(); 

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null); 
    setLoading(true); 

    // Add client-side validation to prevent empty submissions
    if (!groupName || !groupDescription) {
        setError("Group Name and Description are required.");
        setLoading(false);
        return; 
    }

    try {
      const response = await fetch("http://localhost:8080/api/groups", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          title: groupName,        
          description: groupDescription, 
        }),
        // This is the key change for cookie-based authentication
        credentials: 'include', 
      });

      // Specific error handling for authentication and other issues
      if (response.status === 401) {
        throw new Error('You are not logged in. Please log in to create a group.');
      }

      if (!response.ok) {
        // Handle HTTP errors
        const errorData = await response.json(); 
        console.error("Server error details:", errorData);
        throw new Error(errorData.message || "Failed to create group");
      }

      const responseData = await response.json(); 
      console.log("Group created successfully:", responseData);
      alert("Group created successfully!");

      setGroupName("");
      setGroupDescription("");
      router.push('/groups'); 

    } catch (err) {
      console.error("Error creating group:", err);
      setError(err.message || "An unexpected error occurred.");
      // Redirect to login if a 401 error was caught
      if (err.message.includes('You are not logged in')) {
        router.push('/login');
      }
    } finally {
      setLoading(false); 
    }
  };

  return (
    <div className="flex min-h-screen justify-center items-center">
      <div className="bg-white p-8 border border-gray-400 shadow-lg w-full max-w-md">
        <Link href="/groups" className="hover:underline mb-4 block">
          &larr; Back to Groups
        </Link>

        <h1 className="text-2xl font-bold mb-6 text-center">Create New Group</h1>

        {error && (
          <div className="border border-black bg-gray-100 px-4 py-3 mb-4" role="alert">
            <strong className="font-bold">Error:</strong>
            <span className="block sm:inline"> {error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="groupName" className="block text-sm font-bold mb-2">
              Group Name *
            </label>
            <input
              type="text"
              id="groupName"
              className="border border-gray-400 w-full py-2 px-3 focus:border-black focus:outline-none"
              value={groupName}
              onChange={(e) => setGroupName(e.target.value)}
              placeholder="Enter group name"
              required
              disabled={loading}
            />
          </div>
          <div className="mb-6">
            <label htmlFor="groupDescription" className="block text-sm font-bold mb-2">
              Description *
            </label>
            <textarea
              id="groupDescription"
              className="border border-gray-400 w-full py-2 px-3 focus:border-black focus:outline-none h-32 resize-none"
              value={groupDescription}
              onChange={(e) => setGroupDescription(e.target.value)}
              placeholder="Enter group description"
              required
              disabled={loading}
            ></textarea>
          </div>
          <div className="flex gap-3">
            <button
              type="submit"
              className={`flex-1 font-bold py-2 px-4 transition-all duration-200 ${
                loading
                  ? 'border border-gray-400 text-gray-400 cursor-not-allowed'
                  : 'bg-black hover:bg-gray-800 text-white'
              }`}
              disabled={loading}
            >
              {loading ? 'Creating...' : 'Create Group'}
            </button>
            <Link href="/groups" className="flex-1">
              <button
                type="button"
                className="w-full border border-black hover:bg-black hover:text-white font-bold py-2 px-4 transition-all duration-200"
              >
                Cancel
              </button>
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}