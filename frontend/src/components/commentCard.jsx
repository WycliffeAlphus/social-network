'use client';

import { useEffect, useState } from 'react';
import { formatTimeAgo } from './dateFormat';

export default function CommentSection({ postId }) {
    const [comments, setComments] = useState([]);
    const [newComment, setNewComment] = useState('');
    const [loading, setLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState(null);

    // Fetch comments on mount
    useEffect(() => {
        const fetchComments = async () => {
            setLoading(true);
            try {
                const res = await fetch(`http://localhost:8080/api/posts/${postId}/comments`, {
                    method: 'GET',
                    credentials: 'include',
                });
                const data = await res.json();
                setComments(data.comments || []);
            } catch (err) {
                console.error('Fetch error:', err);
                setError('Failed to load comments.');
            } finally {
                setLoading(false);
            }
        };

        fetchComments();
    }, [postId]);

    // Submit a new comment
    const handleSubmit = async (e) => {
        e.preventDefault();
        if (newComment.trim() === '') return;

        setSubmitting(true);
        try {
            const res = await fetch(`http://localhost:8080/api/posts/${postId}/comments`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ content: newComment }),
            });

            const data = await res.json();
            if (res.ok) {
                setComments((prev) => [...prev, data.comment]);
                setNewComment('');
            } else {
                setError(data.error || 'Comment failed.');
            }
        } catch (err) {
            console.error(err);
            setError('Error submitting comment.');
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <div className="mt-4">
            {loading && <p className="text-sm ">Loading comments...</p>}
            {error && <p className="text-sm text-red-500">{error}</p>}

            {!loading && comments.length > 0 && (
                <ul className="space-y-2 mt-2">
                    {comments.map((comment) => (
                        <li key={comment.id} className="border p-2 rounded">
                            <p className="text-sm font-semibold">
                                {comment.userNickname || `${comment.user_first_name} ${comment.user_last_name}`}
                            </p>
                            <p className="text-xs text-gray-500">{formatTimeAgo(comment.created_at)}</p>
                            <p className="text-sm ">{comment.content}</p>
                        </li>
                    ))}
                </ul>
            )}

            <form onSubmit={handleSubmit} className="mt-3 flex flex-col gap-2">
                <textarea
                    className="w-full border rounded p-2 text-sm"
                    rows={3}
                    placeholder="Write a comment..."
                    value={newComment}
                    onChange={(e) => setNewComment(e.target.value)}
                ></textarea>
                <button
                    type="submit"
                    disabled={submitting}
                    className="self-end bg-blue-600 hover:bg-blue-700 text-white text-sm px-3 py-1 rounded"
                >
                    {submitting ? 'Posting...' : 'Post Comment'}
                </button>
            </form>
        </div>
    );
}
