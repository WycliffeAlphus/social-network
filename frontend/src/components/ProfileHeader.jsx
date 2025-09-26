'use client'
import Link from 'next/link'
import { EnvelopeIcon } from "@heroicons/react/24/outline";

export default function ProfileHeader({
  user,
  isOwner,
  isPublic,
  followersCount,
  followingCount,
  followStatus,
  incomingFollowRequestStatus,
  onFollow,
  onCancelRequest,
  onAcceptFollowRequest,
  onDeclineFollowRequest,
  onToggleVisibility
}) {
  const renderFollowButton = () => {
    let label = 'Follow'
    let className = 'bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-xl font-semibold transition-all duration-300 transform hover:scale-105 shadow-lg hover:shadow-xl flex items-center gap-2'
    let onClick = onFollow
    let disabled = false

    if (followStatus === 'accepted') {
      label = 'Following'
      className = 'bg-white border-2 border-blue-200 text-blue-700 px-6 py-3 rounded-xl font-semibold transition-all duration-300 hover:bg-blue-50 hover:border-blue-300 shadow-md hover:shadow-lg flex items-center gap-2'
      onClick = null
      disabled = true
    } else if (followStatus === 'requested') {
      label = 'Requested'
      className = 'bg-blue-100 border-2 border-blue-200 text-blue-600 px-6 py-3 rounded-xl font-semibold transition-all duration-300 hover:bg-blue-50 cursor-pointer shadow-md hover:shadow-lg flex items-center gap-2'
      onClick = onCancelRequest
    }
    
    if (incomingFollowRequestStatus === 'pending') {
      return (
        <div className="flex space-x-4">
          <button onClick={onAcceptFollowRequest} className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
            Accept
          </button>
          <button onClick={onDeclineFollowRequest} className="px-4 py-2 bg-gray-300 text-gray-800 rounded hover:bg-gray-400">
            Decline
          </button>
        </div>
      )
    }

    return (
      <button onClick={onClick} className={className} disabled={disabled}>
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          {followStatus === 'accepted' ? (
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          ) : (
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          )}
        </svg>
        {label}
      </button>
    )
  }

  return (
    <div className="bg-white rounded-2xl shadow-2xl border border-blue-100 p-8 mb-8">
      <div className="flex flex-col lg:flex-row items-start lg:items-center gap-8">
        {/* Profile Picture */}
        <div className="relative group">
          <div className="w-36 h-36 rounded-full overflow-hidden bg-blue-200 shadow-2xl ring-4 ring-white ring-offset-4 ring-offset-blue-50 transition-all duration-300 group-hover:ring-offset-8">
            {user.img_url ? (
              <img
                src={user.img_url}
                alt="avatar"
                className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-110"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-5xl font-bold text-white bg-blue-500">
                {user.first_name[0]}{user.last_name[0]}
              </div>
            )}
          </div>
          <div className="absolute -bottom-2 -right-2 w-8 h-8 bg-green-500 rounded-full border-4 border-white shadow-lg"></div>
        </div>

        {/* User Info */}
        <div className="flex-1 space-y-6">
          {/* Name and Email */}
          <div>
            <h1 className="text-4xl font-bold text-blue-800 mb-2">
              {user.first_name} {user.last_name}
            </h1>
            <div className="flex items-center gap-2 mb-4">
              <svg className="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" />
              </svg>
              <p className="text-blue-600 font-medium text-lg">{user.email}</p>
            </div>
          </div>

          {/* Followers/Following Stats */}
          <div className="flex gap-8">
            <Link
              href={`/followers/${user.id}`}
              className="group bg-white rounded-xl px-6 py-4 border border-blue-100 hover:border-blue-300 transition-all duration-300 hover:shadow-lg transform hover:-translate-y-1"
            >
              <div className="text-center">
                <div className="text-2xl font-bold text-blue-700 group-hover:text-blue-800 transition-colors">
                  {followersCount}
                </div>
                <div className="text-blue-600 font-medium text-sm uppercase tracking-wide">
                  Followers
                </div>
              </div>
            </Link>

            <Link
              href={`/following/${user.id}`}
              className="group bg-white rounded-xl px-6 py-4 border border-blue-100 hover:border-blue-300 transition-all duration-300 hover:shadow-lg transform hover:-translate-y-1"
            >
              <div className="text-center">
                <div className="text-2xl font-bold text-blue-700 group-hover:text-blue-800 transition-colors">
                  {followingCount}
                </div>
                <div className="text-blue-600 font-medium text-sm uppercase tracking-wide">
                  Following
                </div>
              </div>
            </Link>
          </div>

          {!isOwner &&
            <Link href={`/messages/${user.id}`}>
              <EnvelopeIcon className="h-6 w-6 text-gray-500" />
            </Link>
          }

          {/* Actions */}
          <div className="pt-4">
            {isOwner ? (
              <div className="bg-white rounded-xl p-6 border border-blue-100 shadow-sm">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-blue-500 rounded-lg flex items-center justify-center">
                      <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                    </div>
                    <div>
                      <div className="font-semibold text-gray-800">Profile Visibility</div>
                      <div className="text-sm text-gray-600">
                        {isPublic ? 'Anyone can see your profile' : 'Only followers can see your profile'}
                      </div>
                    </div>
                  </div>

                  <label className="relative inline-flex items-center cursor-pointer">
                    <input
                      type="checkbox"
                      className="sr-only peer"
                      checked={isPublic}
                      onChange={onToggleVisibility}
                    />
                    <div className="w-14 h-7 bg-gray-200 rounded-full peer peer-checked:bg-blue-500 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all after:shadow-md peer-checked:after:translate-x-7 peer-checked:after:border-white hover:shadow-lg transition-shadow" />
                  </label>
                </div>
              </div>
            ) : (
              <div className="flex justify-start">
                {renderFollowButton()}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}