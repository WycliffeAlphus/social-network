'use client'

import { useState, useEffect } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'

export default function FollowingPage() {
    const { id } = useParams()
    const [following, setFollowing] = useState([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        const fetchFollowing = async () => {
            try {
                const response = await fetch(`http://localhost:8080/api/following/${id}`, {
                    method: 'GET',
                    credentials: 'include'
                })
                if (!response.ok) throw new Error('Failed to fetch following')
                const data = await response.json()
                setFollowing(Array.isArray(data?.Users) ? data.Users : [])
            } catch (err) {
                setError(err.message || 'Failed to load following')
            } finally {
                setLoading(false)
            }
        }
        fetchFollowing()
    }, [id])

    if (loading) return <div className="flex justify-center py-8">Loading...</div>
    if (error) return <div className="flex justify-center py-8 text-red-500">Error: {error}</div>

    return (
        <div className="max-w-2xl mx-auto py-8 px-4">
            <h1 className="text-2xl font-bold mb-6">Following</h1>
            <div className="space-y-4">
                {following.length > 0 ? (
                    following.map(user => (
                        <div key={user.ID} className="flex items-center gap-4 p-3 hover:bg-gray-100 rounded-lg">
                            <div className="w-12 h-12 rounded-full bg-gray-200 overflow-hidden flex items-center justify-center">
                                <span className="text-lg font-medium">
                                    {user.ID.substring(0, 2).toUpperCase()}
                                </span>
                            </div>
                            <div>
                                <Link href={`/profile/${user.ID}`} className="font-medium hover:underline">
                                    User ID: {user.ID}
                                </Link>
                                {user.Status && (
                                    <p className="text-sm text-gray-500">Status: {user.Status}</p>
                                )}
                            </div>
                        </div>
                    ))
                ) : (
                    <p>Not following anyone yet</p>
                )}
            </div>
        </div>
    )
}