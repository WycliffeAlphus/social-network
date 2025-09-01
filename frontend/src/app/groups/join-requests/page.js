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
      
      // Fetch groups where current user is creator
      const groupsResponse = await fetch("http://localhost:8080/api/groups", {
        credentials: 'include',
      });

      if (groupsResponse.status === 401) {
        router.push('/login');
        return;
      }

      if (!groupsResponse.ok) {
        throw new Error('Failed to fetch groups');
      }

      const groupsData = await groupsResponse.json();
      const userGroups = groupsData.data || [];
      
      // Get current user to filter groups they created
      const userResponse = await fetch("http://localhost:8080/api/profile/currentuser", {
        credentials: 'include',
      });
      
      if (userResponse.ok) {
        const currentUser = await userResponse.json();
        const createdGroups = userGroups.filter(group => group.creator_id === currentUser.id);
        setGroups(createdGroups);
        
        // For each group, fetch pending join requests
        const requestsData = {};
        for (const group of createdGroups) {
          try {
            const requestsResponse = await fetch(`http://localhost:8080/api/groups/${group.id}/join-requests`, {
              credentials: 'include',
            });

            if (requestsResponse.ok) {
              const requestsResult = await requestsResponse.json();
              requestsData[group.id] = requestsResult.data || [];
            } else {
              console.error(`Failed to fetch requests for group ${group.id}`);
              requestsData[group.id] = [];
            }
          } catch (err) {
            console.error(`Failed to fetch requests for group ${group.id}:`, err);
            requestsData[group.id] = [];
          }
        }
        setJoinRequests(requestsData);
      }

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
      <div className="flex min-h-screen justify-center items-center bg-gray-100">
        <p className="text-xl text-gray-700">Loading join requests...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-screen justify-center items-center bg-gray-100">
        <p className="text-xl text-red-500">Error: {error}</p>
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
          className="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded transition duration-300 ease-in-out"
        >
          Back to Groups
        </button>
      </div>

      {groups.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-600 text-lg">You haven't created any groups yet.</p>
          <button
            onClick={() => router.push('/groups/create')}
            className="mt-4 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-300 ease-in-out"
          >
            Create Your First Group
          </button>
        </div>
      ) : (
        <div className="space-y-6">
          {groups.map((group) => {
            const requests = joinRequests[group.id] || [];
            
            return (
              <div key={group.id} className="bg-white rounded-lg shadow-md p-6">
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <h2 className="text-xl font-semibold text-gray-900">{group.title}</h2>
                    <p className="text-gray-600">{group.description}</p>
                  </div>
                  <span className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
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
                        <div key={request.user_id} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                          <div className="flex items-center gap-3">
                            <UserIcon className="h-8 w-8 text-gray-400" />
                            <div>
                              <p className="font-medium text-gray-900">{request.user_name}</p>
                              <p className="text-sm text-gray-500">
                                Requested: {new Date(request.requested_at).toLocaleDateString()}
                              </p>
                            </div>
                          </div>
                          
                          <div className="flex gap-2">
                            <button
                              onClick={() => handleAcceptRequest(group.id, request.user_id, request.user_name)}
                              disabled={isProcessing}
                              className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                                isProcessing === 'accepting'
                                  ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                                  : 'bg-green-500 hover:bg-green-600 text-white hover:shadow-md'
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
                              className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                                isProcessing === 'rejecting'
                                  ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                                  : 'bg-red-500 hover:bg-red-600 text-white hover:shadow-md'
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
