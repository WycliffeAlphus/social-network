import Link from "next/link";
import { HomeIcon, BellIcon, ChatBubbleLeftIcon, UserCircleIcon } from '@heroicons/react/24/solid';
import { useState, useEffect, useRef } from 'react';
import NotificationsList from './NotificationsList';
import { getNotifications } from '../lib/notifications';

function Navbar() {
    const [showNotifications, setShowNotifications] = useState(false);
    const [unreadCount, setUnreadCount] = useState(0);
    const notificationRef = useRef(null);

    useEffect(() => {
        // Fetch initial unread count
        const fetchUnreadCount = async () => {
            const notifications = await getNotifications();
            const count = notifications.filter(n => !n.is_read).length;
            setUnreadCount(count);
        };
        fetchUnreadCount();

        // Poll for new notifications every 30 seconds
        const interval = setInterval(fetchUnreadCount, 30000);

        return () => clearInterval(interval);
    }, []);

    useEffect(() => {
        // Close dropdown when clicking outside
        function handleClickOutside(event) {
            if (notificationRef.current && !notificationRef.current.contains(event.target)) {
                setShowNotifications(false);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, [notificationRef]);

    const toggleNotifications = () => {
        setShowNotifications(!showNotifications);
        if (!showNotifications) {
            // When opening the list, reset count and refetch can be added here
            setUnreadCount(0); 
        }
    };

    return (
        <nav className='bg-blue-900 dark:text-white p-2 flex border-b justify-between border-gray-100/20' >
            <div className="left">
                <Link href="/" className='mr-2'><HomeIcon className="h-6 w-6 " /></Link>
            </div>
            <div className="right flex items-center">
                <Link href="/messages"><ChatBubbleLeftIcon className="h-6 w-6 mr-4" /></Link>
                
                <div className="relative" ref={notificationRef}>
                    <button onClick={toggleNotifications} className="relative">
                        <BellIcon className="h-6 w-6 mr-4" />
                        {unreadCount > 0 && (
                            <span className="absolute top-0 right-3 -mt-1 -mr-1 flex justify-center items-center w-4 h-4 bg-red-500 text-white text-xs rounded-full">
                                {unreadCount}
                            </span>
                        )}
                    </button>
                    {showNotifications && <NotificationsList />}
                </div>

                <Link href="/profile"><UserCircleIcon className="h-6 w-6 mr-3" /></Link>
            </div>
        </nav>
    );
}

export default Navbar;