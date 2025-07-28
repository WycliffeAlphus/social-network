"use client"

import Link from "next/link"
import { useEffect, useState } from "react"

export default function FollowSuggestion() {
    const [loading, setLoading] = useState(true)
    const [availableUsers, setAvailableUsers] = useState([])
    const [visibility, setVisibility] = useState("public")
    const [followingMap, setFollowingMap] = useState({})

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
            }
        } catch (err) {
            console.error('Error fetching available users:', err)
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
                // update map to turn button to following
                setFollowingMap((prev) => ({ ...prev, [userId]: true }))
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

                        let buttonLabel = "Follow";

                        if (followingMap[otherUser.id]) {
                            buttonLabel = "Following"
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
                                    </div>
                                </div>
                                <button
                                    onClick={() => handleFollow(otherUser.id)}
                                    className={`px-4 py-2 rounded-3xl text-white
                                        ${buttonLabel === "Following"
                                            ? "border border-gray-400 hover:bg-gray-400 hover:text-black"
                                            : "bg-blue-500 hover:bg-blue-600"
                                        }
                                        `}
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
