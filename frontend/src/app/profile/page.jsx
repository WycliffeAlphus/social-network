"use client";
import { useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function Profile() {
    const searchParams = useSearchParams();
    const id = searchParams.get("id");
    const [profile, setProfile] = useState(null);
    const [error, setError] = useState("");
    console.log(id);


    useEffect(() => {
        fetch(`http://localhost:8080/api/profile?id=${id ? `${id}` : ''}`, {
            credentials: 'include',
        })
            .then(res => {
                if (!res.ok) throw new Error("Not allowed");

                return res.json();
            })
            .then(data => setProfile(data.data))
            .catch(err => setError(err.message));
    }, []);

    if (error) return <p>{error}</p>;
    if (!profile) return <p>Loading...</p>;

   return (
    <div className="p-6">
        {profile.profile_visibility === "private" ? (
            <div>
                <h1>Private Profile</h1>
                <p>This profile is private. Only you can see the details.</p>
            </div>
        ) : (
            <div>
                <h1>{profile.first_name} {profile.last_name}</h1>
                <p>{profile.about}</p>
            </div>
        )}
    </div>
);
}