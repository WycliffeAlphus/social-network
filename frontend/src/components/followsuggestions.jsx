"use client"

import Link from "next/link"
import { useEffect, useState } from "react"

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
        try {
            const response = await fetch('http://localhost:8080/api/users/follow', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ userId }),
                credentials: 'include'
            })

            if (response.ok) {
                const data = await response.json();
                setFollowStatusMap(prev => ({ ...prev, [userId]: data.status }))
            }
        } catch (err) {
            console.error('Error following user:', err)
        }
    }

    if (loading) {
        return <div className="container mx-auto p-4">Loading...</div>
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
                        let isDisabled = false;

                        if (followStatus === 'accepted') {
                            buttonLabel = "Following"
                            buttonClass = "border border-gray-400 hover:bg-gray-400 hover:text-black"
                        } else if (followStatus === 'requested') {
                            buttonLabel = "Requested"
                            buttonClass = "bg-gray-300 text-gray-600 cursor-not-allowed"
                            isDisabled = true
                        } else if (otherUser.followsMe && visibility !== "private") {
                            buttonLabel = "Follow back"
                        }

                        return (
                            <div key={otherUser.id} className="flex items-center justify-between p-4">
                                <div className="flex items-center">
                                    {otherUser.avatar?.String ? (
                                        <img
                                            src={otherUser.avatar.String}
                                            className="w-12 h-12 object-cover rounded-full mr-3"
                                            alt="User avatar"
                                        />
                                    ) : null}
                                    <div>
                                        <p className="text-sm text-gray-500">
                                            {otherUser.firstName} {otherUser.lastName}
                                        </p>
                                        {followStatus === 'requested' && (
                                            <p className="text-xs text-gray-400 mt-1">
                                                This account is private
                                            </p>
                                        )}
                                    </div>
                                </div>
                                <button
                                    onClick={() => !isDisabled && handleFollow(otherUser.id)}
                                    disabled={isDisabled}
                                    className={`px-4 py-2 rounded-3xl text-white ${buttonClass}`}
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
