import { http } from "@/lib/http";
import type { ChatMessage } from "../store/types";

interface MessageDTO {
  id: string;
  conversation_id: string;
  sender_id: string;
  recipient_id?: string;
  content: string;
  created_at: string;
}

export async function getMessages(
  conversationID: string,
): Promise<ChatMessage[]> {
  const { data } = await http.get<MessageDTO[]>(
    `/chat/conversations/${conversationID}/messages`,
  );

  return data.map((message) => ({
    id: message.id,
    conversationID: message.conversation_id,
    senderID: message.sender_id,
    recipientID: message.recipient_id,
    content: message.content,
    createdAt: message.created_at,
  }));
}
