"use client"

import { useUser } from "@/context/user-context";
import { useParams } from "next/navigation";
import { useRouter } from 'next/navigation';
import { useEffect } from "react";

export default function Messages() {
    const currentUserId = useUser()
    const { id } = useParams()
    const router = useRouter()

    console.log(currentUserId)

    useEffect(() => {
        if (currentUserId && currentUserId === id) {
            router.push('/messages')
            return
        }

        // fetchContacts();
    }, [id, currentUserId, router]);

    return (
        <div>page</div>
    )
}