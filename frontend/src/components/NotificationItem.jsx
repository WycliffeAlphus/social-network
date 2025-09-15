import React from 'react';
import Link from 'next/link';
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

    const handleFollowRequest = async (action, followerId) => {
        try {
            const response = await fetch(`http://localhost:8080/api/follow/${action}` , {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ followerId })
            });

            if (response.ok) {
                onRead(notification.id);
            } else {
                const errorData = await response.json();
                console.error(`Failed to ${action} request: ${errorData.message || 'Unknown error'}`);
            }
        } catch (err) {
            console.error(`Error ${action} follow request:`, err);
        }
    };

    const handleGroupInvite = async (action, invitationID) => {
        try {
            const response = await fetch(`http://localhost:8080/api/groups/invites/${invitationID}/${action}` , {
                method: 'POST',
                credentials: 'include',
            });

            if (response.ok) {
                onRead(notification.id);
            } else {
                const errorData = await response.json();
                console.error(`Failed to ${action} group invite: ${errorData.message || 'Unknown error'}`);
            }
        } catch (err) {
            console.error(`Error ${action} group invite:`, err);
        }
    };

    const handleFollowBack = async (userId) => {
        try {
            const response = await fetch('http://localhost:8080/api/users/follow', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ userId })
            });

            if (response.ok) {
                onRead(notification.id);
            } else {
                const errorData = await response.json();
                console.error(`Failed to follow back: ${errorData.message || 'Unknown error'}`);
            }
        } catch (err) {
            console.error('Error following back:', err);
        }
    };

    const renderActionButtons = () => {
        switch (notification.type) {
            case 'new_post':
            case 'new_comment':
            case 'new_reaction':
                return (
                    <Link href={`/post/${notification.post_id.String}`}>
                        <a className="text-sm text-blue-600 hover:underline">View Post</a>
                    </Link>
                );
            case 'follow_request':
                return (
                    <div className="flex space-x-2 mt-2">
                        <button onClick={() => handleFollowRequest('accept', notification.actor_id)} className="px-3 py-1 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition text-xs">Accept</button>
                        <button onClick={() => handleFollowRequest('decline', notification.actor_id)} className="px-3 py-1 bg-gray-200 text-gray-800 rounded-full hover:bg-gray-300 transition text-xs">Decline</button>
                    </div>
                );
            case 'group_invite':
                return (
                    <div className="flex space-x-2 mt-2">
                        <button onClick={() => handleGroupInvite('accept', notification.id)} className="px-3 py-1 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition text-xs">Accept</button>
                        <button onClick={() => handleGroupInvite('decline', notification.id)} className="px-3 py-1 bg-gray-200 text-gray-800 rounded-full hover:bg-gray-300 transition text-xs">Decline</button>
                    </div>
                );
            case 'new_follower':
                return (
                    <div className="flex space-x-2 mt-2">
                        <button onClick={() => handleFollowBack(notification.actor_id)} className="px-3 py-1 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition text-xs">Follow Back</button>
                    </div>
                );
            default:
                return null;
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
                {renderActionButtons()}
            </div>
        </div>
    );
};

export default NotificationItem;
