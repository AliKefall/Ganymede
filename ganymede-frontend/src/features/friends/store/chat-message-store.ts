import { create } from "zustand";

export interface ChatMessage {
  id: string;
  senderID: string;
  recipientID: string;
  content: string;
  createdAt: string;
}

interface ChatMessagesState {
  messages: Record<string, ChatMessage[]>;

  addMessage: (friendID: string, message: ChatMessage) => void;

  setMessages: (friendID: string, messages: ChatMessage[]) => void;

  clearConversation: (friendID: string) => void;

  clearAll: () => void;
}

export const useChatMessagesStore =
    create<ChatMessagesState>((set) => ({
    messages: {},

    addMessage: (friendID, message) =>
    set((state) => ({
        messages: {
            ...state.messages,
            [friendID]: [
                ...(state.messages[friendID] ?? []),
                message,
            ]
        }
    })),

    setMessages: (friendID, messages) =>
    set((state) => ({
        messages: {
            ...state.messages,
            [friendID]: messages,
        }
    })),

    clearConversation: (friendID) =>
    set((state) => ({
        messages: {
            ...state.messages,
            [friendID]: [],
        },
    })),


    clearAll: () =>
    set({
        messages: {},
    })
}))
