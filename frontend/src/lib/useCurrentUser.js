import { useState, useEffect } from 'react';

export function useCurrentUser() {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function fetchCurrentUser() {
            try {
                const response = await fetch('http://localhost:8080/api/profile/current', {
                    credentials: 'include',
                });
                if (!response.ok) {
                    throw new Error('Failed to fetch current user');
                }
                const userData = await response.json();
                setUser(userData);
            } catch (error) {
                console.error(error);
            } finally {
                setLoading(false);
            }
        }

        fetchCurrentUser();
    }, []);

    return { user, loading };
}