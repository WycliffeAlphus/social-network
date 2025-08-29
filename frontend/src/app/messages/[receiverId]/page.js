"use client"

import { useUser } from "@/context/user-context";
import { useParams } from "next/navigation";
import { useRouter } from 'next/navigation';
import { useEffect, useRef, useState } from "react";
import { PaperAirplaneIcon, PhotoIcon, FaceSmileIcon, ExclamationTriangleIcon } from "@heroicons/react/24/outline";
import EmojiPicker from 'emoji-picker-react';
import { getWebSocket } from "@/components/ws";
import { sendMessage } from "@/components/ws";

export default function Messages() {
    const currentUserId = useUser()
    const { receiverId } = useParams()
    const router = useRouter()
    const [message, setMessage] = useState('')
    const [showEmojiPicker, setShowEmojiPicker] = useState(false)
    const emojiPickerRef = useRef(null)
    const emojiButtonRef = useRef(null)
    const [isDarkMode, setIsDarkMode] = useState(false)
    const [isLoading, setIsLoading] = useState(true)
    const [receiverExists, setReceiverExists] = useState(false)
    const [hasFollowRelationship, setHasFollowRelationship] = useState(false)

    useEffect(() => {
        if (currentUserId && currentUserId === receiverId) {
            router.push('/messages')
            return
        }

        // check user existence and follow relationship
        const checkUserAndRelationship = async () => {
            if (!currentUserId || !receiverId) return

            try {
                setIsLoading(true)
                const response = await fetch(`http://localhost:8080/api/follow-relationship?receiverId=${receiverId}`, {
                    credentials: 'include',
                })

                if (response.ok) {
                    const data = await response.json()
                    console.log(data)
                    setReceiverExists(data.messageReceiverExists)
                    setHasFollowRelationship(data.has_follow_relationship)
                }
            } catch (error) {
                console.error('Error checking relationship:', error)
                setReceiverExists(false)
                setHasFollowRelationship(false)
            } finally {
                setIsLoading(false)
            }
        }

        checkUserAndRelationship()
    }, [receiverId, currentUserId, router]);

    const onEmojiClick = (emojiData) => {
        setMessage(prevMessage => prevMessage + emojiData.emoji)
    }

    // this closes the emoji picker when clicking outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (emojiPickerRef.current && !emojiPickerRef.current.contains(event.target)) {
                setShowEmojiPicker(false)
            }
        }

        document.addEventListener('click', handleClickOutside)
        return () => {
            document.removeEventListener('click', handleClickOutside)
        }
    }, [showEmojiPicker])

    // calculate emoji picker position based on emoji button
    const getEmojiPickerPosition = () => {
        if (!emojiButtonRef.current) return {};

        const buttonRect = emojiButtonRef.current.getBoundingClientRect();
        return {
            bottom: `calc(100% - ${buttonRect.top}px)`,
            left: `${buttonRect.left}px`,
            transform: 'translateY(-10px)'
        };
    }

    useEffect(() => {
        const checkDarkMode = () => {
            setIsDarkMode(window.matchMedia('(prefers-color-scheme: dark)').matches);
        }

        // check initially
        checkDarkMode();

        // listen for system theme changes
        const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');

        const handleThemeChange = (e) => {
            setIsDarkMode(e.matches);
        };

        mediaQuery.addEventListener('change', handleThemeChange);

        return () => {
            mediaQuery.removeEventListener('change', handleThemeChange);
        };
    }, []);

    // closes emoji picker when window is resized
    useEffect(() => {
        const handleResize = () => {
            setShowEmojiPicker(false);
        };

        window.addEventListener('resize', handleResize);
        return () => {
            window.removeEventListener('resize', handleResize);
        };
    }, []);

    // prevent scrolling when emoji picker is open
    useEffect(() => {
        if (showEmojiPicker) {
            // add event listener to prevent scrolling
            const preventScroll = (e) => {
                // allow scrolling within the emoji picker
                if (emojiPickerRef.current && emojiPickerRef.current.contains(e.target)) {
                    return true; // allow scrolling inside emoji picker
                }

                e.preventDefault();
                e.stopPropagation();
                return false;
            };

            // prevent scroll on various events
            document.addEventListener('wheel', preventScroll, { passive: false });
            document.addEventListener('touchmove', preventScroll, { passive: false });
            document.addEventListener('keydown', (e) => {
                if (['Space', 'ArrowUp', 'ArrowDown', 'PageUp', 'PageDown'].includes(e.code)) {
                    e.preventDefault();
                }
            });

            // prevent scrollbar dragging by disabling overflow on ALL scrollable containers
            const scrollableContainers = document.querySelectorAll('.overflow-y-auto, .overflow-auto, [style*="overflow"], [class*="scroll"]');
            const originalStyles = [];

            scrollableContainers.forEach((container, index) => {
                originalStyles[index] = container.style.overflow;
                container.style.overflow = 'hidden';
            });

            // cleanup function to restore scrolling
            return () => {
                document.removeEventListener('wheel', preventScroll);
                document.removeEventListener('touchmove', preventScroll);

                // restore all scrollable containers
                scrollableContainers.forEach((container, index) => {
                    container.style.overflow = originalStyles[index];
                });
            };
        }
    }, [showEmojiPicker]);

    const handleSubmit = async (e) => {
        e.preventDefault()

        if (!currentUserId || !receiverId || !message.trim()) {
            return
        }

        const msg = {
            from: currentUserId,
            to: receiverId,
            content: message.trim()
        }

        sendMessage(msg)
        setMessage('') // clear input after sending message
    }

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-screen">
                <svg aria-hidden="true" className="w-8 h-8 text-gray-300 animate-spin dark:text-gray-600 fill-blue-600" viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" fill="currentColor" />
                    <path d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" fill="currentFill" />
                </svg>
            </div>
        )
    }

    if (!receiverExists) {
        return (
            <div></div>
        )
    }

    if (!hasFollowRelationship) {
        return (
            <div className="flex items-center justify-center h-screen">
                <div className="text-center p-6 mx-4">
                    <ExclamationTriangleIcon className="h-12 w-12 mx-auto mb-3" />
                    <h3 className="text-lg font-semibold mb-2">Whoa! Hold Up ⛔</h3>
                    <p className="">
                        🕵️‍♂️ No following, no messaging. Them's the rules. This ain't Tinder 😤.
                    </p>
                </div>
            </div>
        )
    }

    return (
        <div className="flex flex-col h-screen">
            {/* Chat header */}
            <div className="border-b border-gray-400 p-4">
                <h2 className="text-xl font-semibold">Chat</h2>
            </div>

            {/* Messages area */}
            <div className="flex-1 overflow-y-auto scrollbar p-3 space-y-4">
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
            <form onSubmit={handleSubmit} className="border-t border-gray-400 p-4">
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
                        ref={emojiButtonRef}
                        type="button"
                        className="p-2 rounded-full cursor-pointer"
                        onClick={() => setShowEmojiPicker(!showEmojiPicker)}
                    >
                        <FaceSmileIcon className="h-6 w-6 text-blue-500" />
                    </button>
                    <input
                        type="text"
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        placeholder="Type a message..."
                        className="flex-1 border border-gray-400 rounded-full py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                    <button
                        type="submit"
                        disabled={!message.trim()}
                        className="p-2 rounded-full bg-blue-500 text-white hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                    >
                        <PaperAirplaneIcon className="h-6 w-6" />
                    </button>
                </div>

                {/* Emoji Picker */}
                {showEmojiPicker && (
                    <div
                        ref={emojiPickerRef}
                        className="fixed scrollbar bottom-16 left-16 z-77"
                        style={getEmojiPickerPosition()}
                    >
                        <EmojiPicker
                            onEmojiClick={onEmojiClick}
                            width={350}
                            height={350}
                            theme={isDarkMode ? 'dark' : 'light'}
                        />
                    </div>
                )}
            </form>
        </div>
    )
}