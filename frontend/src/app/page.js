"use client";

import { useEffect, useState } from "react";
import FollowSuggestion from "../components/followsuggestions";
import Rightbar from "../components/rightbar";
import PostCard from "../components/postCard"; // Reusable post component

export default function Home() {
  const [showComments, setShowComments] = useState(false);
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchFeeds = async () => {
      setLoading(true);
      try {
        const response = await fetch("http://localhost:8080/api/feeds", {
          method: "GET",
          credentials: "include",
        });

        if (!response.ok) throw new Error("Failed to fetch feeds");
        const data = await response.json();
        setPosts(data.data.posts || []); // fallback to empty array
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchFeeds();
  }, []);

 
  return (
    <div className="flex min-h-screen">
      <main className="flex-1 border-x mr-[20px] border-gray-400">
        <div className="lg:hidden">
          <FollowSuggestion />
        </div>

        <div className="p-4 border-t lg:border-0">
          {loading && <p className="text-sm text-gray-500">Loading feeds...</p>}
          {error && <p className="text-sm text-red-500">{error}</p>}
          {!loading && posts.length === 0 && <p>No posts available.</p>}

          {!loading &&
            posts.map((post) => (
              <PostCard
                key={post.id}
                post={post}
                showComments={showComments}
                setShowComments={setShowComments}
              />
            ))}
        </div>
      </main>

      <Rightbar />
    </div>
  );
}
