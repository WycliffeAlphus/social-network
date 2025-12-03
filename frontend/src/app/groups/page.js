"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { UserPlusIcon, ClockIcon, UserGroupIcon, BellIcon } from '@heroicons/react/24/outline';

export default function GroupsPage() {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [currentUser, setCurrentUser] = useState(null);
  const [joinRequests, setJoinRequests] = useState({});
  const [pendingRequests, setPendingRequests] = useState({});
  const [groupMemberships, setGroupMemberships] = useState({});
  const [pendingJoinCounts, setPendingJoinCounts] = useState({});
  const [activeTab, setActiveTab] = useState('all'); // 'all' or 'my'
  const router = useRouter();

  useEffect(() => {
    async function fetchData() {
      try {
        // First, fetch current user
        const userResponse = await fetch("http://localhost:8080/api/profile/currentuser", {
          credentials: 'include',
        });
        let userData = null;
        if (userResponse.ok) {
          const rawUserData = await userResponse.json();
          // Extract the actual user ID from current_user_id field
          userData = {
            ...rawUserData,
            id: rawUserData.current_user_id
          };
          setCurrentUser(userData);
        }

        // Then fetch groups
        const filterParam = activeTab === 'my' ? '?filter=my' : '';
        const response = await fetch(`http://localhost:8080/api/groups${filterParam}`, {
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
          // Fetch membership status for each group, passing userData
          await fetchMembershipStatus(responseData.data, userData);
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

    async function fetchMembershipStatus(groupsList, user) {
      const memberships = {};
      const joinCounts = {};
      for (const group of groupsList) {
        try {
          const response = await fetch(`http://localhost:8080/api/groups/${group.id}/membership`, {
            credentials: 'include',
          });
          if (response.ok) {
            const data = await response.json();
            memberships[group.id] = data.is_member;
          }
        } catch (e) {
          console.error(`Failed to fetch membership for group ${group.id}:`, e);
          memberships[group.id] = false;
        }

        // Fetch pending join requests count for groups created by current user
        if (user && group.creator_id === user.id) {
          try {
            const requestsResponse = await fetch(`http://localhost:8080/api/groups/${group.id}/join-requests`, {
              credentials: 'include',
            });
            if (requestsResponse.ok) {
              const requestsData = await requestsResponse.json();
              joinCounts[group.id] = requestsData.data ? requestsData.data.length : 0;
            }
          } catch (e) {
            console.error(`Failed to fetch join requests for group ${group.id}:`, e);
            joinCounts[group.id] = 0;
          }
        }
      }
      setGroupMemberships(memberships);
      setPendingJoinCounts(joinCounts);
    }

    fetchData();
  }, [router, activeTab]);

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

      await response.json();
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

  const isGroupMember = (group) => {
    return groupMemberships[group.id] === true;
  };

  const canJoinGroup = (group) => {
    return currentUser && !isGroupMember(group) && !joinRequests[group.id];
  };

  if (loading) {
    return (
      <div className="flex min-h-screen justify-center items-center">
        <p className="text-xl">Loading groups...</p>
      </div>
    );
  }

  if (error && !loading) {
    return (
      <div className="flex min-h-screen justify-center items-center">
        <p className="text-xl">Error: {error}</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Groups</h1>
        <div className="flex gap-3">
          <Link href="/groups/invites">
            <button className="border border-black hover:bg-black hover:text-white font-bold py-2 px-4 transition duration-300 ease-in-out flex items-center gap-2">
              <BellIcon className="h-5 w-5" />
              Invitations
            </button>
          </Link>
          <Link href="/groups/join-requests">
            <button className="border border-black hover:bg-black hover:text-white font-bold py-2 px-4 transition duration-300 ease-in-out flex items-center gap-2">
              <UserGroupIcon className="h-5 w-5" />
              Manage Requests
            </button>
          </Link>
          <Link href="/groups/create">
            <button className="bg-black hover:bg-gray-800 text-white font-bold py-2 px-4 transition duration-300 ease-in-out">
              Create New Group
            </button>
          </Link>
        </div>
      </div>

      {/* Tab Navigation */}
      <div className="mb-6 border-b border-gray-400">
        <nav className="flex space-x-8" aria-label="Tabs">
          <button
            onClick={() => setActiveTab('all')}
            className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${activeTab === 'all'
              ? 'border-black text-black'
              : 'border-transparent text-gray-500 hover:text-black hover:border-gray-400'
              }`}
          >
            All Groups
          </button>
          <button
            onClick={() => setActiveTab('my')}
            className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${activeTab === 'my'
              ? 'border-black text-black'
              : 'border-transparent text-gray-500 hover:text-black hover:border-gray-400'
              }`}
          >
            My Groups
          </button>
        </nav>
      </div>

      {(!Array.isArray(groups) || groups.length === 0) ? (
        <p className="text-gray-600 text-center">No groups found. Start by creating a new one!</p>
      ) : (
        <div className="gap-6">
          {groups.map((group) => (
            <div key={group.id} className="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300 mb-4" >
              <Link href={`/groups/${group.id}`}>

                <div className="flex justify-between items-start mb-3">
                  <h2 className="text-xl font-semibold text-gray-900">{group.title}</h2>

                  <span className={`px-2 py-1 text-xs rounded-full ${group.privacy_setting === 'public' ? 'bg-green-100 text-green-800' :
                    group.privacy_setting === 'private' ? 'bg-yellow-100 text-yellow-800' :
                      'bg-red-100 text-red-800'
                    }`}>
                    {group.privacy_setting}
                  </span>
                </div>
              </Link>

              <div className="flex flex-col gap-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-500">
                    Created: {new Date(group.created_at).toLocaleDateString()}
                  </span>

                  {isGroupCreator(group) ? (
                    <div className="flex items-center gap-2">
                      <span className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
                        Your Group
                      </span>
                      {pendingJoinCounts[group.id] > 0 && (
                        <span className="bg-orange-100 text-orange-800 px-3 py-1 rounded-full text-sm font-medium flex items-center gap-1">
                          <BellIcon className="h-4 w-4" />
                          {pendingJoinCounts[group.id]} request{pendingJoinCounts[group.id] > 1 ? 's' : ''}
                        </span>
                      )}
                    </div>
                  ) : canJoinGroup(group) ? (
                    <button
                      onClick={() => handleJoinRequest(group.id)}
                      disabled={pendingRequests[group.id]}
                      className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all duration-200 ${pendingRequests[group.id]
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

                {isGroupMember(group) && (
                  <Link href={`/groups/${group.id}`}>
                    <button className="w-full bg-green-500 hover:bg-green-600 text-white font-medium py-2 px-4 rounded-lg transition-all duration-200 flex items-center justify-center gap-2">
                      <UserGroupIcon className="h-5 w-5" />
                      View Group
                    </button>
                  </Link>
                )}
              </div>

            </div>
          ))}
        </div>
      )}
    </div>
  );
}