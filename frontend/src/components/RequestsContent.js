// app/components/RequestsContent.js
// This component needs to be a Client Component because it contains interactive elements (buttons).
// It will be rendered within app/friends/page.js which is already a client component.

export default function RequestsContent() {
  // Dummy data for requests. In a real application, you'd fetch this from an API.
  const requests = [
    { id: 1, name: "Alice Smith", username: "alices", avatar: "https://randomuser.me/api/portraits/women/1.jpg" },
    { id: 2, name: "Bob Johnson", username: "bobj", avatar: "https://randomuser.me/api/portraits/men/2.jpg" },
    { id: 3, name: "Charlie Brown", username: "charlieb", avatar: "https://randomuser.me/api/portraits/men/3.jpg" },
    { id: 4, name: "Diana Prince", username: "dianap", avatar: "https://randomuser.me/api/portraits/women/4.jpg" },
  ];

  // Handler functions for the buttons
  const handleAccept = (requestId, requesterName) => {
    alert(`Accepted request from ${requesterName} (ID: ${requestId})!`);
    // In a real app, send API request and update state
  };

  const handleDecline = (requestId, requesterName) => {
    alert(`Declined request from ${requesterName} (ID: ${requestId}).`);
    // In a real app, send API request and update state
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-semibold mb-4 text-gray-800">Pending Friend Requests</h2>

      {requests.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {requests.map((request) => (
            <div key={request.id} className="bg-gray-50 rounded-lg shadow-sm p-5 flex flex-col items-center text-center">
              {/* Avatar Section */}
              <div className="w-24 h-24 mb-4 rounded-full overflow-hidden border-2 border-blue-500 flex items-center justify-center bg-gray-200">
                <img
                  src={request.avatar}
                  alt={`${request.name}'s avatar`}
                  className="w-full h-full object-cover"
                />
              </div>

              {/* User Info */}
              <h3 className="text-xl font-semibold text-gray-900">{request.name}</h3>
              <p className="text-gray-500 mb-4">@{request.username}</p>

              {/* Action Buttons */}
              <div className="flex space-x-4">
                <button
                  onClick={() => handleAccept(request.id, request.name)}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition duration-200"
                >
                  Accept
                </button>
                <button
                  onClick={() => handleDecline(request.id, request.name)}
                  className="px-4 py-2 bg-gray-300 text-gray-800 rounded-md hover:bg-gray-400 transition duration-200"
                >
                  Decline
                </button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <p className="text-center text-gray-600 text-lg mt-10">
          You have no pending friend requests at the moment.
        </p>
      )}
    </div>
  );
}