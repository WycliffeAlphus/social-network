"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { CheckIcon, XMarkIcon, UserIcon, ClockIcon } from '@heroicons/react/24/outline';

export default function JoinRequestsPage() {
  const [groups, setGroups] = useState([]);
  const [joinRequests, setJoinRequests] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [processingRequests, setProcessingRequests] = useState({});
  const router = useRouter();

  useEffect(() => {
    fetchGroupsAndRequests();
  }, []);

  const fetchGroupsAndRequests = async () => {
    try {
      setLoading(true);

      // First, get current user
      const userResponse = await fetch("http://localhost:8080/api/profile/currentuser", {
        credentials: 'include',
      });

      if (userResponse.status === 401) {
        router.push('/login');
        return;
      }

      if (!userResponse.ok) {
        throw new Error('Failed to fetch current user');
      }

      const currentUser = await userResponse.json();
      console.log('Current user data:', currentUser);

      // Extract user ID from current_user_id field
      const userId = currentUser.current_user_id;
      console.log('Extracted user ID:', userId);

      // Then fetch all groups
      const groupsResponse = await fetch("http://localhost:8080/api/groups", {
        credentials: 'include',
      });

      if (!groupsResponse.ok) {
        throw new Error('Failed to fetch groups');
      }

      const groupsData = await groupsResponse.json();
      const userGroups = groupsData.data || [];
      console.log('All groups:', userGroups);

      // Filter groups created by current user
      const createdGroups = userGroups.filter(group => {
        console.log(`Group ${group.id}: creator_id=${group.creator_id}, userId=${userId}, matches=${group.creator_id === userId}`);
        return group.creator_id === userId;
      });
      console.log('Created groups after filter:', createdGroups);
      setGroups(createdGroups);

      // For each group, fetch pending join requests
      const requestsData = {};
      for (const group of createdGroups) {
        try {
          console.log(`Fetching join requests for group ${group.id}...`);
          const requestsResponse = await fetch(`http://localhost:8080/api/groups/${group.id}/join-requests`, {
            credentials: 'include',
          });

          console.log(`Response status for group ${group.id}:`, requestsResponse.status);

          if (requestsResponse.ok) {
            const requestsResult = await requestsResponse.json();
            console.log(`Requests result for group ${group.id}:`, requestsResult);
            requestsData[group.id] = requestsResult.data || [];
            console.log(`Stored ${requestsData[group.id].length} requests for group ${group.id}`);
          } else {
            console.error(`Failed to fetch requests for group ${group.id} - Status: ${requestsResponse.status}`);
            try {
              const errorData = await requestsResponse.json();
              console.error(`Error details:`, errorData);
            } catch (e) {
              console.error(`Could not parse error response`);
            }
            requestsData[group.id] = [];
          }
        } catch (err) {
          console.error(`Failed to fetch requests for group ${group.id}:`, err);
          requestsData[group.id] = [];
        }
      }
      console.log('Final requestsData:', requestsData);
      setJoinRequests(requestsData);

    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleAcceptRequest = async (groupId, userId, userName) => {
    const requestKey = `${groupId}-${userId}`;
    
    try {
      setProcessingRequests(prev => ({ ...prev, [requestKey]: 'accepting' }));
      
      const response = await fetch(`http://localhost:8080/api/groups/${groupId}/join?action=accept`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ user_id: userId }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to accept join request');
      }

      // Remove the request from the list
      setJoinRequests(prev => ({
        ...prev,
        [groupId]: prev[groupId].filter(req => req.user_id !== userId)
      }));

      alert(`Successfully accepted ${userName}'s join request!`);
      
    } catch (error) {
      console.error('Error accepting join request:', error);
      alert(`Error: ${error.message}`);
    } finally {
      setProcessingRequests(prev => ({ ...prev, [requestKey]: null }));
    }
  };

  const handleRejectRequest = async (groupId, userId, userName) => {
    const requestKey = `${groupId}-${userId}`;

    try {
      setProcessingRequests(prev => ({ ...prev, [requestKey]: 'rejecting' }));

      const response = await fetch(`http://localhost:8080/api/groups/${groupId}/join?action=reject`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ user_id: userId }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to reject join request');
      }

      // Remove the request from the list
      setJoinRequests(prev => ({
        ...prev,
        [groupId]: prev[groupId].filter(req => req.user_id !== userId)
      }));

      alert(`Successfully rejected ${userName}'s join request.`);

    } catch (error) {
      console.error('Error rejecting join request:', error);
      alert(`Error: ${error.message}`);
    } finally {
      setProcessingRequests(prev => ({ ...prev, [requestKey]: null }));
    }
  };

  if (loading) {
    return (
      <div className="flex min-h-screen justify-center items-center">
        <p className="text-xl">Loading join requests...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-screen justify-center items-center">
        <p className="text-xl">Error: {error}</p>
      </div>
    );
  }

  const totalPendingRequests = Object.values(joinRequests).reduce(
    (total, requests) => total + requests.length, 0
  );

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold">Join Requests</h1>
          <p className="text-gray-600 mt-1">
            Manage join requests for your groups ({totalPendingRequests} pending)
          </p>
        </div>
        <button
          onClick={() => router.back()}
          className="border border-black hover:bg-black hover:text-white font-bold py-2 px-4 transition duration-300 ease-in-out"
        >
          Back to Groups
        </button>
      </div>

      {groups.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-600 text-lg">You haven't created any groups yet.</p>
          <button
            onClick={() => router.push('/groups/create')}
            className="mt-4 bg-black hover:bg-gray-800 text-white font-bold py-2 px-4 transition duration-300 ease-in-out"
          >
            Create Your First Group
          </button>
        </div>
      ) : (
        <div className="space-y-6">
          {groups.map((group) => {
            const requests = joinRequests[group.id] || [];

            return (
              <div key={group.id} className="border border-gray-400 p-6 hover:shadow-lg transition-shadow">
                <div className="flex justify-between items-start mb-4">
                  <h2 className="text-xl font-semibold">{group.title}</h2>
                  <span className="bg-black text-white px-3 py-1 text-sm font-medium">
                    {requests.length} pending
                  </span>
                </div>

                {requests.length === 0 ? (
                  <p className="text-gray-500 italic">No pending join requests for this group.</p>
                ) : (
                  <div className="space-y-3">
                    {requests.map((request) => {
                      const requestKey = `${group.id}-${request.user_id}`;
                      const isProcessing = processingRequests[requestKey];

                      return (
                        <div key={request.user_id} className="flex items-center justify-between p-4 border border-gray-300 hover:border-gray-400">
                          <div className="flex items-center gap-3">
                            <UserIcon className="h-8 w-8 text-gray-400" />
                            <div>
                              <p className="font-medium">{request.user_name}</p>
                              <p className="text-sm text-gray-500">
                                Requested: {new Date(request.requested_at).toLocaleDateString()}
                              </p>
                            </div>
                          </div>

                          <div className="flex gap-2">
                            <button
                              onClick={() => handleAcceptRequest(group.id, request.user_id, request.user_name)}
                              disabled={isProcessing}
                              className={`flex items-center gap-2 px-4 py-2 font-medium transition-all duration-200 ${
                                isProcessing === 'accepting'
                                  ? 'border border-gray-400 text-gray-400 cursor-not-allowed'
                                  : 'bg-black hover:bg-gray-800 text-white'
                              }`}
                            >
                              {isProcessing === 'accepting' ? (
                                <>
                                  <ClockIcon className="h-4 w-4 animate-spin" />
                                  Accepting...
                                </>
                              ) : (
                                <>
                                  <CheckIcon className="h-4 w-4" />
                                  Accept
                                </>
                              )}
                            </button>

                            <button
                              onClick={() => handleRejectRequest(group.id, request.user_id, request.user_name)}
                              disabled={isProcessing}
                              className={`flex items-center gap-2 px-4 py-2 font-medium transition-all duration-200 ${
                                isProcessing === 'rejecting'
                                  ? 'border border-gray-400 text-gray-400 cursor-not-allowed'
                                  : 'border border-black hover:bg-black hover:text-white'
                              }`}
                            >
                              {isProcessing === 'rejecting' ? (
                                <>
                                  <ClockIcon className="h-4 w-4 animate-spin" />
                                  Rejecting...
                                </>
                              ) : (
                                <>
                                  <XMarkIcon className="h-4 w-4" />
                                  Reject
                                </>
                              )}
                            </button>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
