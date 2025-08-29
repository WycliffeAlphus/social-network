let sock = null

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
        console.log("WebSocket connected.");
    }

    sock.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        console.log(msg)
    };

    sock.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
}

export function getWebSocket() {
    return sock;
}
