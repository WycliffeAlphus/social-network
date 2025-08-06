"use client";

export default function PostCard({ post, showComments, setShowComments }) {
  return (
    <div className="max-w-xl mx-auto p-4 bg-[var(--background)] border border-[var(--foreground)] rounded-2xl shadow mb-6">
      <div className="flex items-center mb-3">
        <img
          src={post.creatorimg || "profile.jpg"}
          alt="User Avatar"
          className="w-10 h-10 rounded-full mr-3"
        />
        <div>
          <h4 className="text-sm font-semibold">{post.creator || "Jane Doe"}</h4>
          <p className="text-xs">{post.createdat}</p>
        </div>
      </div>

      <p className="text-sm mb-3">{post.content}</p>

      <div className="flex justify-between text-sm border-t pt-2">
        <button className="hover:text-blue-600">👍</button>
        <button className="hover:text-blue-600">👎</button>
        <button
          onClick={() => setShowComments((prev) => !prev)}
          className="hover:text-blue-600"
        >
          💬
        </button>
      </div>

      {showComments && (
        <div className="mt-4">
          <textarea
            className="w-full p-2 border border-gray-300 rounded-md mb-2 text-sm"
            rows="2"
            placeholder="Write a comment..."
          ></textarea>
          <button className="text-sm bg-blue-500 text-white px-3 py-1 rounded hover:bg-blue-600 mb-4">
            Post
          </button>

          {/* Static comments for now — replace with real data later */}
          <div className="space-y-2">
            <div className="flex items-start gap-2">
              <img src="profile.jpg" className="w-8 h-8 rounded-full" />
              <div className="bg-[var(--background)] p-2 rounded-lg text-sm">
                <strong>Alex:</strong> Sounds amazing! Was it about AI in rural schools too?
              </div>
            </div>
            <div className="flex items-start gap-2">
              <img src="profile.jpg" className="w-8 h-8 rounded-full" />
              <div className="bg-[var(--background)] p-2 rounded-lg text-sm">
                <strong>Maria:</strong> I was there too! The part on teacher automation was 🔥.
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
