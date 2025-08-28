let sock

export default function connectWebsocket(id) {
    sock = new WebSocket(`ws://localhost:8080/ws`);

    sock.onopen = () => {
        sock.send(JSON.stringify({
            from: id,
            to: "",
            content: "",
            timestamp: new Date().toISOString()
        }))
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
