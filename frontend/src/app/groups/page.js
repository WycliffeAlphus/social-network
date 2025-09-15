"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function GroupsPage() {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const router = useRouter();

  useEffect(() => {
    async function fetchGroups() {
      try {
        const response = await fetch("http://localhost:8080/api/groups", {
          credentials: 'include',
        });

        if (response.status === 401) {
          throw new Error('You are not logged in. Please log in to view groups.');
        }

        if (!response.ok) {
          const errorData = await response.json();
          console.error("Server error details:", errorData);
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const responseData = await response.json();
        console.log("Received data from API:", responseData);

        if (responseData && Array.isArray(responseData.data)) {
          setGroups(responseData.data);
        } else {
          console.error("API did not return a valid groups array:", responseData);
          setGroups([]);
        }

      } catch (e) {
        setError(e.message);
        if (e.message.includes('You are not logged in')) {
          router.push('/login');
        }
      } finally {
        setLoading(false);
      }
    }
    fetchGroups();
  }, [router]);

  if (loading) {
    return (
      <div className="flex min-h-screen justify-center items-center bg-gray-100">
        <p className="text-xl text-gray-700">Loading groups...</p>
      </div>
    );
  }

  if (error && !loading) {
    return (
      <div className="flex min-h-screen justify-center items-center bg-gray-100">
        <p className="text-xl text-red-500">Error: {error}</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Groups</h1>
        <Link href="/groups/create">
          <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-300 ease-in-out">
            Create New Group
          </button>
        </Link>
      </div>

      {(!Array.isArray(groups) || groups.length === 0) ? (
        <p className="text-gray-600 text-center">No groups found. Start by creating a new one!</p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {groups.map((group) => (
            <div key={group.id} className="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300">
              <a href={`/groups/${group.id}/group`}>
                <h2 className="text-xl font-semibold mb-2 text-gray-900">{group.title}</h2>
              </a>
              <p className="text-gray-600">{group.description}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}