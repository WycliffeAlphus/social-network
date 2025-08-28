"use client"

import { useUser } from "@/context/user-context";
import { useParams } from "next/navigation";
import { useRouter } from 'next/navigation';
import { useEffect, useState } from "react";
import { PaperAirplaneIcon, PhotoIcon, FaceSmileIcon } from "@heroicons/react/24/outline";

export default function Messages() {
    const currentUserId = useUser()
    const { id } = useParams()
    const router = useRouter()

    const [message, setMessage] = useState('')

    useEffect(() => {
        if (currentUserId && currentUserId === id) {
            router.push('/messages')
            return
        }
    }, [id, currentUserId, router]);

    return (
        <div className="flex flex-col h-screen">
            {/* Chat header */}
            <div className="border-b border-gray-200 p-4">
                <h2 className="text-xl font-semibold">Chat</h2>
            </div>

            {/* Messages area */}
            <div className="flex-1 overflow-y-auto p-3 space-y-4">
                {/* Example messages */}
                <div className="flex justify-start">
                    <div className="bg-gray-700 text-white rounded-3xl rounded-bl-sm p-3 max-w-xs">
                        <p>Hello there! How are you doing?</p>
                        <p className="text-xs text-gray-500 mt-1">10:30 AM</p>
                    </div>
                </div>
                <div className="flex justify-end">
                    <div className="bg-blue-500 text-white rounded-3xl rounded-br-sm p-3 max-w-xs">
                        <p>I'm good, thanks for asking!</p>
                        <p className="text-xs text-blue-100 mt-1">10:32 AM</p>
                    </div>
                </div>
            </div>

            {/* Message input form */}
            <form className="border-t border-gray-200 p-4">
                <div className="flex items-center space-x-2">
                    <div>
                        <label htmlFor="chatImage" className="cursor-pointer">
                            <PhotoIcon className="h-5 w-5 text-blue-500" />
                        </label>
                        <input
                            id="chatImage"
                            name="chatImage"
                            type="file"
                            className="hidden"
                            accept="image/*"
                        // onChange={handleFileChange}
                        />
                    </div>
                    <button
                        type="button"
                        className="p-2 rounded-full cursor-pointer"
                    >
                        <FaceSmileIcon className="h-6 w-6 text-blue-500" />
                    </button>
                    <input
                        type="text"
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        placeholder="Type a message..."
                        className="flex-1 border border-gray-300 rounded-full py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                    <button
                        type="submit"
                        disabled={!message.trim()}
                        className="p-2 rounded-full bg-blue-500 text-white hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                    >
                        <PaperAirplaneIcon className="h-6 w-6" />
                    </button>
                </div>
            </form>
        </div>
    )
}