'use client'
import Link from 'next/link'

export default function ProfileHeader({
  user,
  isOwner,
  isPublic,
  followersCount,
  followingCount,
  followStatus,
  onFollow,
  onCancelRequest,
  onToggleVisibility
}) {
  const renderFollowButton = () => {
    let label = 'Follow'
    let className = 'bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded transition'
    let onClick = onFollow

    if (followStatus === 'accepted') {
      label = 'Following'
      className = 'border border-gray-400 text-gray-800 px-4 py-2 rounded transition'
    } else if (followStatus === 'requested') {
      label = 'Requested'
      className = 'bg-gray-300 text-gray-600 px-4 py-2 rounded transition'
      onClick = onCancelRequest
    }

    return (
      <button onClick={onClick} className={className}>
        {label}
      </button>
    )
  }

  return (
    <div className="flex flex-col md:flex-row items-start md:items-center gap-6 mb-8">
      {/* Profile Picture */}
      <div className="w-32 h-32 rounded-full overflow-hidden bg-gray-200 shadow">
        {user.img_url ? (
          <img src={user.img_url} alt="avatar" className="w-full h-full object-cover" />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-4xl text-gray-400">
            {user.first_name[0]}{user.last_name[0]}
          </div>
        )}
      </div>

      {/* Name, Email, Buttons */}
      <div className="flex-1">
        <h1 className="text-3xl font-bold mb-1">{user.first_name} {user.last_name}</h1>
        <p className="text-gray-600 mb-3">{user.email}</p>

        <div className="flex gap-6 mb-4">
          <Link href={`/followers/${user.id}`} className="text-blue-600 hover:underline">
            <strong>{followersCount}</strong> Followers
          </Link>
          <Link href={`/following/${user.id}`} className="text-blue-600 hover:underline">
            <strong>{followingCount}</strong> Following
          </Link>
        </div>

        {isOwner ? (
          <div className="flex items-center gap-2">
            <span className="text-sm text-gray-600">Set account to public</span>
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                className="sr-only peer"
                checked={isPublic}
                onChange={onToggleVisibility}
              />
              <div className="w-11 h-6 bg-gray-200 rounded-full peer peer-checked:bg-blue-500 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:after:translate-x-full" />
            </label>
          </div>
        ) : (
          renderFollowButton()
        )}
      </div>
    </div>
  )
}
