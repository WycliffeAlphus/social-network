'use client';

import { useState, useEffect } from 'react';
import { UserIcon } from '@heroicons/react/24/solid';

export default function FollowingContent({ userId = 'me' }) {
  const [following, setFollowing] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchFollowing();
  }, [userId]);

  const fetchFollowing = async () => {
    try {
      setLoading(true);
      const response = await fetch(`/api/users/${userId}/following`, {
        credentials: 'include',
      });

      if (response.ok) {
        const data = await response.json();
        setFollowing(data.data || []);
      } else if (response.status === 401) {
        setError('Please log in to view who you are following.');
      } else if (response.status === 403) {
        setError('You do not have permission to view this user\'s following list.');
      } else {
        setError('Failed to load following list.');
      }
    } catch (err) {
      setError('An error occurred while loading following list.');
      console.error('Error fetching following:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleUnfollow = async (followedUserId) => {
    try {
      const response = await fetch(`/api/users/${followedUserId}/follow`, {
        method: 'DELETE',
        credentials: 'include',
      });

      if (response.ok) {
        // Remove the unfollowed user from the list
        setFollowing(following.filter(f => f.user_id !== followedUserId));
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

  const handleFollow = async (followUserId) => {
    try {
      const response = await fetch(`/api/users/${followUserId}/follow`, {
        method: 'POST',
        credentials: 'include',
      });

      if (response.ok) {
        const result = await response.json();
        alert(result.message || 'Successfully followed user');
        // Optionally refresh the list
        fetchFollowing();
      } else {
        const errorData = await response.text();
        alert(`Failed to follow: ${errorData}`);
      }
    } catch (err) {
      alert('An error occurred while following');
      console.error('Error following user:', err);
    }
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-semibold mb-4 text-gray-800">Following</h2>
        <div className="flex justify-center items-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <span className="ml-2 text-gray-600">Loading following list...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-semibold mb-4 text-gray-800">Following</h2>
        <div className="text-center py-8">
          <p className="text-red-600 text-lg">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-semibold mb-4 text-gray-800">
        Following ({following.length})
      </h2>

      {following.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {following.map((followedUser) => (
            <div key={followedUser.user_id} className="bg-gray-50 rounded-lg shadow-sm p-5 flex flex-col items-center text-center">
              {/* Avatar Section */}
              <div className="w-24 h-24 mb-4 rounded-full overflow-hidden border-2 border-blue-500 flex items-center justify-center bg-gray-200">
                {followedUser.img_url ? (
                  <img
                    src={followedUser.img_url}
                    alt={`${followedUser.first_name}'s avatar`}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <UserIcon className="w-12 h-12 text-gray-400" />
                )}
              </div>

              {/* User Info */}
              <h3 className="text-xl font-semibold text-gray-900">
                {followedUser.first_name} {followedUser.last_name}
              </h3>
              {followedUser.nickname && (
                <p className="text-gray-500 mb-2">@{followedUser.nickname}</p>
              )}
              <p className="text-gray-600 text-sm mb-4">{followedUser.email}</p>

              {/* Profile Visibility Badge */}
              <div className="mb-4">
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                  followedUser.profile_visibility === 'public' 
                    ? 'bg-green-100 text-green-800' 
                    : 'bg-yellow-100 text-yellow-800'
                }`}>
                  {followedUser.profile_visibility === 'public' ? 'Public' : 'Private'}
                </span>
              </div>

              {/* Action Buttons */}
              <div className="flex space-x-2">
                <button
                  onClick={() => window.location.href = `/profile/${followedUser.user_id}`}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition duration-200 text-sm"
                >
                  View Profile
                </button>
                {userId === 'me' && (
                  <button
                    onClick={() => handleUnfollow(followedUser.user_id)}
                    className="px-4 py-2 bg-red-500 text-white rounded-md hover:bg-red-600 transition duration-200 text-sm"
                  >
                    Unfollow
                  </button>
                )}
              </div>

              {/* Follow Date */}
              <p className="text-xs text-gray-400 mt-2">
                Following since {new Date(followedUser.created_at).toLocaleDateString()}
              </p>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <UserIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-center text-gray-600 text-lg">
            {userId === 'me' ? 'Not following anyone currently' : 'This user is not following anyone'}
          </p>
          <p className="text-center text-gray-500 text-sm mt-2">
            {userId === 'me' ? 'People you follow will appear here' : ''}
          </p>
        </div>
      )}
    </div>
  );
}
