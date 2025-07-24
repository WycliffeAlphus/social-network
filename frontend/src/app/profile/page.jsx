"use client";
import { useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import {
    HomeIcon, BellIcon, ChatBubbleLeftIcon,
    UserIcon, UserGroupIcon, PlusIcon
} from '@heroicons/react/24/solid';
import Navbar from '../components/navbar';
import Sidebar from '../components/sidebar';

export default function Profile() {
    const searchParams = useSearchParams();
    const id = searchParams.get("id");
    const [profile, setProfile] = useState(null);
    const [error, setError] = useState("");

    useEffect(() => {
        fetch(`http://localhost:8080/api/profile?id=${id ? `${id}` : ''}`, {
            credentials: 'include',
        })
            .then(res => {
                if (!res.ok) throw new Error("Not allowed");
                return res.json();
            })
            .then(data => {
                console.log("Profile data received:", data.data); 
                setProfile(data.data);
            })
            .catch(err => setError(err.message));
    }, [id]); // Add id as dependency

    if (error) return <p>Error: {error}</p>;
    if (!profile) return <p>Loading...</p>;


    const isPrivateProfile = profile.profile_visibility === "private";

    return (
        <div className="min-h-screen flex flex-col">
            <Navbar />
            <div className="flex flex-1">
                <Sidebar />
                <main className="flex-1 p-4">
                    <div className="p-6">
                        {isPrivateProfile ? (
                            <div>
                                <div className="bg-gray-900 text-white p-6">
                                    <div className="flex items-center justify-between max-w-6xl mx-auto">
                                        <div className="flex items-center space-x-6">
                                            <div className="relative">
                                                <img
                                                    src={profile.avatar || "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=120&h=120&fit=crop&crop=face"}
                                                    alt={`${profile.first_name} ${profile.last_name}`}
                                                    className="w-20 h-20 rounded-full object-cover border-2 border-gray-700"
                                                />
                                                <div className="absolute bottom-1 right-1 w-4 h-4 bg-green-500 rounded-full border-2 border-gray-900"></div>
                                            </div>

                                            <div>
                                                <h1 className="text-2xl font-bold mb-1">
                                                    {profile.first_name} {profile.last_name}
                                                </h1>
                                                <div className="flex items-center space-x-4 text-gray-400 text-sm">
                                                    <span>Profile is private</span>
                                                </div>
                                            </div>
                                        </div>

                                        <div className="flex items-center space-x-3">
                                            <button className="flex items-center space-x-2 bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded-lg font-medium transition-colors">
                                                <UserIcon className="w-5 h-5" />
                                                <span>Add friend</span>
                                            </button>

                                            <button className="flex items-center space-x-2 bg-gray-700 hover:bg-gray-600 px-4 py-2 rounded-lg font-medium transition-colors">
                                                <ChatBubbleLeftIcon className="w-5 h-5" />
                                                <span>Message</span>
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        ) : (
                            <div>
                                <h1>{profile.first_name} {profile.last_name}</h1>
                                <p>{profile.about}</p>
                                {/* Add more public profile content here */}
                            </div>
                        )}
                    </div>
                </main>
            </div>
        </div>
    );
}