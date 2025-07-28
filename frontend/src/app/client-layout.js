'use client'

import { useEffect, useState } from "react";
import Sidebar from "../components/sidebar";

export default function ClientLayout({ children }) {
  const [data, setData] = useState(null);

  useEffect(() => {
    async function fetchUser() {
      try {
        const res = await fetch("http://localhost:8080/api/profile/current", {
          method: 'GET',
          credentials: 'include'
        })
        if (res.ok) {
          const data = await res.json()
          // console.log(data)
          setData(data)
        }
      } catch (error) {
        console.error("error fetching user: ", error)
      }
    }

    fetchUser()
  }, [])

  if (!data) return <div>Loading...</div>

  return (
    <div className="h-screen px-[1%] md:px-[12%] 2xl:px-[18%]">
      <div className="flex min-h-screen">
        <Sidebar data={data} />
        <main className="flex-1">
          {children}
        </main>
      </div>
    </div>
  );
}