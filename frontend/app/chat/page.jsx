import { Suspense } from "react";
import ChatPage from "./ChatPage";

export default function Chat() {
  return (
    <Suspense fallback={<div>Loading chat...</div>}>
      <ChatPage />
    </Suspense>
  );
}
