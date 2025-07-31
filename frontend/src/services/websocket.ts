import { useEffect, useRef, useState } from 'react';

const useWebSocket = (url) => {
    const [socket, setSocket] = useState(null);
    const [messages, setMessages] = useState([]);
    const [error, setError] = useState(null);
    const isMounted = useRef(true);

    useEffect(() => {
        const ws = new WebSocket(url);

        ws.onopen = () => {
            console.log('WebSocket connection established');
        };

        ws.onmessage = (event) => {
            if (isMounted.current) {
                setMessages((prevMessages) => [...prevMessages, event.data]);
            }
        };

        ws.onerror = (event) => {
            setError(event);
        };

        ws.onclose = () => {
            console.log('WebSocket connection closed');
        };

        setSocket(ws);

        return () => {
            isMounted.current = false;
            ws.close();
        };
    }, [url]);

    const sendMessage = (message) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(message);
        }
    };

    return { messages, error, sendMessage };
};

export default useWebSocket;