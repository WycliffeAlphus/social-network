"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { CalendarIcon, ClockIcon, UserGroupIcon, PlusIcon, CheckIcon, XMarkIcon } from '@heroicons/react/24/outline';

export default function GroupEventsPage() {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showCreateEvent, setShowCreateEvent] = useState(false);
  const [responding, setResponding] = useState({});
  const params = useParams();
  const router = useRouter();
  const { id: groupId } = params;

  useEffect(() => {
    if (groupId) {
      fetchEvents();
    }
  }, [groupId]);

  const fetchEvents = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/groups/${groupId}/events`, {
        credentials: 'include',
      });

      if (response.status === 401) {
        router.push('/login');
        return;
      }

      if (!response.ok) {
        throw new Error('Failed to fetch events');
      }

      const responseData = await response.json();
      setEvents(responseData.data || []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const handleRespond = async (eventId, response) => {
    try {
      setResponding(prev => ({ ...prev, [eventId]: true }));

      const res = await fetch(`http://localhost:8080/api/events/${eventId}/respond`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ response }),
      });

      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.message || 'Failed to respond to event');
      }

      // Refresh events after responding
      await fetchEvents();
    } catch (error) {
      console.error('Error responding to event:', error);
      alert(`Error: ${error.message}`);
    } finally {
      setResponding(prev => ({ ...prev, [eventId]: false }));
    }
  };

  const formatEventTime = (timeString) => {
    const date = new Date(timeString);
    return date.toLocaleString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (loading) {
    return (
      <div className="flex min-h-screen justify-center items-center">
        <p className="text-xl">Loading events...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-screen justify-center items-center">
        <p className="text-xl">Error: {error}</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Group Events</h1>
        <div className="flex gap-3">
          <Link href={`/groups/${groupId}`}>
            <button className="border border-black hover:bg-black hover:text-white font-bold py-2 px-4 transition duration-300 ease-in-out">
              Back to Group
            </button>
          </Link>
          <button
            onClick={() => setShowCreateEvent(true)}
            className="bg-black hover:bg-gray-800 text-white font-bold py-2 px-4 transition duration-300 ease-in-out flex items-center gap-2"
          >
            <PlusIcon className="h-5 w-5" />
            Create Event
          </button>
        </div>
      </div>

      {events.length === 0 ? (
        <div className="text-center py-12">
          <CalendarIcon className="h-16 w-16 mx-auto text-gray-400 mb-4" />
          <p className="text-gray-600 mb-4">No events yet. Create the first event for this group!</p>
        </div>
      ) : (
        <div className="space-y-6">
          {events.map((event) => (
            <div key={event.id} className="border border-gray-400 p-6 hover:shadow-lg transition-shadow duration-300">
              <div className="flex justify-between items-start mb-4">
                <div className="flex-1">
                  <h2 className="text-2xl font-semibold mb-2">{event.title}</h2>
                  <p className="text-gray-600 mb-3">{event.description}</p>

                  <div className="flex items-center gap-2 text-gray-700 mb-4">
                    <ClockIcon className="h-5 w-5" />
                    <span className="font-medium">{formatEventTime(event.event_time)}</span>
                  </div>

                  <div className="flex items-center gap-6 text-sm">
                    <div className="flex items-center gap-2">
                      <CheckIcon className="h-5 w-5 text-gray-600" />
                      <span className="font-medium">{event.going_count} Going</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <XMarkIcon className="h-5 w-5 text-gray-600" />
                      <span className="font-medium">{event.not_going_count} Not Going</span>
                    </div>
                  </div>
                </div>
              </div>

              <div className="border-t border-gray-300 pt-4 mt-4">
                <p className="text-sm font-medium mb-3">Your Response:</p>
                <div className="flex gap-3">
                  <button
                    onClick={() => handleRespond(event.id, 'going')}
                    disabled={responding[event.id]}
                    className={`flex items-center gap-2 px-4 py-2 font-medium transition-all duration-200 ${
                      event.user_response === 'going'
                        ? 'bg-black text-white'
                        : responding[event.id]
                        ? 'border border-gray-400 text-gray-400 cursor-not-allowed'
                        : 'border border-black hover:bg-black hover:text-white'
                    }`}
                  >
                    <CheckIcon className="h-5 w-5" />
                    Going
                  </button>
                  <button
                    onClick={() => handleRespond(event.id, 'not_going')}
                    disabled={responding[event.id]}
                    className={`flex items-center gap-2 px-4 py-2 font-medium transition-all duration-200 ${
                      event.user_response === 'not_going'
                        ? 'bg-black text-white'
                        : responding[event.id]
                        ? 'border border-gray-400 text-gray-400 cursor-not-allowed'
                        : 'border border-black hover:bg-black hover:text-white'
                    }`}
                  >
                    <XMarkIcon className="h-5 w-5" />
                    Not Going
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {showCreateEvent && (
        <CreateEventModal
          groupId={groupId}
          onClose={() => setShowCreateEvent(false)}
          onEventCreated={() => {
            setShowCreateEvent(false);
            fetchEvents();
          }}
        />
      )}
    </div>
  );
}

function CreateEventModal({ groupId, onClose, onEventCreated }) {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    event_time: '',
  });
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (!formData.title || !formData.event_time) {
      setError('Title and event time are required');
      return;
    }

    try {
      setCreating(true);

      // Convert to ISO 8601 format
      const eventTime = new Date(formData.event_time).toISOString();

      const response = await fetch(`http://localhost:8080/api/groups/${groupId}/events`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: formData.title,
          description: formData.description,
          event_time: eventTime,
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to create event');
      }

      onEventCreated();
    } catch (error) {
      console.error('Error creating event:', error);
      setError(error.message);
    } finally {
      setCreating(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-white shadow-xl p-6 max-w-md w-full border border-gray-400" onClick={(e) => e.stopPropagation()}>
        <div className="flex justify-between items-center mb-4 border-b border-gray-400 pb-4">
          <h2 className="text-2xl font-bold">Create Event</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-black text-2xl font-bold"
          >
            &times;
          </button>
        </div>

        {error && (
          <div className="mb-4 p-3 border border-black bg-gray-100">
            <p className="text-sm">{error}</p>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">
              Event Title *
            </label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              className="w-full p-2 border border-gray-400 focus:border-black focus:outline-none"
              placeholder="Enter event title"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">
              Description
            </label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="w-full p-2 border border-gray-400 focus:border-black focus:outline-none"
              placeholder="Enter event description"
              rows="3"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">
              Event Date & Time *
            </label>
            <input
              type="datetime-local"
              value={formData.event_time}
              onChange={(e) => setFormData({ ...formData, event_time: e.target.value })}
              className="w-full p-2 border border-gray-400 focus:border-black focus:outline-none"
              required
            />
          </div>

          <div className="flex gap-3 pt-4">
            <button
              type="submit"
              disabled={creating}
              className={`flex-1 py-2 font-medium transition-all duration-200 ${
                creating
                  ? 'border border-gray-400 text-gray-400 cursor-not-allowed'
                  : 'bg-black hover:bg-gray-800 text-white'
              }`}
            >
              {creating ? 'Creating...' : 'Create Event'}
            </button>
            <button
              type="button"
              onClick={onClose}
              className="flex-1 border border-black hover:bg-black hover:text-white py-2 font-medium transition-all duration-200"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
