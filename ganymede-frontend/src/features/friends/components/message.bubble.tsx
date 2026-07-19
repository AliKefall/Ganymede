import { useAuthStore } from "@/features/auth/auth-store";
import { ChatMessage } from "../store/message-store";

interface Props {
  message: ChatMessage;
}

export function MessageBubble({ message }: Props) {
  const me = useAuthStore((s) => s.user?.user_id);

  const mine = message.senderID === me;

  return (
    <div className={`flex ${mine ? "justify-end" : "justify-start"}`}>
      <div
        className={`
                    max-w-[70%]
                    rounded-xl
                    px-4
                    py-2
                    text-sm
                    ${mine ? "bg-primary text-primary-foreground" : "bg-muted"}
                `}
      >
        {message.content}
      </div>
    </div>
  );
}
