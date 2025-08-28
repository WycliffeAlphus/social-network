'use client';

import { useParams, useRouter } from 'next/navigation';
import { useState } from 'react';

export default function GroupEventsCreatePage() {
  const params = useParams();
  const router = useRouter();
  const groupId = params?.id;

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [time, setTime] = useState('');
  const [location, setLocation] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    setSuccess(null);

    try {
      // Convert local input to RFC3339
      const when = new Date(time).toISOString();

      const res = await fetch(`http://localhost:8080/api/groups/${groupId}/events`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description, time: when, location }),
      });
      const data = await res.json();
      if (!res.ok) {
        throw new Error(data?.error || 'Failed to create event');
      }
      setSuccess('Event created');
      setTitle('');
      setDescription('');
      setTime('');
      setLocation('');
      // Optionally navigate back to group page
      // router.push(`/groups/${groupId}`)
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="p-4 max-w-xl mx-auto">
      <h1 className="text-xl font-semibold mb-4">Create Group Event</h1>
      {error && <p className="text-sm text-red-600 mb-2">{error}</p>}
      {success && <p className="text-sm text-green-600 mb-2">{success}</p>}
      <form onSubmit={handleSubmit} className="space-y-3">
        <div>
          <label className="block text-sm mb-1">Title</label>
          <input className="w-full border rounded p-2 text-sm" value={title} onChange={(e) => setTitle(e.target.value)} required />
        </div>
        <div>
          <label className="block text-sm mb-1">Description</label>
          <textarea className="w-full border rounded p-2 text-sm" rows={3} value={description} onChange={(e) => setDescription(e.target.value)} />
        </div>
        <div>
          <label className="block text-sm mb-1">Date and Time</label>
          <input type="datetime-local" className="w-full border rounded p-2 text-sm" value={time} onChange={(e) => setTime(e.target.value)} required />
        </div>
        <div>
          <label className="block text-sm mb-1">Location</label>
          <input className="w-full border rounded p-2 text-sm" value={location} onChange={(e) => setLocation(e.target.value)} placeholder={'HQ conference room'} />
        </div>
        <button type="submit" disabled={submitting} className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">
          {submitting ? 'Creating...' : 'Create Event'}
        </button>
      </form>
    </div>
  );
}


