"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import CreatePost from "@/components/createpost";
import PostCard from "@/components/postCard"; // Import PostCard
import { UserPlusIcon } from '@heroicons/react/24/outline';

export default function GroupPage() {
  const [group, setGroup] = useState(null);
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showCreatePost, setShowCreatePost] = useState(false);
  const [showInviteModal, setShowInviteModal] = useState(false);
  const [currentUser, setCurrentUser] = useState(null);
  const [isMember, setIsMember] = useState(false);
  const params = useParams();
  const router = useRouter();
  const { id } = params;

  useEffect(() => {
    if (id) {
      checkMembership();
    }
  }, [id]);

  const checkMembership = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/groups/${id}/membership`, {
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Failed to check membership");
      }
      const data = await response.json();
      if (data.is_member) {
        setIsMember(true);
        fetchCurrentUser();
        fetchGroupDetails();
        fetchGroupPosts();
      } else {
        setError("You are not a member of this group");
        setLoading(false);
      }
    } catch (e) {
      setError(e.message);
      setLoading(false);
    }
  };

  const fetchCurrentUser = async () => {
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
  };

  const fetchGroupDetails = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/groups/${id}`, {
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Failed to fetch group details");
      }
      const data = await response.json();
      setGroup(data);
    } catch (e) {
      setError(e.message);
    }
  };

  const fetchGroupPosts = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/groups/${id}/posts`, {
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Failed to fetch group posts");
      }
      const data = await response.json();
      setPosts(data || []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const isGroupCreator = currentUser && group && group.creator_id === currentUser.id;

  return (
    <div className="container mx-auto p-4">
      {loading && <p>Loading...</p>}
      {error && (
        <div className="text-center py-12">
          <p className="text-red-500 text-xl mb-4">{error}</p>
          <button
            onClick={() => router.push('/groups')}
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
          >
            Back to Groups
          </button>
        </div>
      )}
      {!error && group && (
        <div>
          <div className="flex justify-between items-start mb-4">
            <div>
              <h1 className="text-3xl font-bold">{group.title}</h1>
              <p className="text-gray-600">{group.description}</p>
              {isGroupCreator && (
                <span className="inline-block bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium mt-2">
                  You are the creator
                </span>
              )}
            </div>
          </div>

          <div className="flex gap-3">
            <button
              onClick={() => setShowCreatePost(true)}
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-300 ease-in-out"
            >
              Create Post
            </button>

            {isGroupCreator && (
              <button
                onClick={() => setShowInviteModal(true)}
                className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition duration-300 ease-in-out flex items-center gap-2"
              >
                <UserPlusIcon className="h-5 w-5" />
                Invite Members
              </button>
            )}
          </div>
        </div>
      )}
      {showCreatePost && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <CreatePost onClose={() => setShowCreatePost(false)} groupId={id} />
        </div>
      )}
      {showInviteModal && (
        <InviteMembersModal
          groupId={id}
          onClose={() => setShowInviteModal(false)}
        />
      )}
      <div className="mt-8">
        {posts.map((post) => (
          <PostCard
            key={post.id}
            post={post}
          />
        ))}
      </div>
    </div>
  );
}

function InviteMembersModal({ groupId, onClose }) {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [inviting, setInviting] = useState({});

  useEffect(() => {
    fetchAvailableUsers();
  }, []);

  const fetchAvailableUsers = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/users/available", {
        credentials: 'include',
      });
      if (response.ok) {
        const data = await response.json();
        setUsers(data || []);
      }
    } catch (e) {
      console.error("Failed to fetch users:", e);
    } finally {
      setLoading(false);
    }
  };

  const handleInvite = async (userId) => {
    try {
      setInviting(prev => ({ ...prev, [userId]: true }));

      const response = await fetch(`http://localhost:8080/api/groups/${groupId}/invite`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ user_id: userId }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to send invite');
      }

      alert('Invite sent successfully!');

    } catch (error) {
      console.error('Error sending invite:', error);
      alert(`Error: ${error.message}`);
    } finally {
      setInviting(prev => ({ ...prev, [userId]: false }));
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl p-6 max-w-md w-full max-h-[80vh] overflow-y-auto">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-bold">Invite Members</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-2xl"
          >
            &times;
          </button>
        </div>

        {loading ? (
          <p>Loading users...</p>
        ) : users.length === 0 ? (
          <p className="text-gray-600">No users available to invite.</p>
        ) : (
          <div className="space-y-3">
            {users.map((user) => (
              <div key={user.id} className="flex justify-between items-center p-3 bg-gray-50 rounded-lg">
                <div>
                  <p className="font-medium">{user.fname} {user.lname}</p>
                  <p className="text-sm text-gray-500">@{user.nickname}</p>
                </div>
                <button
                  onClick={() => handleInvite(user.id)}
                  disabled={inviting[user.id]}
                  className={`px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                    inviting[user.id]
                      ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                      : 'bg-blue-500 hover:bg-blue-600 text-white'
                  }`}
                >
                  {inviting[user.id] ? 'Inviting...' : 'Invite'}
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
