"use client";

import { useAuthStore } from "@/features/auth/auth-store";
import type { ChatMessage } from "../store/types";

interface Props {
  message: ChatMessage;
}

export function MessageBubble({ message }: Props) {
  const me = useAuthStore((state) => state.user?.user_id);
  const mine = message.senderID === me;
  const sentAt = new Date(message.createdAt).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });

  return (
    <div className={`flex ${mine ? "justify-end" : "justify-start"}`}>
      <div
        className={`max-w-[75%] rounded-2xl px-5 py-2 ${
          mine ? "bg-primary text-primary-foreground" : "bg-muted"
        }`}
      >
        <p className="wrap-break-word text-shadow-md">{message.content}</p>
        <p
          className={`mt-1 text-right text-[10px] ${
            mine ? "text-primary-foreground/70" : "text-muted-foreground"
          }`}
        >
          {sentAt}
        </p>
      </div>
    </div>
  );
}
