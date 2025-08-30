"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import CreatePost from "@/components/createpost";
import PostCard from "@/components/postCard"; // Import PostCard

export default function GroupPage() {
  const [group, setGroup] = useState(null);
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showCreatePost, setShowCreatePost] = useState(false);
  const params = useParams();
  const { id } = params;

  useEffect(() => {
    if (id) {
      fetchGroupDetails();
      fetchGroupPosts();
    }
  }, [id]);

  const fetchGroupDetails = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/groups/${id}`, {
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Failed to fetch group details");
      }
      const data = await response.json();
      setGroup(data);
    } catch (e) {
      setError(e.message);
    }
  };

  const fetchGroupPosts = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/groups/${id}/posts`, {
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Failed to fetch group posts");
      }
      const data = await response.json();
      setPosts(data || []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-4">
      {loading && <p>Loading...</p>}
      {error && <p>Error: {error}</p>}
      {group && (
        <div>
          <h1 className="text-3xl font-bold">{group.title}</h1>
          <p>{group.description}</p>
          <button
            onClick={() => setShowCreatePost(true)}
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-300 ease-in-out my-4"
          >
            Create Post
          </button>
        </div>
      )}
      {showCreatePost && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
          <CreatePost onClose={() => setShowCreatePost(false)} groupId={id} />
        </div>
      )}
      <div className="mt-8">
        {posts.map((post) => (
          <PostCard
            key={post.id}
            post={post}
          />
        ))}
      </div>
    </div>
  );
}
