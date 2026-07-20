import { chatActions } from "../store/actions";

interface NewMessagePayload {
  id: string;
  conversation_id: string;
  sender_id: string;
  recipient_id: string;
  content: string;
  created_at: string;
}

export function handleNewMessage(payload: NewMessagePayload) {
  chatActions.setConversation(payload.sender_id, payload.conversation_id);
  chatActions.setConversation(payload.recipient_id, payload.conversation_id);
  chatActions.addMessage(payload.conversation_id, {
    id: payload.id,
    conversationID: payload.conversation_id,
    senderID: payload.sender_id,
    recipientID: payload.recipient_id,
    content: payload.content,
    createdAt: payload.created_at,
  });
}
