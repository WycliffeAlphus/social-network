"use client";

import { useState } from "react";
import CommentSection from "./commentCard";
import { formatTimeAgo } from "./dateFormat";
import Reaction from "./reaction";

export default function PostCard({ post }) { // Removed showComments, setShowComments from props
  // Initialize with existing counts if available, otherwise 0
  const [likes, setLikes] = useState(post.likeCount || 0);
  const [dislikes, setDislikes] = useState(post.dislikeCount || 0);
  const [userReaction, setUserReaction] = useState(post.userReaction || "");
  const [isLoading, setIsLoading] = useState(false);
  const [commentCount, setCommentCount] = useState(post.commentcount || post.commentCount || 0);
  const [showComments, setShowComments] = useState(false); // Internal state for comments

  const handleLikeSubmit = async (e) => {
    e.preventDefault();
    if (isLoading) return; // Prevent multiple clicks

    setIsLoading(true);
    try {
      const response = await Reaction("like", post.id);
      console.log("Like response:", response);

      if (response.success) {
        setLikes(response.likeCount);
        setDislikes(response.dislikeCount);
        setUserReaction(response.userReaction);
      } else {
        console.error("Failed to update reaction:", response.message);
      }
    } catch (error) {
      console.error("Error submitting like:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDislikeSubmit = async (e) => {
    e.preventDefault();
    if (isLoading) return; // Prevent multiple clicks

    setIsLoading(true);
    try {
      const response = await Reaction("dislike", post.id);
      console.log("Dislike response:", response);

      if (response.success) {
        setLikes(response.likeCount);
        setDislikes(response.dislikeCount);
        setUserReaction(response.userReaction);
      } else {
        console.error("Failed to update reaction:", response.message);
      }
    } catch (error) {
      console.error("Error submitting dislike:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-xl mx-auto p-4 bg-[var(--background)] border border-[var(--foreground)] rounded-2xl shadow mb-6">
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center">
          <img
            src={post.creatorimg || "profile.jpg"}
            alt="User Avatar"
            className="w-10 h-10 rounded-full mr-3"
          />
          <div>
            <h4 className="text-sm font-semibold">{post.creator || "Jane Doe"}</h4>
          </div>
        </div>
        <p className="text-xs text-gray-500">{formatTimeAgo(post.createdat)}</p>
      </div>
      
      <p className="text-sm mb-3 break-words">{post.content}</p>
      
      {post.imageurl?.Valid && (() => {
        const raw = (post.imageurl.String || "");
        const stripped = raw.replace(/^\/?frontend\/public/, "");
        const src = stripped.startsWith("/") ? stripped : `/${stripped}`;
        return (
          <img
            src={src}
            alt="Post image"
            className="w-full rounded mt-2"
          />
        );
      })()}

      <div className="flex justify-between text-sm border-t pt-2">
        <button
          className={`hover:text-blue-600 ${userReaction === 'like' ? 'text-blue-600' : ''} ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}`}
          onClick={handleLikeSubmit}
          disabled={isLoading}
        >
          👍{likes}
        </button>

        <button
          className={`hover:text-blue-600 ${userReaction === 'dislike' ? 'text-red-600' : ''} ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}`}
          onClick={handleDislikeSubmit}
          disabled={isLoading}
        >
          👎{dislikes}
        </button>

        <button
          onClick={() => setShowComments((prev) => !prev)}
          className="hover:text-blue-600"
        >
          💬{commentCount}
        </button>
      </div>

      {showComments && (
        <CommentSection postId={post.id} onCountChange={setCommentCount} />
      )}
    </div>
  );
}
