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

    try {
      const response = await fetch("http://localhost:8080/api/groups", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          // If you're using cookies for authentication, they will be sent automatically by the browser.
          // If you're using a token (e.g., JWT) in localStorage, you'd add it here:
          // "Authorization": `Bearer ${localStorage.getItem('authToken')}`,
        },
        body: JSON.stringify({
          title: groupName,        
          description: groupDescription, 
        }),
      });

      if (!response.ok) {
        // Handle HTTP errors
        const errorData = await response.json(); // Assuming backend sends JSON error details
        throw new Error(errorData.message || "Failed to create group");
      }

      const responseData = await response.json(); // Assuming your backend returns some data on success
      console.log("Group created successfully:", responseData);
      alert("Group created successfully!");

     
      setGroupName("");
      setGroupDescription("");
      router.push('/groups'); 

    } catch (err) {
      console.error("Error creating group:", err);
      setError(err.message || "An unexpected error occurred.");
    } finally {
      setLoading(false); 
    }
  };

  return (
    <div className="flex min-h-screen justify-center items-center bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
        {}
        <Link href="/" className="text-blue-500 hover:underline mb-4 block">
          &larr; Back to Home
        </Link>

        <h1 className="text-2xl font-bold mb-6 text-center">Create New Group</h1>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
            <strong className="font-bold">Error:</strong>
            <span className="block sm:inline"> {error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="groupName" className="block text-gray-700 text-sm font-bold mb-2">
              Group Name:
            </label>
            <input
              type="text"
              id="groupName"
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              value={groupName}
              onChange={(e) => setGroupName(e.target.value)}
              required
              disabled={loading}
            />
          </div>
          <div className="mb-6">
            <label htmlFor="groupDescription" className="block text-gray-700 text-sm font-bold mb-2">
              Description:
            </label>
            <textarea
              id="groupDescription"
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline h-32 resize-none"
              value={groupDescription}
              onChange={(e) => setGroupDescription(e.target.value)}
              required
              disabled={loading} // Disable input while loading
            ></textarea>
          </div>
          <div className="flex items-center justify-between">
            <button
              type="submit"
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={loading} // Disable button while loading
            >
              {loading ? 'Creating...' : 'Create Group'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}