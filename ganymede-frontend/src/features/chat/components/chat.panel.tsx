"use client";

import { SyntheticEvent, useState } from "react";
import { SendHorizonal, X } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { websocketManager } from "@/lib/websocket";

import { useMessages } from "../hooks/use-messages";
import { useChatStore } from "../store/chat-store";
import {
  useSelectedConversation,
  useSelectedConversationID,
} from "../store/selectors";
import { MessageBubble } from "./message.bubble";

export function ChatPanel() {
  const isOpen = useChatStore((state) => state.isOpen);
  const selectedFriend = useChatStore((state) => state.selectedFriend);
  const close = useChatStore((state) => state.close);
  const conversationID = useSelectedConversationID();
  const messages = useSelectedConversation();

  const [message, setMessage] = useState("");
  const [sendError, setSendError] = useState<string | null>(null);

  useMessages(conversationID);

  if (!selectedFriend) return null;

  function handleSubmit(e: SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();

    const content = message.trim();
    if (!content || !selectedFriend) return;

    try {
      websocketManager.send("send_message", {
        recipient_id: selectedFriend.id,
        conversation_id: conversationID,
        content,
      });
      setMessage("");
      setSendError(null);
    } catch (error) {
      setSendError(
        error instanceof Error ? error.message : "Message could not be sent",
      );
    }
  }

  return (
    <aside
      className={`fixed right-80 top-0 z-40 flex h-screen w-[420px] flex-col border-l bg-background shadow-xl transition-transform duration-300 ${
        isOpen ? "translate-x-0" : "translate-x-full"
      }`}
    >
      <header className="flex items-center justify-between border-b px-5 py-4">
        <div className="flex items-center gap-3">
          <div
            className={`h-3 w-3 rounded-full ${selectedFriend.online ? "bg-green-500" : "bg-zinc-500"}`}
          />
          <div>
            <h2 className="font-semibold">{selectedFriend.username}</h2>
            <p className="text-xs text-muted-foreground">
              {selectedFriend.online ? "Online" : "Offline"}
            </p>
          </div>
        </div>
        <Button variant="ghost" size="icon" onClick={close}>
          <X className="h-5 w-5" />
        </Button>
      </header>

      <div className="flex-1 overflow-y-auto px-4 py-6">
        <div className="flex flex-col gap-3">
          {messages.length === 0 ? (
            <p className="text-center text-sm text-muted-foreground">
              No messages yet. Say hello!
            </p>
          ) : (
            messages.map((msg) => <MessageBubble key={msg.id} message={msg} />)
          )}
        </div>
      </div>

      <form onSubmit={handleSubmit} className="border-t p-4">
        {sendError && (
          <p className="mb-2 text-xs text-destructive">{sendError}</p>
        )}
        <div className="flex gap-2">
          <Input
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Type a message..."
            maxLength={500}
          />
          <Button type="submit" size="icon" disabled={!message.trim()}>
            <SendHorizonal className="h-4 w-4" />
          </Button>
        </div>
      </form>
    </aside>
  );
}
