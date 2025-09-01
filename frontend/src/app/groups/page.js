"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { UserPlusIcon, CheckIcon, ClockIcon, UserGroupIcon } from '@heroicons/react/24/outline';

export default function GroupsPage() {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [currentUser, setCurrentUser] = useState(null);
  const [joinRequests, setJoinRequests] = useState({});
  const [pendingRequests, setPendingRequests] = useState({});
  const router = useRouter();

  useEffect(() => {
    async function fetchCurrentUser() {
      try {
        const response = await fetch("http://localhost:8080/api/profile/currentuser", {
          credentials: 'include',
        });
        if (response.ok) {
          const userData = await response.json();
          setCurrentUser(userData);
        }
      } catch (e) {
        console.error("Failed to fetch current user:", e);
      }
    }

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

    fetchCurrentUser();
    fetchGroups();
  }, [router]);

  const handleJoinRequest = async (groupId) => {
    try {
      setPendingRequests(prev => ({ ...prev, [groupId]: true }));

      const response = await fetch(`http://localhost:8080/api/groups/${groupId}/join`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to send join request');
      }

      const result = await response.json();
      setJoinRequests(prev => ({ ...prev, [groupId]: 'pending' }));
      alert('Join request sent successfully!');

    } catch (error) {
      console.error('Error sending join request:', error);
      alert(`Error: ${error.message}`);
    } finally {
      setPendingRequests(prev => ({ ...prev, [groupId]: false }));
    }
  };

  const isGroupCreator = (group) => {
    return currentUser && group.creator_id === currentUser.id;
  };

  const canJoinGroup = (group) => {
    return currentUser && group.creator_id !== currentUser.id && !joinRequests[group.id];
  };

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
        <div className="flex gap-3">
          <Link href="/groups/join-requests">
            <button className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition duration-300 ease-in-out flex items-center gap-2">
              <UserGroupIcon className="h-5 w-5" />
              Manage Requests
            </button>
          </Link>
          <Link href="/groups/create">
            <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-300 ease-in-out">
              Create New Group
            </button>
          </Link>
        </div>
      </div>

      {(!Array.isArray(groups) || groups.length === 0) ? (
        <p className="text-gray-600 text-center">No groups found. Start by creating a new one!</p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {groups.map((group) => (
            <div key={group.id} className="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300">
              <div className="flex justify-between items-start mb-3">
                <h2 className="text-xl font-semibold text-gray-900">{group.title}</h2>
                <span className={`px-2 py-1 text-xs rounded-full ${
                  group.privacy_setting === 'public' ? 'bg-green-100 text-green-800' :
                  group.privacy_setting === 'private' ? 'bg-yellow-100 text-yellow-800' :
                  'bg-red-100 text-red-800'
                }`}>
                  {group.privacy_setting}
                </span>
              </div>

              <p className="text-gray-600 mb-4">{group.description}</p>

              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-500">
                  Created: {new Date(group.created_at).toLocaleDateString()}
                </span>

                {isGroupCreator(group) ? (
                  <span className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
                    Your Group
                  </span>
                ) : canJoinGroup(group) ? (
                  <button
                    onClick={() => handleJoinRequest(group.id)}
                    disabled={pendingRequests[group.id]}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                      pendingRequests[group.id]
                        ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                        : 'bg-blue-500 hover:bg-blue-600 text-white hover:shadow-md'
                    }`}
                  >
                    {pendingRequests[group.id] ? (
                      <>
                        <ClockIcon className="h-4 w-4 animate-spin" />
                        Sending...
                      </>
                    ) : (
                      <>
                        <UserPlusIcon className="h-4 w-4" />
                        Join Group
                      </>
                    )}
                  </button>
                ) : joinRequests[group.id] === 'pending' ? (
                  <span className="bg-yellow-100 text-yellow-800 px-3 py-1 rounded-full text-sm font-medium flex items-center gap-1">
                    <ClockIcon className="h-4 w-4" />
                    Request Pending
                  </span>
                ) : (
                  <span className="bg-gray-100 text-gray-600 px-3 py-1 rounded-full text-sm font-medium">
                    Already Member
                  </span>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}