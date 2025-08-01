export default function ProfileDetails({ user }) {
  return (
    <div className="bg-white rounded-lg shadow p-6 mb-8">
      <h2 className="text-xl font-semibold mb-4">About</h2>
      {user.about ? (
        <p className="text-gray-700 whitespace-pre-line">{user.about}</p>
      ) : (
        <p className="text-gray-400 italic">No bio yet</p>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-6">
        <div>
          <h3 className="text-sm font-medium text-gray-500">Date of Birth</h3>
          <p>{user.dob || 'Not specified'}</p>
        </div>
        <div>
          <h3 className="text-sm font-medium text-gray-500">Member Since</h3>
          <p>{new Date(user.created_at).toLocaleDateString()}</p>
        </div>
      </div>
    </div>
  )
}
