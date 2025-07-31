"use client";

import Link from "next/link";
import {
  HomeIcon, BellIcon, ChatBubbleLeftIcon,
  UserIcon, UserGroupIcon, PlusIcon,
  UserCircleIcon,
  UserPlusIcon
} from '@heroicons/react/24/outline';
import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import CreatePost from "./createpost";

export default function Sidebar({ data }) {
  const [showCreateOptions, setShowCreateOptions] = useState(false);
  const [showAccountOptions, setShowAccountOptions] = useState(false);
  const [showCreatePosts, setShowCreatePosts] = useState(false);
  const createOptionsRef = useRef(null);
  const accountsRef = useRef(null);
  const router = useRouter();

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

  const handleLogout = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/logout", {
        method: "POST",
        credentials: "include",
      });
      if (response.ok) {
        router.push("/login");
      } else {
        alert("Logout failed.");
      }
    } catch (err) {
      alert("Logout error: " + err.message);
    }
  };

  useEffect(() => {
    if (showCreatePosts) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'auto';
    }

    // Cleanup function to reset overflow when component unmounts
    return () => {
      document.body.style.overflow = 'auto';
    };
  }, [showCreatePosts]);

  return (
    <>
      <aside className="sticky w-[5rem] xl:w-[12rem] pt-[3rem] top-0 h-[100vh] overflow-y-auto">
        <ul className="w-full">
          <li className="p-4 border-gray-100/20">
            <Link href="/" className="flex w-full pb-4 items-center items-center justify-center md:justify-start">
              <HomeIcon className="h-6 w-6 mr-2" />
              <span className="hidden xl:inline">Home</span>
            </Link>
          </li>
          <li className="p-4 border-gray-100/20">
            <Link href="/follow-requests" className="flex w-full pb-4 items-center items-center justify-center md:justify-start">
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
            <Link href={`/profile/${data.current_user_id}`} className="flex w-full pb-4 items-center items-center justify-center md:justify-start">
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
              <div ref={createOptionsRef} className="fixed flex flex-col gap-2 shadow-[0_0_12px_rgba(0,0,0,0.5),0_0_12px_rgba(255,255,255,0.5)] dark:bg-black bg-white rounded-lg p-6">
                <p
                  onClick={() => {
                    setShowCreatePosts(true)
                    setShowCreateOptions(false)
                  }}
                  className="cursor-pointer text-sm hover:text-blue-500"
                >
                  + New Post
                </p>
                <p className="cursor-pointer text-sm hover:text-blue-500">
                  + New Group
                </p>
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
              <div ref={accountsRef} className="fixed shadow-[0_0_12px_rgba(0,0,0,0.5),0_0_12px_rgba(255,255,255,0.5)] dark:bg-black bg-white rounded-lg p-6">
                <p
                  className="cursor-pointer text-sm hover:text-red-700 text-red-600"
                  onClick={handleLogout}
                >
                  Log out
                </p>
              </div>
            )}
          </li>
        </ul>
      </aside>

      {showCreatePosts && (
        <div
          className="fixed flex justify-center inset-0 p-7 z-50 overflow-auto bg-blue-300/20 backdrop-blur-xs"
          onClick={() => setShowCreatePosts(false)}
        >
          <CreatePost onClose={() => setShowCreatePosts(false)} />
        </div>
      )}
    </>
  );
}
