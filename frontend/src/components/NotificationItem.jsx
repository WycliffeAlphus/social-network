import React from 'react';
import { markNotificationAsRead } from '../lib/notifications';

// A simple component to format the time since the notification was created
const TimeAgo = ({ date }) => {
    // Basic time formatting, can be replaced with a library like date-fns
    const seconds = Math.floor((new Date() - new Date(date)) / 1000);
    let interval = seconds / 31536000;
    if (interval > 1) return Math.floor(interval) + "y ago";
    interval = seconds / 2592000;
    if (interval > 1) return Math.floor(interval) + "mo ago";
    interval = seconds / 86400;
    if (interval > 1) return Math.floor(interval) + "d ago";
    interval = seconds / 3600;
    if (interval > 1) return Math.floor(interval) + "h ago";
    interval = seconds / 60;
    if (interval > 1) return Math.floor(interval) + "m ago";
    return Math.floor(seconds) + "s ago";
};

const NotificationItem = ({ notification, onRead }) => {
    const itemClasses = `p-3 flex items-start gap-3 border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800`;

    const handleMarkAsRead = async () => {
        if (notification.is_read) return;
        const success = await markNotificationAsRead(notification.id);
        if (success) {
            onRead(notification.id);
        }
    };

    return (
        <div className={itemClasses} onClick={handleMarkAsRead}>
            {!notification.is_read && (
                <div className="w-2.5 h-2.5 bg-blue-500 rounded-full mt-1.5 flex-shrink-0"></div>
            )}
            <div className={`flex-grow ${notification.is_read ? 'ml-5' : ''}`}>
                <p className="text-sm text-gray-700 dark:text-gray-300">
                    {notification.message}
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    <TimeAgo date={notification.created_at} />
                </p>
            </div>
        </div>
    );
};

export default NotificationItem;
