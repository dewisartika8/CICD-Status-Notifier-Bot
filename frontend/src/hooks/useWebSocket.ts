import { useEffect, useRef, useState } from 'react';

const useWebSocket = (url) => {
    const [socket, setSocket] = useState(null);
    const [messages, setMessages] = useState([]);
    const [error, setError] = useState(null);
    const socketRef = useRef(null);

    useEffect(() => {
        const connect = () => {
            const ws = new WebSocket(url);
            socketRef.current = ws;

            ws.onopen = () => {
                console.log('WebSocket connection established');
            };

            ws.onmessage = (event) => {
                setMessages((prevMessages) => [...prevMessages, event.data]);
            };

            ws.onerror = (event) => {
                setError(event);
            };

            ws.onclose = () => {
                console.log('WebSocket connection closed, reconnecting...');
                setTimeout(connect, 1000); // Reconnect after 1 second
            };
        };

        connect();

        return () => {
            if (socketRef.current) {
                socketRef.current.close();
            }
        };
    }, [url]);

    const sendMessage = (message) => {
        if (socketRef.current) {
            socketRef.current.send(message);
        }
    };

    return { messages, error, sendMessage };
};

export default useWebSocket;