"use client";

import FollowSuggestion from "@/components/followsuggestions";
import Loading from '@/components/loading'
import Rightbar from "@/components/rightbar";
import Link from "next/link";
import { useEffect, useState } from "react";

export default function FollowRequestsPage() {
  const [followRequests, setFollowRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchFollowRequests = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/follow-requests', {
          method: 'GET',
          credentials: 'include'
        });
        if (!response.ok) {
          throw new Error('Failed to fetch follow requests');
        }
        const data = await response.json();
        setFollowRequests(data.data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchFollowRequests();
  }, []);

  const handleAccept = async (followerId) => {
    try {
      const response = await fetch('http://localhost:8080/api/follow/accept', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ followerId })
      });

      if (response.ok) {
        // Remove the accepted request from the list
        setFollowRequests(prev => prev.filter(request => request.follower_id !== followerId));
      } else {
        const errorData = await response.json();
        console.error(`Failed to accept request: ${errorData.message || 'Unknown error'}`);
      }
    } catch (err) {
      console.error('Error accepting follow request:', err);
    }
  };

  const handleDecline = async (followerId) => {
    try {
      const response = await fetch('http://localhost:8080/api/follow/decline', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ followerId })
      });

      if (response.ok) {
        // Remove the declined request from the list
        setFollowRequests(prev => prev.filter(request => request.follower_id !== followerId));
      } else {
        const errorData = await response.json();
        console.error(`Failed to decline request: ${errorData.message || 'Unknown error'}`);
      }
    } catch (err) {
      console.error('Error declining follow request:', err);
    }
  };

  if (loading) return <div className="flex items-center justify-center h-screen"><Loading /></div>

  return (
    <div className="flex min-h-screen">
      <main className="flex-1 border-x mr-[20px] border-gray-400">
        <div className="lg:hidden">
          <FollowSuggestion />
        </div>
        <div className="p-4 border-t lg:border-0 border-gray-400">
          <h1 className="text-xl font-bold mb-4">Follow Requests</h1>

          {error ? (
            <div className="text-red-500 p-4">{error}</div>
          ) : followRequests.length === 0 ? (
            <div className="text-gray-500 p-4">No pending follow requests</div>
          ) : (
            <div className="space-y-4">
              {followRequests.map((request) => (
                <div key={request.follower_id} className="flex items-center justify-between p-3">
                  <Link
                    href={`/profile/${request.follower_id}`}
                    className="group flex items-center gap-3 p-3 w-fit"
                  >
                    {request.follower_avatar ? (
                      <img
                        src={request.follower_avatar}
                        alt={`${request.follower_fname} ${request.follower_lname}`}
                        className="w-12 h-12 rounded-full object-cover"
                      />
                    ) : (
                      <div className="w-12 h-12 rounded-full bg-gray-200 overflow-hidden flex items-center justify-center">
                        <span className="text-lg text-black font-medium">
                          {request.follower_fname?.charAt(0).toUpperCase()}
                          {request.follower_lname?.charAt(0).toUpperCase()}
                        </span>
                      </div>
                    )}
                    <div>
                      <span className="font-medium group-hover:text-[#4169e1]">
                        {request.follower_fname} {request.follower_lname}
                      </span>
                    </div>
                  </Link>
                  <div className="flex space-x-2">
                    <button
                      onClick={() => handleAccept(request.follower_id)}
                      className="px-4 py-1 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition"
                    >
                      Accept
                    </button>
                    <button
                      onClick={() => handleDecline(request.follower_id)}
                      className="px-4 py-1 bg-gray-200 text-gray-800 rounded-full hover:bg-gray-300 transition"
                    >
                      Decline
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>
      <Rightbar />
    </div>
  );
}
