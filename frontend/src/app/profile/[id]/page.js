'use client'

import { useState, useEffect } from 'react'
import { use } from 'react'
import Link from 'next/link'

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
                    method: 'GET',
                    credentials: 'include'
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
                credentials: 'include'
            })
            if (response.ok) {
                const data = await response.json()
                setFollowStatus(data.status)
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
            })

            if (response.ok) {
                setFollowStatus('pending')
            }
        } catch (err) {
            console.error('Error following user:', err)
        }
    }

    const handleToggleVisibility = async () => {
        const newVisibility = profileVisibility === 'public' ? 'private' : 'public'

        try {
            const response = await fetch(`http://localhost:8080/api/profile/update`, {
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
                throw new Error('Failed to update visibility')
            }
        } catch (err) {
            console.error('Error updating profile visibility:', err)
        }
    }

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
        let isDisabled = false;

        if (followStatus === 'accepted') {
            buttonLabel = "Following";
            buttonClass = "border border-gray-400 hover:bg-gray-400 hover:text-black text-gray-800 px-4 py-2 rounded";
        } else if (followStatus === 'requested') {
            buttonLabel = "Requested";
            buttonClass = "bg-gray-300 text-gray-600 px-4 py-2 rounded cursor-not-allowed";
            isDisabled = true;
        }

        return (
            <div className="flex flex-col items-center justify-center h-screen">
                <h1 className="text-2xl font-bold mb-4">Private Profile</h1>
                <p className="text-gray-600 mb-6">{profileData.first_name} {profileData.last_name}'s profile is private. Follow them to view their profile.</p>
                {followStatus === 'requested' && (
                    <p className="text-sm text-gray-500 mb-4">This account is private</p>
                )}
                <button 
                    onClick={() => !isDisabled && handleFollow()}
                    disabled={isDisabled}
                    className={buttonClass}
                >
                    {buttonLabel}
                </button>
            </div>
        )
    }

    return (
        <div className="max-w-4xl mx-auto py-8 px-4">
            {/* Profile Header */}
            <div className="flex flex-col md:flex-row items-start md:items-center gap-6 mb-8">
                <div className="w-32 h-32 rounded-full overflow-hidden bg-gray-200">
                    {profileData.img_url ? (
                        <img
                            src={profileData.img_url}
                            alt={`${profileData.first_name}'s profile`}
                            className="w-full h-full object-cover"
                        />
                    ) : (
                        <div className="w-full h-full flex items-center justify-center text-4xl text-gray-400">
                            {profileData.first_name.charAt(0)}{profileData.last_name.charAt(0)}
                        </div>
                    )}
                </div>

                <div className="flex-1">
                    <div className="flex items-center gap-4 mb-2">
                        <h1 className="text-3xl font-bold">
                            {profileData.first_name} {profileData.last_name}
                        </h1>
                        {profileData.nickname && (
                            <span className="text-gray-500">({profileData.nickname})</span>
                        )}
                    </div>

                    <p className="text-gray-600 mb-4">{profileData.email}</p>

                    <div className="flex gap-6 mb-4">
                        <Link
                            href={`/followers/${id}`}
                            className="hover:text-blue-500 underline"
                        >
                            Followers <span className="font-semibold">{followersCount}</span>
                        </Link>
                        <Link
                            href={`/following/${id}`}
                            className="hover:text-blue-500 underline"
                        >
                            Following <span className="font-semibold">{followingCount}</span>
                        </Link>
                    </div>

                    {isOwner && (
                        <div className="flex items-center gap-4 mb-4">
                            <div className="flex items-center gap-2">
                                <span className="text-sm text-gray-600">
                                    Set account to public
                                </span>
                                <label className="relative inline-flex items-center cursor-pointer">
                                    <input
                                        type="checkbox"
                                        className="sr-only peer"
                                        checked={profileVisibility=== 'public'}
                                        onChange={handleToggleVisibility}
                                    />
                                    <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full
                                    peer peer-checked:after:translate-x-full peer-checked:after:border-white
                                    after:content-[''] after:absolute after:top-[2px] after:left-[2px]
                                    after:bg-white after:border-gray-300 after:border after:rounded-full
                                    after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-500"
                                    >
                                    </div>
                                </label>
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {/* Profile Details */}
            <div className="bg-white rounded-lg shadow p-6 mb-8">
                <h2 className="text-xl font-semibold mb-4">About</h2>
                {profileData.about ? (
                    <p className="text-gray-700 whitespace-pre-line">{profileData.about}</p>
                ) : (
                    <p className="text-gray-400 italic">No bio yet</p>
                )}

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-6">
                    <div>
                        <h3 className="text-sm font-medium text-gray-500">Date of Birth</h3>
                        <p>{profileData.dob || 'Not specified'}</p>
                    </div>
                    <div>
                        <h3 className="text-sm font-medium text-gray-500">Member Since</h3>
                        <p>{new Date(profileData.created_at).toLocaleDateString()}</p>
                    </div>
                </div>
            </div>
        </div>
    )
}