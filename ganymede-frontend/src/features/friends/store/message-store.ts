import { create } from "zustand";

export interface ChatMessage {
  id: string;
  conversationID: string;
  senderID: string;
  recipientID: string;

  content: string;

  createdAt: string;

  pending?: boolean;
}

interface MessagesState {
  messages: Record<string, ChatMessage[]>;

  setMessages: (
    conversationID: string,
    messages: ChatMessage[],
  ) => void;

  addMessage: (
    conversationID: string,
    message: ChatMessage,
  ) => void;

  clear(): void;
}

export const useMessagesStore =
  create<MessagesState>((set) => ({
    messages: {},

    setMessages(conversationID, messages) {
      set((state) => ({
        messages: {
          ...state.messages,
          [conversationID]: messages,
        },
      }));
    },

    addMessage(conversationID, message) {
      set((state) => ({
        messages: {
          ...state.messages,
          [conversationID]: [
            ...(state.messages[conversationID] ?? []),
            message,
          ],
        },
      }));
    },

    clear() {
      set({
        messages: {},
      });
    },
  }));
