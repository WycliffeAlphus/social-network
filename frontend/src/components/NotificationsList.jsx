import React, { useState, useEffect } from 'react';
import { getNotifications, markNotificationsAsRead, getFollowStatuses, getGroupInviteStatuses } from '../lib/notifications';
import NotificationItem from './NotificationItem';

const NotificationsList = () => {
    const [notifications, setNotifications] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [followStatuses, setFollowStatuses] = useState({});
    const [groupInviteStatuses, setGroupInviteStatuses] = useState({});

    useEffect(() => {
        const fetchNotifications = async () => {
            setIsLoading(true);
            const data = await getNotifications();
            console.log('Notifications data:', data);
            setNotifications(data || []);

            const followRequestUserIds = data
                .filter(n => n.type === 'follow_request')
                .map(n => n.actor_id);

            if (followRequestUserIds.length > 0) {
                const statuses = await getFollowStatuses(followRequestUserIds);
                setFollowStatuses(statuses);
            }

            const groupInviteIds = data
                .filter(n => n.type === 'group_invite')
                .map(n => n.content_id);

            if (groupInviteIds.length > 0) {
                const statuses = await getGroupInviteStatuses(groupInviteIds);
                setGroupInviteStatuses(statuses);
            }

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

    const handleMarkAsRead = (notificationId) => {
        setNotifications(notifications.map(n =>
            n.id === notificationId ? { ...n, is_read: true } : n
        ));
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
                        <NotificationItem 
                            key={notification.id} 
                            notification={notification} 
                            onRead={handleMarkAsRead} 
                            followStatus={followStatuses[notification.actor_id]}
                            groupInviteStatus={groupInviteStatuses[notification.content_id]}
                        />
                    ))
                )}
            </div>
        </div>
    );
};

export default NotificationsList;
