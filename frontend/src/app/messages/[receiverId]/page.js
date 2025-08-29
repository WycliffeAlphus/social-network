"use client"

import { useUser } from "@/context/user-context";
import { useParams } from "next/navigation";
import { useRouter } from 'next/navigation';
import { useEffect, useRef, useState } from "react";
import { PaperAirplaneIcon, PhotoIcon, FaceSmileIcon } from "@heroicons/react/24/outline";
import EmojiPicker from 'emoji-picker-react';
import { getWebSocket } from "@/components/ws";

export default function Messages() {
    const currentUserId = useUser()
    const { receiverId } = useParams()
    const router = useRouter()
    const [message, setMessage] = useState('')
    const [showEmojiPicker, setShowEmojiPicker] = useState(false)
    const emojiPickerRef = useRef(null)
    const emojiButtonRef = useRef(null)
    const [isDarkMode, setIsDarkMode] = useState(false)

    useEffect(() => {
        if (currentUserId && currentUserId === receiverId) {
            router.push('/messages')
            return
        }
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

        const sock = getWebSocket()
        if (sock && sock.readyState === WebSocket.OPEN) {
            sock.send(JSON.stringify(msg))
            setMessage('') // clear input after sending message
        } else {
            console.error("WebSocket not connected")
        }
    }

    return (
        <div className="flex flex-col h-screen">
            {/* show the below if valid */}
            {/* if not valid show cannot send message to the user */}
            {/* if ID not in the db, giza */}

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