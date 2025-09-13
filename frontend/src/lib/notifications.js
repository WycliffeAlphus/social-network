const API_URL = 'http://localhost:8080';

export async function getNotifications() {
  try {
    const response = await fetch(`${API_URL}/api/notifications`, {
      credentials: 'include',
    });
    if (!response.ok) {
      const errorBody = await response.text();
      throw new Error(`Failed to fetch notifications: ${response.status} ${response.statusText} - ${errorBody}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching notifications:', error);
    return [];
  }
}

export async function markNotificationsAsRead() {
  try {
    const response = await fetch(`${API_URL}/api/notifications/read`, {
      method: 'PUT',
      credentials: 'include',
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

export async function markNotificationAsRead(id) {
  try {
    const response = await fetch(`${API_URL}/api/notifications/${id}/read`, {
      method: 'PUT',
      credentials: 'include',
    });
    if (!response.ok) {
      throw new Error('Failed to mark notification as read');
    }
    return true;
  } catch (error) {
    console.error('Error marking notification as read:', error);
    return false;
  }
}
