'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import Loading from './loading';
import { useParams } from 'next/navigation';

export function ContactedChats() {
  const [contacts, setContacts] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const { receiverId } = useParams()

  useEffect(() => {
    const fetchContacts = async () => {
      try {
        const response = await fetch('http://localhost:3000/api/users', {
          credentials: 'include',
        });

        if (!response.ok) {
          throw new Error('Failed to fetch contacted users')
        }

        const data = await response.json();

        setContacts(data);
      } catch (error) {
        setError(error.message)
      } finally {
        setLoading(false)
      }
    }

    fetchContacts();

    const handleMessageSent = () => [
      fetchContacts()
    ]

    window.addEventListener('messageEvent', handleMessageSent)

    return () => {
      window.removeEventListener('messageEvent', handleMessageSent)
    }
  }, []);

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <div className="text-red-500">{error}</div>
  }

  return (
    <div className="flex flex-col gap-3">
      {contacts && contacts.length > 0 ? (
        contacts.map((contact) => (
          <Link
            key={contact.id}
            href={`/messages/${contact.id}`}
            className={`p-3 flex items-center gap-3 transition-colors rounded-xl
            ${receiverId === contact.id
                ? 'bg-gray-200 dark:bg-gray-600'
                : 'hover:bg-gray-600/30'
              }
            `}
          >
            {contact.avatar?.String ? (
              <img
                src={contact.avatar.String}
                alt="User avatar"
                className="w-12 h-12 rounded-full object-cover"
              />
            ) : (
              <div className="w-[clamp(2rem,4vw,3.5rem)] h-[clamp(2rem,4vw,3.5rem)] rounded-full bg-gray-200 overflow-hidden flex items-center justify-center">
                <span className="text-lg text-black font-medium">
                  {contact.firstname?.charAt(0).toUpperCase()}
                  {contact.lastname?.charAt(0).toUpperCase()}
                </span>
              </div>
            )}

            <div>
              <span className="font-medium">
                {contact.firstname} {contact.lastname}
              </span>
            </div>
          </Link>
        ))
      ) : (
        <p>You haven't messaged anyone yet.</p>
      )}
    </div>
  );
} 
