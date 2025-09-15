export const formattedMessageDate = (timestamp = new Date()) => {
    const date = new Date(timestamp);

    return new Intl.DateTimeFormat('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
        hour12: true
    }).format(date);
}
