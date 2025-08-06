'use client'

import { use, useEffect, useState } from 'react'
import ProfileHeader from '@/components/ProfileHeader'
import ProfileDetails from '@/components/ProfileDetails'

export default function Profile({ params }) {
  const { id } = use(params)

  const [profileData, setProfileData] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [isOwner, setIsOwner] = useState(false)
  const [followsMe, setFollowsMe] = useState(false)
  const [profileVisibility, setProfileVisibility] = useState('private')
  const [followersCount, setFollowersCount] = useState(0)
  const [followingCount, setFollowingCount] = useState(0)
  const [followStatus, setFollowStatus] = useState('not_following')

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await fetch(`http://localhost:8080/api/profile/${id}`, {
          credentials: 'include',
        })
        const data = await res.json()
        setProfileData(data.profile)
        setIsOwner(data.current_user_id === data.profile.id)
        setProfileVisibility(data.profile.profile_visibility)
        setFollowsMe(data.follows_me)
        setFollowersCount(data.followers_count)
        setFollowingCount(data.following_count)

        if (data.current_user_id !== data.profile.id) {
          await fetchFollowStatus()
        }
      } catch (err) {
        setError(err.message)
      } finally {
        setLoading(false)
      }
    }

    fetchProfile()
  }, [id])

  const fetchFollowStatus = async () => {
    try {
      const res = await fetch(`http://localhost:8080/api/follow-status/${id}`, {
        credentials: 'include',
      })
      if (!res.ok) {
        console.error('Failed to fetch follow status:', res.status, res.statusText)
        return
      }
      const data = await res.json()
      setFollowStatus(data.status)
    } catch (err) {
      console.error('Failed to fetch follow status:', err)
    }
  }

  const handleFollow = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/users/follow', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ userId: id }),
        credentials: 'include',
      });
      if (res.ok) {
        setFollowStatus('requested');
      } else {
        console.error('Failed to follow user');
      }
    } catch (error) {
      console.error('Error following user:', error);
    }
  };

  const handleCancelRequest = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/follow/cancel', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ userId: id }),
        credentials: 'include',
      });
      if (res.ok) {
        setFollowStatus('not_following');
      } else {
        console.error('Failed to cancel follow request');
      }
    } catch (error) {
      console.error('Error cancelling follow request:', error);
    }
  };

  const handleToggleVisibility = async () => {
    const newVisibility = profileVisibility === 'public' ? 'private' : 'public';
    try {
      const res = await fetch('http://localhost:8080/api/profile/update', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ profile_visibility: newVisibility }),
      });
      if (res.ok) {
        setProfileVisibility(newVisibility);
      } else {
        console.error('Failed to update profile visibility');
      }
    } catch (error) {
      console.error('Error updating profile visibility:', error);
    }
  };

  if (loading) return <div className="text-center py-10">Loading...</div>
  if (error) return <div className="text-center py-10 text-red-500">{error}</div>

  const showProfile = isOwner || profileVisibility === 'public' || followsMe

  return (
    <div className="max-w-4xl mx-auto py-10 px-4">
      {showProfile ? (
        <>
          <ProfileHeader
            user={profileData}
            isOwner={isOwner}
            isPublic={profileVisibility === 'public'}
            followersCount={followersCount}
            followingCount={followingCount}
            followStatus={followStatus}
            onFollow={handleFollow}
            onCancelRequest={handleCancelRequest}
            onToggleVisibility={handleToggleVisibility}
          />
          <ProfileDetails user={profileData} />
        </>
      ) : (
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-2">Private Profile</h1>
          <p className="mb-4">{profileData.first_name} {profileData.last_name}'s profile is private.</p>
          <button
            onClick={followStatus === 'requested' ? handleCancelRequest : handleFollow}
            className="bg-blue-500 text-white px-4 py-2 rounded"
          >
            {followStatus === 'requested' ? 'Cancel Request' : 'Follow'}
          </button>
        </div>
      )}
    </div>
  )
}
