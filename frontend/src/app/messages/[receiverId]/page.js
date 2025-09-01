"use client"

import { useUser } from "@/context/user-context";
import { useParams } from "next/navigation";
import { useRouter } from 'next/navigation';
import { useEffect, useRef, useState } from "react";
import { PaperAirplaneIcon, PhotoIcon, FaceSmileIcon, ExclamationTriangleIcon } from "@heroicons/react/24/outline";
import EmojiPicker from 'emoji-picker-react';
import { sendMessage } from "@/components/ws";
import Loading from "@/components/loading";
import { formattedMessageDate } from "@/lib/messageDate";

export default function Messages() {
    const currentUserId = useUser().current_user_id
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
    const [messages, setMessages] = useState([])
    const [messagesLoading, setMessagesLoading] = useState(true)
    const [receiverData, setReceiverData] = useState(null)
    const messagesEndRef = useRef(null)

    // auto-scroll to bottom when messages change
    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
    }

    useEffect(() => {
        scrollToBottom()
    }, [messages])

    useEffect(() => {
        if (currentUserId && currentUserId === receiverId) {
            router.push('/messages')
            return
        }

        // check user existence and follow relationship
        const checkUserAndRelationship = async () => {
            if (!currentUserId || !receiverId) return

            try {
                const response = await fetch(`http://localhost:8080/api/follow-relationship?receiverId=${receiverId}`, {
                    credentials: 'include',
                })

                if (response.ok) {
                    const data = await response.json()
                    console.log(data)
                    setHasFollowRelationship(data.hasFollowRelationship)
                    setReceiverData(data.receiverData)
                    setReceiverExists(data.messageReceiverExists)
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

    useEffect(() => {
        if (!currentUserId || !receiverId || currentUserId === receiverId) {
            return;
        }

        const getConversation = async () => {
            try {
                const response = await fetch(`http://localhost:8080/api/conversations?receiverId=${receiverId}`, {
                    credentials: 'include',
                })

                if (response.status === 400) {
                    return
                }

                if (!response.ok) {
                    throw new Error('Failed to fetch profile')
                }

                const messages = await response.json()
                setMessages(messages)
            } catch (error) {
                console.error('Error fetching conversation:', error);
            } finally {
                setMessagesLoading(false)
            }
        }

        getConversation()
    }, [receiverId, currentUserId])

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

            // prevent arrow keys and spacebar from scrolling the page
            const preventKeyScroll = (e) => {
                if (['Space', 'ArrowUp', 'ArrowDown', 'PageUp', 'PageDown'].includes(e.code)) {
                    e.preventDefault();
                }
            };

            // prevent scroll on various events
            document.addEventListener('wheel', preventScroll, { passive: false });
            document.addEventListener('touchmove', preventScroll, { passive: false });
            document.addEventListener('keydown', preventKeyScroll);

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
                document.removeEventListener('keydown', preventKeyScroll);

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
            content: message.trim(),
            timestamp: formattedMessageDate(),
        }

        sendMessage(msg)
        setMessages(prevMessages => [...prevMessages, msg]); // adds the message to the UI immediately (optimistic update)
        setMessage('') // clear input after sending message

        // create a custom event that will be listened accross the app when a message is sent
        window.dispatchEvent(new CustomEvent('messageEvent'))
    }

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-screen">
                <Loading />
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
            <div className="border-b flex items-center gap-3 border-gray-400 p-4">
                {receiverData.avatar?.String ? (
                    <img
                        src={receiverData.avatar.String}
                        alt="User avatar"
                        className="w-12 h-12 rounded-full object-cover"
                    />
                ) : (
                    <div className="w-[clamp(2rem,4vw,3.5rem)] h-[clamp(2rem,4vw,3.5rem)] rounded-full bg-gray-200 overflow-hidden flex items-center justify-center">
                        <span className="text-lg text-black font-medium">
                            {receiverData.firstname?.charAt(0).toUpperCase()}
                            {receiverData.lastname?.charAt(0).toUpperCase()}
                        </span>
                    </div>
                )}

                <div>
                    <span className="font-medium group-hover:text-[#4169e1]">
                        {receiverData.firstname} {receiverData.lastname}
                    </span>
                </div>
            </div>

            {/* Messages area */}
            {!messagesLoading ? (
                <div className="flex-1 overflow-y-auto scrollbar p-3 space-y-4">
                    {messages && messages.length > 0 ? (
                        messages.map((message) => (
                            <div key={`${message.timestamp}-${message.from}`} className={`flex ${message.from === currentUserId ? 'justify-end' : 'justify-start'}`}>
                                <div className={`text-white rounded-3xl p-3 max-w-xs ${message.from === currentUserId ? 'bg-blue-500 rounded-br-sm' : 'bg-gray-700 rounded-bl-sm'}`}>
                                    <p>{message.content}</p>
                                    <p className="text-xs font-bold text-gray-300 mt-1">{formattedMessageDate(message.timestamp)}</p>
                                </div>
                            </div>
                        ))
                    ) : (
                        <h3>Begin your conversation</h3>
                    )}
                    <div ref={messagesEndRef} /> {/* invsible element for auto-scrolling */}
                </div>
            ) : (
                <Loading />
            )}

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