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
    const [isPublic, setIsPublic] = useState(false)
    const [profileVisibility, setProfileVisibility] = useState('private')
    const [followersCount, setFollowersCount] = useState(0)
    const [followingCount, setFollowingCount] = useState(0)
    const [followStatus, setFollowStatus] = useState('not_following')

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const response = await fetch(`http://localhost:8080/api/profile/${id}`, {
                    credentials: 'include',
                })

                if (!response.ok) {
                    throw new Error('Failed to fetch profile')
                }

                const data = await response.json()
                setProfileData(data.profile)
                setIsOwner(data.current_user_id === data.profile.id)
                setProfileVisibility(data.profile.profile_visibility)
                setIsPublic(data.profile.profile_visibility === 'public')
                setFollowsMe(data.follows_me)
                setFollowersCount(data.followers_count)
                setFollowingCount(data.following_count)

                // Fetch follow status if not the owner
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
            const response = await fetch(`http://localhost:8080/api/follow-status/${id}`, {
                credentials: 'include',
            })
            if (response.ok) {
                const data = await response.json()
                setFollowStatus(data.status)
            } else {
                console.error('Failed to fetch follow status:', response.status, response.statusText)
                return
            }
        } catch (err) {
            console.error('Error fetching follow status:', err)
        }
    }

    const handleFollow = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/users/follow', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ userId: id }),
                credentials: 'include'
            });
            if (response.ok) {
                const data = await response.json();
                setFollowStatus(data.status);
            } else {
                console.error('Failed to follow user');
            }
        } catch (err) {
            console.error('Error following user:', err);
        }
    };

    const handleCancelRequest = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/follow/cancel', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ userId: id }),
                credentials: 'include'
            });

            if (response.ok) {
                setFollowStatus('not_following')
            } else {
                console.error('Failed to cancel follow request');
            }
        } catch (err) {
            console.error('Error cancelling follow request:', err);
        }
    };

    const handleToggleVisibility = async () => {
        const newVisibility = profileVisibility === 'public' ? 'private' : 'public'

        try {
            const response = await fetch('http://localhost:8080/api/profile/update', {
                method: 'PUT',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ "profile_visibility": newVisibility }),
            })

            if (response.ok) {
                setProfileVisibility(newVisibility)
                setIsPublic(newVisibility === 'public')
            } else {
                console.error('Failed to update profile visibility')
            }
        } catch (err) {
            console.error('Error updating profile visibility:', err)
        }
    };

    if (loading) {
        return <div className="flex justify-center items-center h-screen">Loading...</div>
    }

    if (error) {
        return <div className="flex justify-center items-center h-screen text-red-500">{error}</div>
    }

    const shouldShowProfile = isOwner || isPublic || followsMe

    if (!shouldShowProfile) {
        let buttonLabel = "Follow";
        let buttonClass = "bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded";
        let disabled = false;

        if (followStatus === 'accepted') {
            buttonLabel = "Following";
            buttonClass = "border border-gray-400 hover:bg-gray-400 hover:text-black text-gray-800 px-4 py-2 rounded";
            disabled = true;
        } else if (followStatus === 'requested') {
            buttonLabel = "Requested";
            buttonClass = "bg-gray-300 text-gray-600 hover:bg-gray-400 hover:text-black px-4 py-2 rounded";
        }

        return (
            <div className="flex flex-col items-center justify-center h-screen">
                <h1 className="text-2xl font-bold mb-4">Private Profile</h1>
                <p className="text-gray-600 mb-6">{profileData.first_name} {profileData.last_name}'s profile is private. Follow them to view their profile.</p>
                {followStatus === 'requested' && (
                    <p className="text-sm text-gray-500 mb-4">This account is private</p>
                )}
                <button
                    onClick={() => {
                        if (followStatus === 'requested') {
                            handleCancelRequest()
                        } else {
                            handleFollow()
                        }
                    }}
                    className={buttonClass}
                    disabled={disabled}
                >
                    {buttonLabel}
                </button>
            </div>
        )
    }

    return (
        <div className="max-w-4xl mx-auto py-10 px-4">
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
        </div>
    )
}
