export interface ChatMessage {
  id: string;
  conversationID: string;
  senderID: string;
  recipientID?: string;
  content: string;
  createdAt: string;
  pending?: boolean;
}
