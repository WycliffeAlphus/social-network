let sock

export default function connectWebsocket(id) {
    sock = new WebSocket(`ws://localhost:8080/ws`);

    sock.onopen = () => {
        if (sock.readyState === WebSocket.OPEN) { // <-- this explicity checks the socket's on-open status to prevent race condition when useffect is called multiple times
            sock.send(JSON.stringify({
                from: id,
                to: "",
                content: ""
            }));
        }
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
