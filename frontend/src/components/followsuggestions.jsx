"use client"

import Link from "next/link"
import { useEffect, useState } from "react"
import Loading from "./loading"

export default function FollowSuggestion() {
    const [loading, setLoading] = useState(true)
    const [availableUsers, setAvailableUsers] = useState([])
    const [visibility, setVisibility] = useState("public")
    const [followStatusMap, setFollowStatusMap] = useState({})

    useEffect(() => {
        const loadUser = async () => {
            // Fetch available users and follow stats
            await Promise.all([
                fetchAvailableUsers(),
            ])
            setLoading(false)
        }
        loadUser()
    }, [])

    const fetchAvailableUsers = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/users/available', {
                credentials: 'include'
            })
            if (response.ok) {
                const data = await response.json()
                setAvailableUsers(data.users)
                setVisibility(data.visibility)

                // Fetch follow status for each user
                await Promise.all(
                    data.users.map(async (user) => {
                        await fetchFollowStatus(user.id)
                    })
                )
            }
        } catch (err) {
            console.error('Error fetching available users:', err)
        }
    }

    const fetchFollowStatus = async (userId) => {
        try {
            const response = await fetch(`http://localhost:8080/api/follow-status/${userId}`, {
                credentials: 'include'
            })
            if (response.ok) {
                const data = await response.json()
                setFollowStatusMap(prev => ({ ...prev, [userId]: data.status }))
            }
        } catch (err) {
            console.error('Error fetching follow status:', err)
        }
    }

    const handleFollow = async (userId) => {
        const userToFollow = availableUsers.find(user => user.id === userId);

        // Optimistically update the UI
        const newStatus = userToFollow && userToFollow.visibility === 'private' ? 'requested' : 'accepted';
        setFollowStatusMap(prev => ({ ...prev, [userId]: newStatus }));

        try {
            const response = await fetch('http://localhost:8080/api/users/follow', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ userId }),
                credentials: 'include'
            });

            if (!response.ok) {
                // If the request fails, revert the change
                setFollowStatusMap(prev => ({ ...prev, [userId]: 'not_following' }));
                console.error('Failed to follow user');
            }
            // No need to parse the response if we are updating optimistically
        } catch (err) {
            // If the request fails, revert the change
            setFollowStatusMap(prev => ({ ...prev, [userId]: 'not_following' }));
            console.error('Error following user:', err)
        }
    };

    const handleCancelRequest = async (userId) => {
        try {
            const response = await fetch('http://localhost:8080/api/follow/cancel', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ userId }),
                credentials: 'include'
            });

            if (response.ok) {
                setFollowStatusMap(prev => ({ ...prev, [userId]: 'not_following' }))
            } else {
                console.error('Failed to cancel follow request');
            }
        } catch (err) {
            console.error('Error canceling follow request:', err)
        }
    };

    if (loading) {
        return <div className="flex items-center justify-center h-screen"><Loading /></div>
    }

    return (
        <div className="border-0 lg:border border-gray-400 rounded-xl p-4">
            <h2 className="text-xl font-semibold mb-4">Profiles you can follow</h2>
            <div className="space-y-4">
                {availableUsers && availableUsers.length > 0 ? (
                    availableUsers.map((otherUser) => {
                        const followStatus = followStatusMap[otherUser.id] || 'not_following'

                        let buttonLabel = "Follow";
                        let buttonClass = "bg-blue-500 hover:bg-blue-600";

                        if (followStatus === 'accepted') {
                            buttonLabel = "Following"
                            buttonClass = "border border-gray-400 hover:bg-gray-400 hover:text-black"
                            isDisabled = true;
                        } else if (followStatus === 'requested') {
                            buttonLabel = "Requested"
                            buttonClass = "bg-gray-300 text-gray-600 hover:bg-gray-400 hover:text-black"
                        } else if (otherUser.followsMe && visibility !== "private") {
                            buttonLabel = "Follow back"
                        }

                        return (
                            <div key={otherUser.id} className="flex items-center justify-between p4">
                                <Link
                                    href={`/profile/${otherUser.id}`}
                                    className="group flex items-center gap-3 rounded-lg transition-colors w-fit"
                                >
                                    {otherUser.avatar?.String ? (
                                        <img
                                            src={otherUser.avatar.String}
                                            alt="User avatar"
                                            className="w-12 h-12 rounded-full object-cover"
                                        />
                                    ) : (
                                        <div className="w-12 h-12 rounded-full bg-gray-200 overflow-hidden flex items-center justify-center">
                                            <span className="text-lg text-black font-medium">
                                                {otherUser.firstName?.charAt(0).toUpperCase()}
                                                {otherUser.lastName?.charAt(0).toUpperCase()}
                                            </span>
                                        </div>
                                    )}
                                    <div>
                                        <span className="font-medium group-hover:text-[#4169e1]">
                                            {otherUser.firstName} {otherUser.lastName}
                                        </span>
                                    </div>
                                </Link>
                                <button
                                    onClick={() => {
                                        if (followStatus === 'requested') {
                                            handleCancelRequest(otherUser.id)
                                        } else {
                                            handleFollow(otherUser.id)
                                        }
                                    }}
                                    className={`px-4 py-2 rounded-3xl text-white ${buttonClass}`}
                                    disabled={isDisabled}
                                >
                                    {buttonLabel}
                                </button>
                            </div>
                        )
                    })
                ) : (
                    <p>No users available to follow</p>
                )}
                {availableUsers && availableUsers.length > 3 && (
                    <div>
                        <Link href="/explore/people" className="hover:text-blue-700 text-[#4169e1]">
                            Show more
                        </Link>
                    </div>
                )}
            </div>
        </div>
    )
}