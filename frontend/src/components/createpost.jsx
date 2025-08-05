"use client"

import { showFieldError } from "@/lib/auth";
import { validateAndPreviewFile } from "@/lib/filevalidater";
import { PhotoIcon, XMarkIcon } from "@heroicons/react/24/outline";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";

export default function CreatePost({ onClose }) {
    const [showFollowers, setShowFollowers] = useState(false);
    const [selectedPrivacy, setSelectedPrivacy] = useState('public');
    const [error, setError] = useState('')
    const [loadingFollowers, setLoadingFollowers] = useState(false);
    const [followersError, setFollowersError] = useState(null);
    const [followers, setFollowers] = useState([]);
    const [selectedFollowers, setSelectedFollowers] = useState([]);
    const [preview, setPreview] = useState('') // for image preview
    const router = useRouter()
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

    // fetch followers when private mode is selected
    useEffect(() => {
        if (showFollowers) {
            fetchFollowers();
        }
    }, [showFollowers]);

    const fetchFollowers = async () => {
        setLoadingFollowers(true);
        setFollowersError(null);

        try {
            const response = await fetch('http://localhost:8080/api/followers/currentuser', {
                method: 'GET',
                credentials: 'include'
            });
            if (!response.ok) {
                throw new Error('Failed to fetch followers');
            }
            const data = await response.json();
            setFollowers(data.users || []);
        } catch (err) {
            setFollowersError(err.message);
            console.error('Error fetching followers:', err);
        } finally {
            setLoadingFollowers(false);
        }
    };

    const handleFollowerSelection = (followerId) => {
        setSelectedFollowers(prev => {
            if (prev.includes(followerId)) {
                return prev.filter(id => id !== followerId);
            } else {
                return [...prev, followerId];
            }
        });
    };

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
    }

    const handleFileChange = (e) => {
        showFieldError('postImage', '');
        const file = e.target.files[0];

        if (!file) return;

        const result = validateAndPreviewFile(file, setPreview);

        if (!result.valid) {
            showFieldError('postImage', result.error);
            e.target.value = ''; // clear the file input
            return;
        }

        setFormData(prev => ({ ...prev, postImage: file }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault()
        setError('')

        // ensure at least one follower is selected when post privacy is private
        if (selectedPrivacy === 'private' && selectedFollowers.length === 0) {
            showFieldError('private', "Please select at least one follower for private posts");
            return;
        }

        try {
            const formDataToSend = new FormData()

            // include selected followers if private mode
            if (selectedPrivacy === 'private') {
                formDataToSend.append('allowedFollowers', JSON.stringify(selectedFollowers));
            }

            // append all form data to FormData object
            Object.entries(formData).forEach(([key, value]) => {
                if (value !== null && value !== undefined) {
                    if (key === 'postImage' && value instanceof File) {
                        formDataToSend.append(key, value)
                    } else if (key !== 'avatarImage') {
                        formDataToSend.append(key, value)
                    }
                }
            })

            console.log(Object.fromEntries(formDataToSend.entries()))

            const response = await fetch('http://localhost:8080/api/createpost', {
                method: 'POST',
                credentials: 'include',
                body: formDataToSend,
            })

            if (!response.ok) {
                const data = await response.json();

                if (data.titleerror) {
                    showFieldError('title', data.titleerror);
                }
                if (data.contenterror) {
                    showFieldError('content', data.contenterror);
                }
                if (data.privacyerror) {
                    showFieldError('postPrivacy', data.privacyerror);
                }
                if (data.imageerror) {
                    showFieldError('postImage', data.imageerror);
                }
            }

            if (response.ok) {
                router.push('/') // redirect to home after post creation
                onClose()
            }
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
                    type="button"
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
                    <div id="private-error" className="text-red-500"></div>
                </div>

                {/* Follower Selector (shown only when private is selected) */}
                {showFollowers && (
                    <div className="shadow-[0_0_12px_rgba(0,0,0,0.5),0_0_12px_rgba(255,255,255,0.5)] dark:bg-black bg-white mt-3 p-4 rounded-lg">
                        <h3 className="text-sm mb-3">Select which followers can see this post:</h3>

                        {loadingFollowers && <p className="text-sm text-gray-500">Loading followers...</p>}
                        {followersError && <p className="text-sm text-red-500">{followersError}</p>}

                        {/* Follower List */}
                        <div className="max-h-48 overflow-y-auto space-y-2">
                            {followers.map(follower => (
                                <div key={follower.id} className="group flex items-center justify-between px-7">
                                    <label htmlFor={`${follower.id}`} className="cursor-pointer flex items-center flex-1">
                                        <img
                                            src={follower.ImgURL || "/default-avatar.png"}
                                            alt="folloer image"
                                            className="h-8 w-8 rounded-full"
                                        />
                                        <span className="ml-3 text-sm group-hover:text-blue-500">
                                            {follower.fname} {follower.lname}
                                        </span>
                                    </label>
                                    <input
                                        id={`${follower.id}`}
                                        type="checkbox"
                                        checked={selectedFollowers.includes(follower.id)}
                                        onChange={() => handleFollowerSelection(follower.id)}
                                        className="cursor-pointer h-4 w-4 text-indigo-600 rounded"
                                    />
                                </div>
                            ))}
                        </div>

                        <div className="mt-4">
                            <span className="text-xs text-gray-500">
                                {selectedFollowers.length} {selectedFollowers.length === 1 ? 'follower' : 'followers'} selected
                            </span>
                        </div>
                    </div>
                )}
            </div>

            {/* Image Upload */}
            <div className="mt-6">
                <label className="block text-sm text-gray-700 dark:text-white">Featured Image</label>
                <div className="mt-1 flex items-center">
                    <label htmlFor="postImage" className="cursor-pointer">
                        <PhotoIcon className="h-5 w-5 text-blue-500" />
                    </label>
                    <input
                        id="postImage"
                        name="postImage"
                        type="file"
                        className="hidden"
                        accept="image/*"
                        onChange={handleFileChange}
                    />
                    <p className="ml-3 text-xs text-gray-500">PNG, JPG, GIF up to 20MB</p>
                </div>
                {preview && (
                    <img
                        src={preview}
                        alt="Preview"
                        className="w-35 mt-3 object-cover"
                    />
                )}
                <div id="postImage-error" className="text-red-500"></div>
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
