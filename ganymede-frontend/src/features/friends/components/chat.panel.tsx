"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

import { SendHorizonal, X } from "lucide-react";

import { useState } from "react";

import { useChatStore } from "../store/chat-store";

export function ChatPanel() {
  const isOpen = useChatStore((state) => state.isOpen);
  const selectedFriend = useChatStore((state) => state.selectedFriend);
  const close = useChatStore((state) => state.close);

  const [message, setMessage] = useState("");

  if (!selectedFriend) return null;

  const dummyMessages = [
    {
      id: "1",
      mine: false,
      text: "Hello 👋",
      createdAt: "09:42",
    },
    {
      id: "2",
      mine: true,
      text: "Hi!",
      createdAt: "09:43",
    },
    {
      id: "3",
      mine: false,
      text: "How are you doing?",
      createdAt: "09:44",
    },
  ];

  function handleSubmit(e: React.SyntheticEvent) {
    e.preventDefault();

    if (!message.trim()) return;

    console.log(message);

    setMessage("");
  }

  return (
    <aside
      className={`
        fixed
        right-80
        top-0
        z-40
        flex
        h-screen
        w-[420px]
        flex-col
        border-l
        bg-background
        shadow-xl
        transition-transform
        duration-300
        ${isOpen ? "translate-x-0" : "translate-x-full"}
      `}
    >
      {/* Header */}

      <header className="flex items-center justify-between border-b px-5 py-4">
        <div className="flex items-center gap-3">
          <div
            className={`h-3 w-3 rounded-full ${
              selectedFriend.online ? "bg-green-500" : "bg-zinc-500"
            }`}
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

      {/* Messages */}

      <div className="flex-1 overflow-y-auto px-4 py-6">
        <div className="flex flex-col gap-3">
          {dummyMessages.map((msg) => (
            <div
              key={msg.id}
              className={`flex ${msg.mine ? "justify-end" : "justify-start"}`}
            >
              <div
                className={`
                  max-w-[75%]
                  rounded-2xl
                  px-4
                  py-2
                  ${
                    msg.mine ? "bg-primary text-primary-foreground" : "bg-muted"
                  }
                `}
              >
                <p className="break-words text-sm">{msg.text}</p>

                <p
                  className={`mt-1 text-right text-[10px] ${
                    msg.mine
                      ? "text-primary-foreground/70"
                      : "text-muted-foreground"
                  }`}
                >
                  {msg.createdAt}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Footer */}

      <form onSubmit={handleSubmit} className="border-t p-4">
        <div className="flex gap-2">
          <Input
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Type a message..."
          />

          <Button type="submit" size="icon" disabled={!message.trim()}>
            <SendHorizonal className="h-4 w-4" />
          </Button>
        </div>
      </form>
    </aside>
  );
}
