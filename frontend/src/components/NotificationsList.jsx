import React, { useState, useEffect } from 'react';
import { getNotifications, markNotificationsAsRead } from '../lib/notifications';
import NotificationItem from './NotificationItem';

const NotificationsList = () => {
    const [notifications, setNotifications] = useState([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchNotifications = async () => {
            setIsLoading(true);
            const data = await getNotifications();
            console.log('Notifications data:', data);
            setNotifications(data || []);
            setIsLoading(false);
        };

        fetchNotifications();
    }, []);

    const handleMarkAllAsRead = async () => {
        const success = await markNotificationsAsRead();
        if (success) {
            setNotifications(notifications.map(n => ({ ...n, is_read: true })));
        }
    };

    return (
        <div className="absolute top-full right-0 mt-2 w-80 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg z-20">
            <div className="p-3 flex justify-between items-center border-b border-gray-200 dark:border-gray-700">
                <h3 className="font-semibold text-gray-800 dark:text-white">Notifications</h3>
                <button onClick={handleMarkAllAsRead} className="text-sm text-blue-600 hover:underline">
                    Mark all as read
                </button>
            </div>
            <div className="max-h-96 overflow-y-auto">
                {isLoading ? (
                    <p className="p-4 text-center text-gray-500">Loading...</p>
                ) : notifications.length === 0 ? (
                    <p className="p-4 text-center text-gray-500">No new notifications.</p>
                ) : (
                    notifications.map(notification => (
                        <NotificationItem key={notification.id} notification={notification} />
                    ))
                )}
            </div>
        </div>
    );
};

export default NotificationsList;
