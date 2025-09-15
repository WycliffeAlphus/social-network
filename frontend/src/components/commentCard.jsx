'use client';

import { useEffect, useState } from 'react';
import { formatTimeAgo } from './dateFormat';

export default function CommentSection({ postId, onCountChange = () => {} }) {
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
                onCountChange(Array.isArray(data.comments) ? data.comments.length : 0);
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
                onCountChange((prev) => (typeof prev === 'number' ? prev + 1 : prev));
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

    // Add a reply to the correct parent in local state
    const addReplyToState = (parentId, reply) => {
        const addReplyRecursive = (items) => {
            return items.map((item) => {
                if (item.id === parentId) {
                    const existing = Array.isArray(item.replies) ? item.replies : [];
                    return { ...item, replies: [...existing, reply] };
                }
                if (Array.isArray(item.replies) && item.replies.length > 0) {
                    return { ...item, replies: addReplyRecursive(item.replies) };
                }
                return item;
            });
        };
        setComments((prev) => addReplyRecursive(prev));
        // Do not change top-level comment count when adding a reply
    };

    return (
        <div className="mt-4">
            {loading && <p className="text-sm ">Loading comments...</p>}
            {error && <p className="text-sm text-red-500">{error}</p>}

            {!loading && comments.length > 0 && (
                <ul className="space-y-2 mt-2">
                    {comments.map((comment) => (
                        <CommentItem
                            key={comment.id}
                            comment={comment}
                            postId={postId}
                            onAddReply={addReplyToState}
                            level={0}
                        />
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

function CommentItem({ comment, postId, onAddReply, level = 0 }) {
    const [showReply, setShowReply] = useState(false);
    const [replyText, setReplyText] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [expanded, setExpanded] = useState(false);

    const name =
        comment.userNickname ||
        comment.user_nickname ||
        `${comment.user_first_name || ''} ${comment.user_last_name || ''}`.trim();
    const createdAt = comment.created_at || comment.createdAt;
    const replyCount = Array.isArray(comment.replies) ? comment.replies.length : 0;

    const handleReplySubmit = async (e) => {
        e.preventDefault();
        if (!replyText.trim()) return;
        setIsSubmitting(true);
        try {
            const res = await fetch(`http://localhost:8080/api/posts/${postId}/comments`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ content: replyText, parent_id: comment.id }),
            });

            const data = await res.json();
            if (res.ok) {
                onAddReply(comment.id, data.comment);
                setReplyText('');
                setShowReply(false);
            }
        } catch (err) {
            console.error('Failed to post reply', err);
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <li className="border p-2 rounded overflow-hidden">
            <div className="flex items-start gap-2">
                <div className="flex-1">
                    <div className="flex items-start justify-between">
                        <p className="text-sm font-semibold">{name}</p>
                        <p className="text-xs text-gray-500">{formatTimeAgo(createdAt)}</p>
                    </div>
                    <p className="text-sm break-words">{comment.content}</p>
                    <button
                        className="text-xs text-blue-600 mt-1"
                        onClick={() => setShowReply((p) => !p)}
                    >
                        {showReply ? 'Cancel' : 'Reply'}
                    </button>

                    {replyCount > 0 && (
                        <button
                            className="text-xs text-gray-600 mt-1 ml-3"
                            onClick={() => setExpanded((p) => !p)}
                        >
                            {expanded ? 'Hide replies' : `View replies (${replyCount})`}
                        </button>
                    )}

                    {showReply && (
                        <form onSubmit={handleReplySubmit} className="mt-2 flex flex-col gap-2">
                            <textarea
                                className="w-full border rounded p-2 text-sm"
                                rows={2}
                                placeholder="Write a reply..."
                                value={replyText}
                                onChange={(e) => setReplyText(e.target.value)}
                            ></textarea>
                            <div>
                                <button
                                    type="submit"
                                    disabled={isSubmitting}
                                    className="bg-blue-600 hover:bg-blue-700 text-white text-xs px-2 py-1 rounded"
                                >
                                    {isSubmitting ? 'Posting...' : 'Post Reply'}
                                </button>
                            </div>
                        </form>
                    )}

                    {expanded && Array.isArray(comment.replies) && comment.replies.length > 0 && (
                        <ul className="mt-2 space-y-2 border-l pl-3">
                            {comment.replies.map((child) => (
                                <CommentItem
                                    key={child.id}
                                    comment={child}
                                    postId={postId}
                                    onAddReply={onAddReply}
                                    level={level + 1}
                                />)
                            )}
                        </ul>
                    )}
                </div>
            </div>
        </li>
    );
}
