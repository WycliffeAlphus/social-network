'use client';

import { useEffect, useState } from 'react';
import NotificationItem from '../../components/NotificationItem';
import { getNotifications } from '../../lib/notifications';

function NotificationsPage() {
    const [notifications, setNotifications] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchNotifications = async () => {
            try {
                setLoading(true);
                const data = await getNotifications();
                setNotifications(data || []); // Ensure notifications is an array
            } catch (err) {
                setError('Failed to fetch notifications.');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchNotifications();
    }, []);

    const handleRead = (notificationId) => {
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
                            <NotificationItem key={notification.id} notification={notification} onRead={handleRead} />
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
