"use client";

import { useState, useEffect } from 'react';
import { UserIcon } from '@heroicons/react/24/solid';

export default function FollowersFollowingSection({ userId }) {
  const [activeTab, setActiveTab] = useState('followers');
  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Fetch followers data
  const fetchFollowers = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await fetch(`/api/users/${userId}/followers`, {
        credentials: 'include',
      });
      
      if (!response.ok) {
        if (response.status === 401) {
          throw new Error('You must be logged in to view this content');
        } else if (response.status === 403) {
          throw new Error('Not authorized to view followers');
        } else {
          throw new Error('Failed to load followers');
        }
      }
      
      const data = await response.json();
      setFollowers(data.data || []);
    } catch (err) {
      setError(err.message);
      setFollowers([]);
    } finally {
      setLoading(false);
    }
  };

  // Fetch following data
  const fetchFollowing = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await fetch(`/api/users/${userId}/following`, {
        credentials: 'include',
      });
      
      if (!response.ok) {
        if (response.status === 401) {
          throw new Error('You must be logged in to view this content');
        } else if (response.status === 403) {
          throw new Error('Not authorized to view following list');
        } else {
          throw new Error('Failed to load following list');
        }
      }
      
      const data = await response.json();
      setFollowing(data.data || []);
    } catch (err) {
      setError(err.message);
      setFollowing([]);
    } finally {
      setLoading(false);
    }
  };

  // Load data when component mounts or userId changes
  useEffect(() => {
    if (userId) {
      fetchFollowers();
      fetchFollowing();
    }
  }, [userId]);

  // Handle tab switching
  const handleTabChange = (tab) => {
    setActiveTab(tab);
    setError(''); // Clear any previous errors when switching tabs
  };

  // Render user list
  const renderUserList = (users, emptyMessage) => {
    if (loading) {
      return (
        <div className="flex justify-center items-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      );
    }

    if (error) {
      return (
        <div className="text-center py-8">
          <p className="text-red-600 mb-4">{error}</p>
          <button 
            onClick={activeTab === 'followers' ? fetchFollowers : fetchFollowing}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Try Again
          </button>
        </div>
      );
    }

    if (!users || users.length === 0) {
      return (
        <div className="text-center py-8 text-gray-600">
          <UserIcon className="h-12 w-12 mx-auto mb-4 text-gray-400" />
          <p className="text-lg font-medium">{emptyMessage}</p>
        </div>
      );
    }

    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {users.map((user) => (
          <div key={user.user_id} className="bg-white rounded-lg shadow-sm border p-4 hover:shadow-md transition-shadow">
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 rounded-full overflow-hidden bg-gray-200 flex-shrink-0">
                {user.img_url ? (
                  <img
                    src={user.img_url}
                    alt={`${user.first_name} ${user.last_name}`}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full bg-blue-500 flex items-center justify-center text-white font-semibold">
                    {user.first_name.charAt(0)}{user.last_name.charAt(0)}
                  </div>
                )}
              </div>
              <div className="flex-1 min-w-0">
                <h3 className="text-sm font-semibold text-gray-900 truncate">
                  {user.first_name} {user.last_name}
                </h3>
                {user.nickname && (
                  <p className="text-xs text-gray-500 truncate">@{user.nickname}</p>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    );
  };

  return (
    <div className="bg-gray-50 rounded-lg p-6 mt-6">
      <div className="mb-6">
        <h2 className="text-xl font-bold text-gray-900 mb-4">Connections</h2>
        
        {/* Tab Navigation */}
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            <button
              onClick={() => handleTabChange('followers')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'followers'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Followers ({followers.length})
            </button>
            <button
              onClick={() => handleTabChange('following')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'following'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Following ({following.length})
            </button>
          </nav>
        </div>
      </div>

      {/* Tab Content */}
      <div className="mt-6">
        {activeTab === 'followers' && renderUserList(
          followers, 
          "No followers yet"
        )}
        {activeTab === 'following' && renderUserList(
          following, 
          "Not following anyone yet"
        )}
      </div>
    </div>
  );
}
