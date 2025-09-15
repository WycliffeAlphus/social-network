'use client';

import { useEffect, useState } from 'react';
import NotificationItem from '../../components/NotificationItem';
import { getNotifications, getFollowStatuses, getGroupInviteStatuses } from '../../lib/notifications';

function NotificationsPage() {
    const [notifications, setNotifications] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [followStatuses, setFollowStatuses] = useState({});
    const [groupInviteStatuses, setGroupInviteStatuses] = useState({});

    const fetchNotifications = async () => {
        try {
            setLoading(true);
            const data = await getNotifications();
            setNotifications(data || []);

            if (data && data.length > 0) {
                const userIds = data
                    .filter(n => n.type === 'new_follower' || n.type === 'follow_request')
                    .map(n => n.actor_id);
                
                if (userIds.length > 0) {
                    const statuses = await getFollowStatuses(userIds);
                    setFollowStatuses(statuses);
                }

                const invitationIds = data
                    .filter(n => n.type === 'group_invite')
                    .map(n => n.content_id);

                if (invitationIds.length > 0) {
                    const statuses = await getGroupInviteStatuses(invitationIds);
                    setGroupInviteStatuses(statuses);
                }
            }
        } catch (err) {
            setError('Failed to fetch notifications.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchNotifications();
    }, []);

    const handleRead = (notificationId) => {
        setNotifications(prev => prev.map(n => n.id === notificationId ? { ...n, is_read: true } : n));
    };

    const handleFollowBack = (notificationId, actorId, status) => {
        setFollowStatuses(prev => ({...prev, [actorId]: status}));
        setNotifications(prev => prev.map(n => n.id === notificationId ? { ...n, is_read: true } : n));
    };

    if (loading) {
        return <div className="text-center mt-8">Loading notifications...</div>;
    }

    if (error) {
        return <div className="text-center mt-8 text-red-500">{error}</div>;
    }

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-2xl font-bold mb-4 text-center">Notifications</h1>
            {notifications.length > 0 ? (
                <div className="max-w-md mx-auto">
                    {
                        notifications.map(notification => (
                            <NotificationItem 
                                key={notification.id} 
                                notification={notification} 
                                onRead={handleRead} 
                                onFollowBack={handleFollowBack}
                                followStatus={followStatuses[notification.actor_id]}
                                groupInviteStatus={groupInviteStatuses[notification.content_id]}
                            />
                        ))
                    }
                </div>
            ) : (
                <p className="text-center text-gray-500">No notifications yet.</p>
            )}
        </div>
    );
}

export default NotificationsPage;
