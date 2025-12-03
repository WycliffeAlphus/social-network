"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { CheckIcon, XMarkIcon } from '@heroicons/react/24/outline';

export default function GroupInvitesPage() {
  const [invites, setInvites] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [processing, setProcessing] = useState({});
  const router = useRouter();

  useEffect(() => {
    fetchInvites();
  }, []);

  const fetchInvites = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/groups/invites", {
        credentials: 'include',
      });

      if (response.status === 401) {
        router.push('/login');
        return;
      }

      if (!response.ok) {
        throw new Error('Failed to fetch invitations');
      }

      const responseData = await response.json();
      setInvites(responseData.data || []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const handleAccept = async (groupId) => {
    try {
      setProcessing(prev => ({ ...prev, [groupId]: 'accepting' }));

      const response = await fetch(`http://localhost:8080/api/groups/${groupId}/invites/accept`, {
        method: 'POST',
        credentials: 'include',
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to accept invitation');
      }

      // Remove from list after accepting
      setInvites(prev => prev.filter(invite => invite.id !== groupId));
      alert('Invitation accepted! You are now a member of the group.');
    } catch (error) {
      console.error('Error accepting invitation:', error);
      alert(`Error: ${error.message}`);
    } finally {
      setProcessing(prev => ({ ...prev, [groupId]: null }));
    }
  };

  const handleReject = async (groupId) => {
    try {
      setProcessing(prev => ({ ...prev, [groupId]: 'rejecting' }));

      const response = await fetch(`http://localhost:8080/api/groups/${groupId}/invites/reject`, {
        method: 'POST',
        credentials: 'include',
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to reject invitation');
      }

      // Remove from list after rejecting
      setInvites(prev => prev.filter(invite => invite.id !== groupId));
    } catch (error) {
      console.error('Error rejecting invitation:', error);
      alert(`Error: ${error.message}`);
    } finally {
      setProcessing(prev => ({ ...prev, [groupId]: null }));
    }
  };

  if (loading) {
    return (
      <div className="flex min-h-screen justify-center items-center">
        <p className="text-xl">Loading invitations...</p>
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

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Group Invitations</h1>
        <Link href="/groups">
          <button className="border border-black hover:bg-black hover:text-white font-bold py-2 px-4 transition duration-300 ease-in-out">
            Back to Groups
          </button>
        </Link>
      </div>

      {invites.length === 0 ? (
        <p className="text-center text-gray-600">No pending invitations.</p>
      ) : (
        <div className="space-y-4">
          {invites.map((invite) => (
            <div key={invite.id} className="border border-gray-400 p-6 hover:shadow-lg transition-shadow duration-300">
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <h2 className="text-xl font-semibold mb-2">{invite.title}</h2>
                  <p className="text-gray-600 mb-3">{invite.description}</p>
                  <div className="flex items-center gap-4 text-sm text-gray-500">
                    <span className="border border-black px-2 py-1">{invite.privacy_setting}</span>
                    <span>Created: {new Date(invite.created_at).toLocaleDateString()}</span>
                  </div>
                </div>

                <div className="flex gap-2 ml-4">
                  <button
                    onClick={() => handleAccept(invite.id)}
                    disabled={processing[invite.id]}
                    className={`flex items-center gap-2 px-4 py-2 font-medium transition-all duration-200 ${
                      processing[invite.id] === 'accepting'
                        ? 'border border-gray-400 text-gray-400 cursor-not-allowed'
                        : 'bg-black hover:bg-gray-800 text-white'
                    }`}
                  >
                    <CheckIcon className="h-5 w-5" />
                    {processing[invite.id] === 'accepting' ? 'Accepting...' : 'Accept'}
                  </button>

                  <button
                    onClick={() => handleReject(invite.id)}
                    disabled={processing[invite.id]}
                    className={`flex items-center gap-2 px-4 py-2 font-medium transition-all duration-200 ${
                      processing[invite.id] === 'rejecting'
                        ? 'border border-gray-400 text-gray-400 cursor-not-allowed'
                        : 'border border-black hover:bg-black hover:text-white'
                    }`}
                  >
                    <XMarkIcon className="h-5 w-5" />
                    {processing[invite.id] === 'rejecting' ? 'Rejecting...' : 'Reject'}
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
