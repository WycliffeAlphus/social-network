'use client';

import { useState, useEffect } from 'react';
import { UserIcon } from '@heroicons/react/24/solid';

export default function FollowersContent({ userId = 'me' }) {
  const [followers, setFollowers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchFollowers();
  }, [userId]);

  const fetchFollowers = async () => {
    try {
      setLoading(true);
      const response = await fetch(`/api/users/${userId}/followers`, {
        credentials: 'include',
      });

      if (response.ok) {
        const data = await response.json();
        setFollowers(data.data || []);
      } else if (response.status === 403) {
        setError('You do not have permission to view this user\'s followers.');
      } else {
        setError('Failed to load followers.');
      }
    } catch (err) {
      setError('An error occurred while loading followers.');
      console.error('Error fetching followers:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleUnfollow = async (followerId) => {
    try {
      const response = await fetch(`/api/users/${followerId}/follow`, {
        method: 'DELETE',
        credentials: 'include',
      });

      if (response.ok) {
        // Remove the unfollowed user from the list if viewing own followers
        if (userId === 'me') {
          setFollowers(followers.filter(f => f.user_id !== followerId));
        }
        alert('Successfully unfollowed user');
      } else {
        const errorData = await response.text();
        alert(`Failed to unfollow: ${errorData}`);
      }
    } catch (err) {
      alert('An error occurred while unfollowing');
      console.error('Error unfollowing user:', err);
    }
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-semibold mb-4 text-gray-800">Followers</h2>
        <div className="flex justify-center items-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <span className="ml-2 text-gray-600">Loading followers...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-semibold mb-4 text-gray-800">Followers</h2>
        <div className="text-center py-8">
          <p className="text-red-600 text-lg">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-semibold mb-4 text-gray-800">
        Followers ({followers.length})
      </h2>

      {followers.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {followers.map((follower) => (
            <div key={follower.user_id} className="bg-gray-50 rounded-lg shadow-sm p-5 flex flex-col items-center text-center">
              {/* Avatar Section */}
              <div className="w-24 h-24 mb-4 rounded-full overflow-hidden border-2 border-blue-500 flex items-center justify-center bg-gray-200">
                {follower.img_url ? (
                  <img
                    src={follower.img_url}
                    alt={`${follower.first_name}'s avatar`}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <UserIcon className="w-12 h-12 text-gray-400" />
                )}
              </div>

              {/* User Info */}
              <h3 className="text-xl font-semibold text-gray-900">
                {follower.first_name} {follower.last_name}
              </h3>
              {follower.nickname && (
                <p className="text-gray-500 mb-2">@{follower.nickname}</p>
              )}
              <p className="text-gray-600 text-sm mb-4">{follower.email}</p>

              {/* Profile Visibility Badge */}
              <div className="mb-4">
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                  follower.profile_visibility === 'public' 
                    ? 'bg-green-100 text-green-800' 
                    : 'bg-yellow-100 text-yellow-800'
                }`}>
                  {follower.profile_visibility === 'public' ? 'Public' : 'Private'}
                </span>
              </div>

              {/* Action Buttons */}
              <div className="flex space-x-2">
                <button
                  onClick={() => window.location.href = `/profile/${follower.user_id}`}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition duration-200 text-sm"
                >
                  View Profile
                </button>
                {userId === 'me' && (
                  <button
                    onClick={() => handleUnfollow(follower.user_id)}
                    className="px-4 py-2 bg-gray-300 text-gray-800 rounded-md hover:bg-gray-400 transition duration-200 text-sm"
                  >
                    Remove
                  </button>
                )}
              </div>

              {/* Follow Date */}
              <p className="text-xs text-gray-400 mt-2">
                Following since {new Date(follower.created_at).toLocaleDateString()}
              </p>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <UserIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-center text-gray-600 text-lg">
            {userId === 'me' ? 'You have no followers yet.' : 'This user has no followers.'}
          </p>
          <p className="text-center text-gray-500 text-sm mt-2">
            {userId === 'me' ? 'Share your profile to gain followers!' : ''}
          </p>
        </div>
      )}
    </div>
  );
}
