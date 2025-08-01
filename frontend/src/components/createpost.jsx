"use client"

import { showFieldError } from "@/lib/auth";
import { PhotoIcon, XMarkIcon } from "@heroicons/react/24/outline";
import { useMemo, useState } from "react";

export default function CreatePost({ onClose }) {
    const [showFollowers, setShowFollowers] = useState(false);
    const [selectedPrivacy, setSelectedPrivacy] = useState('public');
    const [error, setError] = useState('')
    const [formData, setFormData] = useState({
        title: '',
        content: '',
        postPrivacy: '',
        postImage: '',
    })

    // memoize required fields to avoid recreation on every render
    const requiredFields = useMemo(() => [
        'title', 'content', 'privacy'
    ], [])

    const handlePrivacyChange = (privacy) => {
        setSelectedPrivacy(privacy);
        setShowFollowers(privacy === 'private');
        setFormData(prev => ({ ...prev, postPrivacy: privacy }));
    };

    const handleChange = (e) => {
        const { name, value } = e.target
        showFieldError(name, '');

        if (requiredFields.includes(name) && !value) {
            showFieldError(name, "This field is needed")
        }

        setFormData(prev => ({ ...prev, [name]: value }))
        console.log(formData)
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setError('')

        try {
            const formDataToSend = new FormData()

            // append all form data to FormData object
            Object.entries(formData).forEach(([key, value]) => {
                if (value !== null && value !== undefined) {
                    if (key === 'postImage' && value instanceof File) {
                        formDataToSend.append(key, value)
                    } else {
                        // console.log(key, value)
                        formDataToSend.append(key, value)
                    }
                }
            })

            console.log(Object.fromEntries(formDataToSend.entries()))
        } catch (err) {
            setError('Post creation failed. Please try again.')
            console.error('Registration error:', err)
        }
    }

    return (
        <form onSubmit={handleSubmit} onClick={(e) => e.stopPropagation()} className="h-fit bg-white dark:bg-black max-w-[40rem] w-[80%] rounded-2xl p-6 space-y-3">
            {/* Page Header */}
            <div className="flex justify-between">
                {error && <p className="text-red-500 mb-4">{error}</p>}
                <h1 className="text-xl font-bold text-gray-800 dark:text-white">Wassup</h1>
                <button
                    onClick={onClose}
                    className="top-4 cursor-pointer right-4 p-1 rounded-full"
                >
                    <XMarkIcon className="h-6 w-6 text-gray-500 dark:text-white hover:text-blue-500" />
                </button>
            </div>
            {/* Title */}
            <div>
                <label htmlFor="title" className="block text-sm text-gray-700 dark:text-white">Post Title</label>
                <input
                    type="text"
                    id="title"
                    name="title"
                    value={formData.title}
                    onChange={handleChange}
                    placeholder="What's your post about?"
                    className="mt-1 w-full py-1 px-3 border-b border-gray-700 focus:outline-none focus:border-b-blue-500"
                    required
                />
                <div id="title-error" className="text-red-500"></div>
            </div>

            {/* Content */}
            <div>
                <label htmlFor="content" className="block text-sm text-gray-700 dark:text-white">Post Content</label>
                <textarea
                    id="content"
                    name="content"
                    value={formData.content}
                    onChange={handleChange}
                    className="mt-1 w-full max-h-[9rem] min-h-[9rem] p-3 border-b border-gray-700 focus:outline-none focus:border-b-blue-500"
                    placeholder="Write your post content here..."
                    required
                ></textarea>
                <div id="content-error" className="text-red-500"></div>
            </div>

            {/* Privacy Settings */}
            <div>
                <h2 className="text-sm mb-3 text-gray-700 dark:text-white">Privacy Settings</h2>

                <div className="space-y-3">
                    {/* Public Option */}
                    <div
                        className={`flex items-center px-3 py-1 border rounded-lg cursor-pointer ${selectedPrivacy === 'public' ? 'border-blue-500' : 'border-gray-700 hover:border-blue-500'}`}
                        onClick={() => handlePrivacyChange('public')}
                    >
                        <input
                            id="public"
                            name="postPrivacy"
                            type="radio"
                            className="h-4 w-4 text-indigo-600 pointer-events-none"
                            checked={selectedPrivacy === 'public'}
                            value="public"
                            onChange={() => handlePrivacyChange('public')}
                        />
                        <label htmlFor="public" className="text-sm ml-3 flex-grow pointer-events-none">
                            <span className="dark:text-white text-gray-700">Public</span>
                            <span className="block text-gray-500">All users will be able to see this post</span>
                        </label>
                    </div>

                    {/* Followers Only Option */}
                    <div
                        className={`flex items-center px-3 py-1 border rounded-lg cursor-pointer ${selectedPrivacy === 'almostprivate' ? 'border-blue-500' : 'border-gray-700 hover:border-blue-500'}`}
                        onClick={() => handlePrivacyChange('almostprivate')}
                    >
                        <input
                            id="almostprivate"
                            name="postPrivacy"
                            type="radio"
                            className="h-4 w-4 text-indigo-600 pointer-events-none"
                            checked={selectedPrivacy === 'almostprivate'}
                            value="almostprivate"
                            onChange={() => handlePrivacyChange('almostprivate')}
                        />
                        <label htmlFor="almostprivate" className="text-sm ml-3 flex-grow pointer-events-none">
                            <span className="dark:text-white text-gray-700">Almost private</span>
                            <span className="block text-gray-500">Only your followers will be able to see this post</span>
                        </label>
                    </div>

                    {/* Private Option */}
                    <div
                        className={`flex items-center px-3 py-1 border rounded-lg cursor-pointer ${selectedPrivacy === 'private' ? 'border-blue-500' : 'border-gray-700 hover:border-blue-500'}`}
                        onClick={() => handlePrivacyChange('private')}
                    >
                        <input
                            id="private"
                            name="postPrivacy"
                            type="radio"
                            className="h-4 w-4 text-indigo-600 pointer-events-none"
                            checked={selectedPrivacy === 'private'}
                            value="private"
                            onChange={() => handlePrivacyChange('private')}
                        />
                        <label htmlFor="private" className="text-sm ml-3 flex-grow pointer-events-none">
                            <span className="dark:text-white text-gray-700">Private</span>
                            <span className="block text-gray-500">Only specific followers you choose will see this post</span>
                        </label>
                    </div>
                </div>

                {/* Follower Selector (shown only when private is selected) */}
                {showFollowers && (
                    <div className="shadow-[0_0_12px_rgba(0,0,0,0.5),0_0_12px_rgba(255,255,255,0.5)] dark:bg-black bg-white mt-3 p-4 rounded-lg">
                        <h3 className="text-sm mb-3">Select which followers can see this post:</h3>

                        {/* Follower List */}
                        <div className="max-h-48 overflow-y-auto space-y-2">
                            <div className="group flex items-center justify-between px-3">
                                <label htmlFor="1" className="cursor-pointer flex items-center flex-1">
                                    <img src="https://randomuser.me/api/portraits/women/44.jpg" alt="Follower" className="h-8 w-8 rounded-full" />
                                    <span className="ml-3 text-sm group-hover:text-blue-500">Sarah Johnson</span>
                                </label>
                                <input id="1" type="checkbox" className="cursor-pointer h-4 w-4 text-indigo-600 rounded" />
                            </div>

                            <div className="flex items-center justify-between p-3 rounded-lg">
                                <div className="flex items-center">
                                    <img src="https://randomuser.me/api/portraits/men/32.jpg" alt="Follower" className="h-8 w-8 rounded-full" />
                                    <span className="ml-3 text-sm ">Michael Chen</span>
                                </div>
                                <input type="checkbox" className="h-4 w-4 text-indigo-600 rounded" />
                            </div>

                            <div className="flex items-center justify-between p-3 rounded-lg">
                                <div className="flex items-center">
                                    <img src="https://randomuser.me/api/portraits/men/75.jpg" alt="Follower" className="h-8 w-8 rounded-full" />
                                    <span className="ml-3 text-sm ">David Rodriguez</span>
                                </div>
                                <input type="checkbox" className="h-4 w-4 text-indigo-600 rounded" />
                            </div>
                        </div>

                        <div className="mt-4">
                            <span className="text-xs text-gray-500">0 followers selected</span>
                        </div>
                    </div>
                )}
            </div>

            {/* Image Upload */}
            <div className="mt-6">
                <label className="block text-sm text-gray-700 dark:text-white">Featured Image</label>
                <div className="mt-1 flex items-center">
                    <label htmlFor="image" className="cursor-pointer">
                        <PhotoIcon className="h-5 w-5 text-blue-500" />
                    </label>
                    <input id="image" type="file" className="hidden" accept="image/*" />
                    <p className="ml-3 text-xs text-gray-500">PNG, JPG, GIF up to 20MB</p>
                </div>
            </div>

            {/* Form Actions - Centered Publish Button */}
            <div className="flex justify-center">
                <button
                    type="submit"
                    className="cursor-pointer py-2 px-4 border border-transparent rounded-3xl shadow-sm text-base text-white bg-blue-500 hover:bg-blue-600"
                >
                    Post
                </button>
            </div>
        </form>
    )
}
