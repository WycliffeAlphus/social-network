"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import SearchBar from "./searchbar";
import {
  HomeIcon, BellIcon, ChatBubbleLeftIcon,
  UserIcon, UserGroupIcon, PlusIcon
} from '@heroicons/react/24/solid';

export default function Sidebar() {
  const router = useRouter();
  const pathname = router.pathname;

  const isGroupsPage = pathname === "/groups";

  return (
    <aside className="w-1/5 bg-blue-900 text-white min-h-full">
      {isGroupsPage ? (
        <div>
          <h2 className="text-lg font-bold mb-4 text-center">Groups</h2>
          <ul className="space-y-2">
            <li>
              <SearchBar />
            </li>
            <li>
              <Link
                href="/create-group"
                className="block px-4 py-2 flex  text-white hover:bg-blue-700 hover:shadow-md transition duration-200"
              >
                <PlusIcon className="h-6 w-6 mr-1" />
                Create Group
              </Link>
            </li>
            <li>
              <Link
                href="/groups"
                className="block px-4 py-2 r  text-white hover:bg-blue-700 hover:shadow-md transition duration-200"
              >
                Your Groups
              </Link>
            </li>
          </ul>
        </div>
      ) : (
        <ul className="w-full">
          <li className="p-4 border-b border-gray-100/20">
            <Link href="/" className="flex w-full  pb-4 items-center">
              <HomeIcon className="h-6 w-6 mr-2" /> Home
            </Link> <SearchBar />
          </li>
          <li className="p-4 border-b border-gray-100/20">
            <Link href="/profile" className="flex w-full  pb-4 items-center">
              <UserIcon className="h-6 w-6 mr-2" /> Profile
            </Link>
          </li>
          <li className="p-4 border-b border-gray-100/20">
            <Link href="/messages" className="flex w-full  pb-4 items-center">
              <ChatBubbleLeftIcon className="h-6 w-6 mr-2" /> Messages
            </Link>
          </li>
          <li className="p-4 border-b border-gray-100/20">
            <Link href="/friends" className="flex w-full pb-4 items-center">
              <div className="relative w-6 h-6 mr-2">
                <UserGroupIcon className="h-5 w-5 text-white" />
                <li className="p-4 border-b border-gray-100/20">
                  <Link href="/friends" className="flex w-full pb-4 items-center">
                    <div className="relative w-6 h-6 mr-2">
                      <UserGroupIcon className="h-5 w-5 text-white" />
                    </div>
                    Friends
                  </Link>
                </li>
              </div>
              Friends
            </Link>
          </li>
          <li className="p-4 border-b border-gray-100/20">
            <Link href="/groups" className="flex w-full  pb-4 items-center">
              <div className="w-7 h-7 flex items-center justify-center rounded-full border-2 border-white mr-1">
                <UserGroupIcon className="h-5 w-5 text-white" />
              </div>
              Groups
            </Link>
          </li>
          <li className="p-4 border-b border-gray-100/20">
            <Link href="/notifications" className="flex w-full  pb-4 items-center">
              <BellIcon className="h-6 w-6 mr-2" /> Notifications
            </Link>
          </li>
        </ul>
      )}
    </aside>
  );
}