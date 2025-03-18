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
    // If there is a connection then return
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) return;

    // Close any existing connection that might be in a bad state
    if (wsRef.current) {
      try {
        wsRef.current.close();
      } catch (err) {
        console.error("Error closing existing WebSocket:", err);
      }
      wsRef.current = null;
    }

    // Connect to the websocket and update the ref
    console.log("Attempting to connect to WebSocket...");
    const ws = new WebSocket("ws://localhost:8080/ws");
    wsRef.current = ws;

    // When the connection opens
    ws.onopen = () => {
      console.log("WebSocket connected successfully");
      setIsConnected(true); 

      if (reconnectRef.current) {
        clearTimeout(reconnectRef.current);
        reconnectRef.current = null;
      }
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        console.log("WebSocket message received:", msg);

        // Handle the msg based on its type
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
    console.log("WebSocket Effect Triggered - isLoggedIn:", isLoggedIn, "loading:", loading, "isConnected:", isConnected);
    if (loading) return;

    if (isLoggedIn) {
      // Check if we need to connect or reconnect
      if (!wsRef.current || wsRef.current.readyState === WebSocket.CLOSED || wsRef.current.readyState === WebSocket.CLOSING) {
        console.log("User is logged in and WebSocket is not connected or in a bad state. Connecting...");
        connectWebSocket();
      }
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
    
    // Set up a periodic connection check for logged-in users
    let connectionCheckInterval;
    if (isLoggedIn) {
      connectionCheckInterval = setInterval(() => {
        if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
          console.log("Periodic check: WebSocket not connected. Attempting to reconnect...");
          connectWebSocket();
        }
      }, 30000); // Check every 30 seconds
    }
    
    return () => {
      if (connectionCheckInterval) {
        clearInterval(connectionCheckInterval);
      }
    };
  }, [isLoggedIn, loading]);

  // Function to allow components to register message handlers
  const addMessageHandler = (type, handler) => {
    handlersRef.current = {
      ...handlersRef.current,
      [type]: handler,
    };
  };

  // Function to send messages with message ID for tracking
  const sendMessage = (message) => {
    // Add a unique message ID if not present
    const messageWithId = {
      ...message,
      messageId: message.messageId || `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    };
    
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      try {
        wsRef.current.send(JSON.stringify(messageWithId));
        console.log("Message sent successfully:", messageWithId.type);
        return true; // Return success status
      } catch (error) {
        console.error("Error sending WebSocket message:", error);
        // Try to reconnect on send error
        if (isLoggedInRef.current) {
          console.log("Attempting to reconnect after send error...");
          setTimeout(connectWebSocket, 1000);
        }
        return false;
      }
    } else {
      console.warn("WebSocket is not connected (state: " + 
                  (wsRef.current ? wsRef.current.readyState : "null") + 
                  "). Message not sent:", messageWithId.type);
      
      // Try to reconnect if we're supposed to be connected
      if (isLoggedInRef.current && (!wsRef.current || wsRef.current.readyState !== WebSocket.CONNECTING)) {
        console.log("Attempting to reconnect...");
        setTimeout(connectWebSocket, 1000);
      }
      
      return false;
    }
  };
  
  // Function to check connection status
  const checkConnection = () => {
    return {
      isConnected: isConnected,
      readyState: wsRef.current ? wsRef.current.readyState : WebSocket.CLOSED
    };
  };

  return (
    <WebSocketContext.Provider value={{ addMessageHandler, sendMessage, isConnected, checkConnection }}>
      {children}
    </WebSocketContext.Provider>
  );
};

export default WebSocketProvider;

export const useWebSocket = () => useContext(WebSocketContext);
