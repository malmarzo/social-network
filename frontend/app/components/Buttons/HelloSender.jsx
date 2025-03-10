import React from "react";
import { useWebSocket } from "@/context/Websocket";

const HelloSender = () => {
    const { sendMessage } = useWebSocket();

  const sendHello = () => {
    const msg = {};
    msg.type = "hello";
    msg.content = "Hello World!";
    sendMessage(msg);
  };

  return (
    <div>
      <button
        className="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200"
        onClick={sendHello}
      >
        Send Hello
      </button>
    </div>
  );
};

export default HelloSender;
