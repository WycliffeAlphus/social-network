"use client";

import Link from "next/link";
import {
  HomeIcon, BellIcon, ChatBubbleLeftIcon,
  UserIcon, UserGroupIcon, PlusIcon,
  UserCircleIcon,
  UserPlusIcon
} from '@heroicons/react/24/outline';
import { useEffect, useRef, useState } from "react";

export default function Sidebar({ data }) {
  console.log(data)
  const [showCreateOptions, setShowCreateOptions] = useState(false);
  const [showAccountOptions, setShowAccountOptions] = useState(false);
  const createOptionsRef = useRef(null);
  const accountsRef = useRef(null);

  useEffect(() => {
    function handleClickOutside(event) {
      if (createOptionsRef.current && !createOptionsRef.current.contains(event.target)) {
        setShowCreateOptions(false)
      }

      if (accountsRef.current && !accountsRef.current.contains(event.target)) {
        setShowAccountOptions(false)
      }
    }

    document.addEventListener("click", handleClickOutside)
    return () => {
      document.removeEventListener("click", handleClickOutside)
    }
  }, [showCreateOptions, showAccountOptions])

  return (
    <aside className="sticky w-[5rem] xl:w-[12rem] pt-[3rem] top-0 text-white h-[100vh] overflow-y-auto">
      <ul className="w-full">
        <li className="p-4 border-gray-100/20">
          <Link href="/" className="flex w-full pb-4 items-center items-center justify-center md:justify-start">
            <HomeIcon className="h-6 w-6 mr-2" />
            <span className="hidden xl:inline">Home</span>
          </Link>
        </li>
        <li className="p-4 border-gray-100/20">
          <Link href="/profile" className="flex w-full pb-4 items-center items-center justify-center md:justify-start">
            <UserPlusIcon className="h-6 w-6 mr-2" />
            <span className="hidden xl:inline">Follow requests</span>
          </Link>
        </li>
        <li className="p-4 border-gray-100/20">
          <Link href="/messages" className="flex w-full pb-4 items-center items-center justify-center md:justify-start">
            <ChatBubbleLeftIcon className="h-6 w-6 mr-2" />
            <span className="hidden xl:inline">Messages</span>
          </Link>
        </li>
        <li className="p-4 border-gray-100/20">
          <Link href="/groups" className="flex w-full pb-4 items-center items-center justify-center md:justify-start">
            <UserGroupIcon className="h-6 w-6 mr-2" />
            <span className="hidden xl:inline">Groups</span>
          </Link>
        </li>
        <li className="p-4 border-gray-100/20">
          <Link href="/notifications" className="flex w-full pb-4 items-center items-center justify-center md:justify-start">
            <BellIcon className="h-6 w-6 mr-2" />
            <span className="hidden xl:inline">Notifications</span>
          </Link>
        </li>
        <li className="p-4 border-gray-100/20">
          <Link href="/profile" className="flex w-full pb-4 items-center items-center justify-center md:justify-start">
            <UserIcon className="h-6 w-6 mr-2" />
            <span className="hidden xl:inline">Profile</span>
          </Link>
        </li>
        <li className="relative p-4 border-gray-100/20">
          <button
            onClick={() => {
              setShowCreateOptions((prev) => !prev)
              setShowAccountOptions(false)
            }}
            className="cursor-pointer flex w-full pb-4 items-center justify-center md:justify-start"
          >
            <div className="bg-white rounded-full mr-2">
              <PlusIcon className="text-black h-6 w-6" />
            </div>
            <span className="hidden xl:inline">Create</span>
          </button>

          {showCreateOptions && (
            <div ref={createOptionsRef} className="fixed flex flex-col gap-2 shadow-[0_0_12px_rgba(65,105,225)] bg-black rounded-lg p-6">
              <Link href="/create/post" className="text-sm hover:text-[#4169e1] text-white">
                + New Post
              </Link>
              <Link href="/create/group" className="text-sm hover:text-[#4169e1] text-white">
                + New Group
              </Link>
            </div>
          )}
        </li>

        <li className="p-4 border-gray-100/20">
          <button
            onClick={() => {
              setShowAccountOptions((prev) => !prev)
              setShowCreateOptions(false);
            }}
            className="cursor-pointer flex w-full mt-33 h-fit items-center justify-center md:justify-start"
          >
            {data?.profile.img_url ? (
              <img
                src={data.profile.img_url}
                alt="User avatar"
                className="h-6 w-6 mr-2 rounded-full object-cover"
              />
            ) : (
              <UserCircleIcon className="h-7 w-7 mr-2" />
            )}
            <span className="hidden xl:inline">Account</span>
          </button>

          {showAccountOptions && (
            <div ref={accountsRef} className="fixed flex flex-col gap-2 shadow-[0_0_12px_rgba(65,105,225)] bg-black rounded-lg p-6">
              <p className="cursor-pointer text-sm hover:text-red-700 text-red-600">
                Log out
              </p>
            </div>
          )}
        </li>
      </ul>
    </aside>
  );
}
