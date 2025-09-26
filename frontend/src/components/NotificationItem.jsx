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

const NotificationItem = ({ notification, onRead, onFollowBack, followStatus, groupInviteStatus }) => {
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
                if (action === 'accept') {
                    const data = await response.json();
                    onFollowBack(notification.id, notification.actor_id, data.status);
                } else {
                    onFollowBack(notification.id, notification.actor_id, 'not_following');
                }
                handleMarkAsRead();
            } else {
                const errorData = await response.json();
                console.error(`Failed to ${action} request: ${errorData.message || 'Unknown error'}`);

                // If the request was already processed, mark as read anyway
                if (response.status === 404 && errorData.message && errorData.message.includes('already processed')) {
                    handleMarkAsRead();
                    if (action === 'accept') {
                        onFollowBack(notification.id, notification.actor_id, 'accepted');
                    } else {
                        onFollowBack(notification.id, notification.actor_id, 'not_following');
                    }
                }
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
                handleMarkAsRead();
            } else {
                const errorData = await response.json();
                console.error(`Failed to ${action} group invite: ${errorData.message || 'Unknown error'}`);

                // If the invite was already processed, mark as read anyway
                if (response.status === 404 || (errorData.message && errorData.message.includes('already processed'))) {
                    handleMarkAsRead();
                }
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
                body: JSON.stringify({ userId: userId, isFollowBack: true })
            });

            if (response.ok) {
                const data = await response.json();
                onFollowBack(notification.id, notification.actor_id, data.status);
                handleMarkAsRead();
            } else {
                const errorData = await response.json();
                console.error(`Failed to follow back: ${errorData.message || 'Unknown error'}`);
            }
        } catch (err) {
            console.error('Error following back:', err);
        }
    };

    const getNotificationLink = () => {
        switch (notification.type) {
            case 'new_post':
            case 'new_comment':
            case 'new_reaction':
                return `/?post_id=${notification.post_id.String}`;
            case 'follow_request':
            case 'new_follower':
            case 'follow_back':
            case 'follow_accepted':
                return `/profile/${notification.actor_id}`;
            case 'group_invite':
            case 'group_join_request':
            case 'group_join_accepted':
            case 'group_event_created':
                return `/groups/${notification.content_id}`;
            default:
                return '#';
        }
    };

    const renderActionButtons = () => {
        // Don't show action buttons for read notifications
        if (notification.is_read) {
            return null;
        }

        switch (notification.type) {
            case 'follow_request':
                if (followStatus === 'not_following' || followStatus === undefined) {
                    return (
                        <div className="flex space-x-2 mt-2">
                            <button onClick={(e) => {e.preventDefault(); handleFollowRequest('accept', notification.actor_id)}} className="px-3 py-1 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition text-xs">Accept</button>
                            <button onClick={(e) => {e.preventDefault(); handleFollowRequest('decline', notification.actor_id)}} className="px-3 py-1 bg-gray-200 text-gray-800 rounded-full hover:bg-gray-300 transition text-xs">Decline</button>
                        </div>
                    );
                } else if (followStatus === 'following' || followStatus === 'requested') {
                    return (
                        <div className="flex space-x-2 mt-2">
                            <button className="px-3 py-1 bg-gray-200 text-gray-800 rounded-full transition text-xs" disabled>Following</button>
                        </div>
                    )
                }
                return null;
            case 'group_invite':
                if (groupInviteStatus === 'pending' || groupInviteStatus === undefined) {
                    return (
                        <div className="flex space-x-2 mt-2">
                            <button onClick={(e) => {e.preventDefault(); handleGroupInvite('accept', notification.content_id)}} className="px-3 py-1 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition text-xs">Accept</button>
                            <button onClick={(e) => {e.preventDefault(); handleGroupInvite('decline', notification.content_id)}} className="px-3 py-1 bg-gray-200 text-gray-800 rounded-full hover:bg-gray-300 transition text-xs">Decline</button>
                        </div>
                    );
                }
                return null;
            case 'new_follower':
                if (followStatus === 'not_following' || followStatus === undefined) {
                    return (
                        <div className="flex space-x-2 mt-2">
                            <button onClick={(e) => {e.preventDefault(); handleFollowBack(notification.actor_id)}} className="px-3 py-1 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition text-xs">Follow Back</button>
                        </div>
                    );
                } else if (followStatus === 'following' || followStatus === 'requested') {
                    return (
                        <div className="flex space-x-2 mt-2">
                            <button className="px-3 py-1 bg-gray-200 text-gray-800 rounded-full transition text-xs" disabled>Following</button>
                        </div>
                    )
                }
                return null;
            default:
                return null;
        }
    };

    return (
        <Link href={getNotificationLink()} className={itemClasses} onClick={handleMarkAsRead}>
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
        </Link>
    );
};

export default NotificationItem;

