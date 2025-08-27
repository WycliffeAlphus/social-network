
export default function connectWebsocket(id) {
    let sock
    sock = new WebSocket(`ws://localhost:8080/ws`);
    
    sock.onopen = () => {
        console.log("hit")
        sock.send(JSON.stringify({
            from: id
        }))

        sock.send(JSON.stringify({
            from: id,
            content: "Test message"
        }));
    }

    sock.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        console.log(msg)
    };

    sock.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
}
