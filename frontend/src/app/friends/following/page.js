"use client";

import Navbar from "../../components/navbar"; 
import Sidebar from "../../components/sidebar"; 
import Link from "next/link";
import { usePathname } from 'next/navigation';

// Import the content components
import FollowingContent from '../../components/FollowingContent';

export default function FollowingPage() {
  const pathname = usePathname();

  // Determine the active tab based on the current URL
  const activeTab = pathname.startsWith('/friends/followers')
    ? 'followers'
    : pathname.startsWith('/friends/following')
      ? 'following'
      : 'requests';

  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      <div className="flex flex-1">
        <Sidebar />
        <main className="flex-1 p-6 bg-gray-100 text-gray-800">
          <h1 className="text-3xl font-bold mb-6 text-blue-800">My Connections</h1>

          {/* Sub-navigation for Friends sections */}
          <nav className="mb-6 border-b border-gray-300">
            <ul className="flex space-x-6">
              <li>
                <Link
                  href="/friends/requests"
                  className={`pb-2 block ${
                    activeTab === 'requests'
                      ? 'border-b-2 border-blue-600 text-blue-600 font-semibold'
                      : 'text-gray-600 hover:text-blue-500'
                  }`}
                >
                  Requests
                </Link>
              </li>
              <li>
                <Link
                  href="/friends/followers"
                  className={`pb-2 block ${
                    activeTab === 'followers' ? 'border-b-2 border-blue-600 text-blue-600 font-semibold' : 'text-gray-600 hover:text-blue-500'
                  }`}
                >
                  Followers
                </Link>
              </li>
              <li>
                <Link
                  href="/friends/following"
                  className={`pb-2 block ${
                    activeTab === 'following' ? 'border-b-2 border-blue-600 text-blue-600 font-semibold' : 'text-gray-600 hover:text-blue-500'
                  }`}
                >
                  Following
                </Link>
              </li>
            </ul>
          </nav>

          {/* Content Area */}
          <FollowingContent userId="me" />
        </main>
      </div>
    </div>
  );
}
