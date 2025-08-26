'use client';

import { useEffect, useState } from 'react';
import NotificationsList from '../../components/NotificationsList';
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
                setNotifications(data);
            } catch (err) {
                setError('Failed to fetch notifications.');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchNotifications();
    }, []);

    if (loading) {
        return <div className="text-center mt-8">Loading notifications...</div>;
    }

    if (error) {
        return <div className="text-center mt-8 text-red-500">{error}</div>;
    }

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-2xl font-bold mb-4">Notifications</h1>
            {notifications.length > 0 ? (
                <NotificationsList notifications={notifications} />
            ) : (
                <p className="text-center text-gray-500">No notifications yet.</p>
            )}
        </div>
    );
}

export default NotificationsPage;
