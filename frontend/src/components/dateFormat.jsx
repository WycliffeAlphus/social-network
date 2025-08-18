export function formatTimeAgo(dateString) {
  const date = new Date(dateString);
  const now = new Date();
  const diff = Math.floor((now - date) / 1000); // in seconds

  if (diff < 60) return 'just now';
  if (diff < 3600) return `${Math.floor(diff / 60)} minute(s) ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)} hour(s) ago`;
  if (diff < 604800) return `${Math.floor(diff / 86400)} day(s) ago`;

  // Fallback to full date: dd/mm/yyyy
  return date.toLocaleDateString('en-GB');
}
