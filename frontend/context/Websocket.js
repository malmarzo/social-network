import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useAuth } from "./AuthContext";
import React from "react";

const WebSocketContext = createContext(null);

const WebSocketProvider = ({ children }) => {
  const { isLoggedIn, loading } = useAuth();
  const wsRef = useRef(null); // Store WebSocket instance
  const reconnectRef = useRef(null); // Store reconnection attempt
  const handlersRef = useRef({}); // Store handlers dynamically
  const [isConnected, setIsConnected] = useState(false); // Track WebSocket connection state
  const isLoggedInRef = useRef(isLoggedIn); // Ref to store latest isLoggedIn state

  useEffect(() => {
    // Update ref with the latest value of isLoggedIn
    isLoggedInRef.current = isLoggedIn;
  }, [isLoggedIn]); // if this value changes it will update it in the ref

  const connectWebSocket = () => {
    //if theere is a connection then return
    if (wsRef.current) return;

    //connect to the websocket and update the ref
    const ws = new WebSocket("ws://localhost:8080/ws");
    wsRef.current = ws;

    //When the connection opens
    ws.onopen = () => {
      console.log("WebSocket connected");
      setIsConnected(true);

      if (reconnectRef.current) {
        clearTimeout(reconnectRef.current);
        reconnectRef.current = null;
      }
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        // console.log("WebSocket message received:", msg);

        //Handle the msg based on its type
        if (msg.type && handlersRef.current[msg.type]) {
          handlersRef.current[msg.type](msg);
        } else {
          console.warn("Unhandled message type:", msg.type);
        }
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected");
      wsRef.current = null;
      setIsConnected(false);

      if (isLoggedInRef.current) {
        console.log("User is logged in, attempting to reconnect...");
        reconnectRef.current = setTimeout(() => {
          connectWebSocket();
        }, 3000);
      } else {
        console.log("User logged out - WebSocket will not reconnect.");
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      ws.close();
    };
  };

  useEffect(() => {
    console.log("WebSocket Effect Triggered - isLoggedIn:", isLoggedIn);
    if (loading) return;

    if (isLoggedIn && !isConnected) {
      connectWebSocket(); // Only connect if the user is logged in and WebSocket isn't connected
    } else if (!isLoggedIn && wsRef.current) {
      console.log("User logged out - closing WebSocket...");
      wsRef.current.close(); // Close WebSocket on logout
      wsRef.current = null;
      setIsConnected(false); // Update connection state on logout

      if (reconnectRef.current) {
        clearTimeout(reconnectRef.current); // Clear reconnection timeout on logout
        reconnectRef.current = null;
      }
    }
  }, [isLoggedIn, loading, isConnected]);

  // Function to allow components to register message handlers
  const addMessageHandler = (type, handler) => {
    handlersRef.current = {
      ...handlersRef.current,
      [type]: handler,
    };
  };

  // Function to send messages
  const sendMessage = (message) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else {
      console.warn("WebSocket is not connected. Message not sent:", message);
    }
  };

  return (
    <WebSocketContext.Provider value={{ addMessageHandler, sendMessage }}>
      {children}
    </WebSocketContext.Provider>
  );
};

export default WebSocketProvider;

export const useWebSocket = () => useContext(WebSocketContext);