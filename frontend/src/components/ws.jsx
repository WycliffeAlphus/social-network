let sock = null
let messageQueue = []

export default function connectWebsocket(id) {
    if (sock && (sock.readyState === WebSocket.OPEN || sock.readyState === WebSocket.CONNECTING)) {
        return // already connected or connecting
    }

    sock = new WebSocket(`ws://localhost:8080/ws`);

    sock.onopen = () => {
        sock.send(JSON.stringify({
            from: id,
            to: "",
            content: ""
        }));

        // send any queued messages
        while (messageQueue.length > 0) {
            const msg = messageQueue.shift()
            sock.send(JSON.stringify(msg))
        }
    }

    sock.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        console.log(msg)
    };

    sock.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    sock.onclose = () => {
        console.log("WebSocket disconnected.");
    };
}

export function getWebSocket() {
    return sock;
}

export function sendMessage(message) {
    const sock = getWebSocket()
    if (sock && sock.readyState === WebSocket.OPEN) {
        sock.send(JSON.stringify(message))
    } else {
        // queue message for when connection is established
        messageQueue.push(message)
        console.log("Message queued, waiting for connection")
    }
}
