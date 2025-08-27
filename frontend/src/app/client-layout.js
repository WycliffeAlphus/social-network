'use client'

import { useEffect, useState } from "react";
import Sidebar from "../components/sidebar";
import { usePathname } from "next/navigation";
import connectWebsocket from "@/components/ws";

export default function ClientLayout({ children }) {
  const [data, setData] = useState(null);
  const pathname = usePathname();

  useEffect(() => {
    // don't fetch user if on auth pages
    if (pathname === '/login' || pathname === '/register') {
      return;
    }

    async function fetchUser() {
      try {
        const res = await fetch("http://localhost:8080/api/profile/currentuser", {
          method: 'GET',
          credentials: 'include'
        })
        if (res.ok) {
          const data = await res.json()
          // console.log(data)
          setData(data)
          connectWebsocket(data.current_user_id)
        }
      } catch (error) {
        console.error("error fetching user: ", error)
      }
    }

    fetchUser()
  }, [])

  // don't show layout for auth pages
  if (pathname === '/login' || pathname === '/register') {
    return <>{children}</>
  }

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