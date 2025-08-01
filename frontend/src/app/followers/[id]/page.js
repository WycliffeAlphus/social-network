'use client'

import { useState, useEffect } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import Rightbar from '@/components/rightbar'
import FollowSuggestion from '@/components/followsuggestions'

export default function FollowersPage() {
    const { id } = useParams()
    const [followers, setFollowers] = useState([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        const fetchFollowers = async () => {
            try {
                const response = await fetch(`http://localhost:8080/api/followers/${id}`, {
                    method: 'GET',
                    credentials: 'include'
                })
                if (!response.ok) throw new Error('Failed to fetch followers')

                const data = await response.json()
                // console.log(data)
                setFollowers(Array.isArray(data?.users) ? data.users : [])
            } catch (err) {
                setError(err.message || 'Failed to load followers')
            } finally {
                setLoading(false)
            }
        }
        fetchFollowers()
    }, [id])

    if (loading) return <div className="flex justify-center py-8">Loading...</div>
    if (error) return <div className="flex justify-center py-8 text-red-500">Error: {error}</div>

    return (
        <div className="flex min-h-screen">
            <main className="flex-1 border-x mr-[20px] border-gray-400">
                <div className="lg:hidden">
                    <FollowSuggestion />
                </div>
                <h1 className="py-[0.7rem] px-[1rem] bg-ble-500 text-2xl font-bold border-t lg:border-0 border-gray-400">Followers</h1>
                <div className="p-7 border-t border-gray-400">
                    {followers.length > 0 ? (
                        followers.map(user => (
                            <Link
                                key={user.id}
                                href={`/profile/${user.id}`}
                                className="group flex items-center gap-3 p-3 rounded-lg transition-colors w-fit"
                            >
                                {user.imgurl ? (
                                    <img
                                        src={user.imgurl}
                                        alt={`${user.fname} ${user.lname}`}
                                        className="w-12 h-12 rounded-full object-cover"
                                    />
                                ) : (
                                    <div className="w-12 h-12 rounded-full bg-gray-200 overflow-hidden flex items-center justify-center">
                                        <span className="text-lg text-black font-medium">
                                            {user.fname?.charAt(0).toUpperCase()}
                                            {user.lname?.charAt(0).toUpperCase()}
                                        </span>
                                    </div>
                                )}
                                <div>
                                    <span className="font-medium group-hover:text-[#4169e1]">
                                        {user.fname} {user.lname}
                                    </span>
                                </div>
                            </Link>
                        ))
                    ) : (
                        <p>No followers yet</p>
                    )}
                </div>
            </main>
            <Rightbar />
        </div>
    )
}