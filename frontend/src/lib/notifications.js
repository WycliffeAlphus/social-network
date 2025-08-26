// This file will contain functions to interact with the notification API endpoints.

export async function getNotifications() {
  try {
    const response = await fetch('/api/notifications');
    if (!response.ok) {
      throw new Error('Failed to fetch notifications');
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching notifications:', error);
    return [];
  }
}

export async function markNotificationsAsRead() {
  try {
    const response = await fetch('/api/notifications/read', {
      method: 'POST',
    });
    if (!response.ok) {
      throw new Error('Failed to mark notifications as read');
    }
    return true;
  } catch (error) {
    console.error('Error marking notifications as read:', error);
    return false;
  }
}
