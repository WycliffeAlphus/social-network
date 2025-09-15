"use client";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";

function Group() {

    const params = useParams();
    const groupId = params?.id;

    const [groupdetails, setGroupDetails] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function GetGroupdetails() {
            try {
                const response = await fetch(`http://localhost:8080/api/groups/${groupId}/view`, {
                    credentials: 'include',
                });
                if (!response.ok) {
                    setError(`HTTP error! status: ${response.status}`);
                }
                const data = await response.json();
                setGroupDetails(data.data);
                setLoading(false);

            } catch (error) {
                setError(error.message);
                setLoading(false);
            }

        }

        return () => {
            GetGroupdetails()
        }
    }, [])

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error}</div>;
    }
    console.log(groupdetails.data);



    return (
        <main className="flex-1 px-4 mt-5 @container">
            <div className="flex flex-col items-center ">
                <div className="flex flex-col items-center justify-center">
                    <p className="text-white text-2xl font-bold leading-tight tracking-[-0.015em]">{groupdetails.title}</p>
                </div>

            </div>
            <div className="py-2">
                <h2 className=" text-xl font-bold leading-tight tracking-[-0.015em] mb-1">About</h2>
                <p className=" text-base font-normal leading-relaxed">
                    {groupdetails.description}
                </p>
            </div>
            <div className="flex flex-col gap-3 py-6 sm:flex-row">
                <button className="flex min-w-[84px] cursor-pointer items-center justify-center overflow-hidden rounded-xl h-12 px-5 bg-[#1919e6] text-[#f8f8fc] text-base font-bold leading-normal tracking-[0.015em] w-full transition-colors hover:bg-blue-800">
                    <span className="truncate">Add Event</span>
                </button>
                <button className="flex min-w-[84px] cursor-pointer items-center justify-center overflow-hidden rounded-xl h-12 px-5 bg-[#e7e7f3] text-[#0e0e1b] text-base font-bold leading-normal tracking-[0.015em] w-full transition-colors hover:bg-zinc-300">
                    <span className="truncate">View Events</span>
                </button>
                <button className="flex min-w-[84px] cursor-pointer items-center justify-center overflow-hidden rounded-xl h-12 px-5 bg-[#e7e7f3] text-[#0e0e1b] text-base font-bold leading-normal tracking-[0.015em] w-full transition-colors hover:bg-zinc-300">
                    <span className="truncate">Invite members</span>
                </button>
            </div>
        </main>
    )
}

export default Group