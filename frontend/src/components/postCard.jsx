"use client";

import CommentSection from "./commentCard";

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
        <CommentSection postId={post.id} />
      )}
    </div>
  );
}
