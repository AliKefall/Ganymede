import { create } from "zustand";
import type { Friend } from "@/features/friends/store/types";
import type { ChatMessage } from "./types";

interface ChatState {
  isOpen: boolean;
  selectedFriend: Friend | null;
  conversationByFriend: Record<string, string>;
  messages: Record<string, ChatMessage[]>;

  open(friend: Friend): void;
  close(): void;
  setConversation(friendID: string, conversationID: string): void;
  setMessages(conversationID: string, messages: ChatMessage[]): void;
  addMessage(conversationID: string, message: ChatMessage): void;
  clear(): void;
}

export const useChatStore = create<ChatState>((set) => ({
  isOpen: false,
  selectedFriend: null,
  conversationByFriend: {},
  messages: {},

  open(friend) {
    set({ isOpen: true, selectedFriend: friend });
  },

  close() {
    set({ isOpen: false, selectedFriend: null });
  },

  setConversation(friendID, conversationID) {
    set((state) => ({
      conversationByFriend: {
        ...state.conversationByFriend,
        [friendID]: conversationID,
      },
    }));
  },

  setMessages(conversationID, messages) {
    set((state) => ({
      messages: { ...state.messages, [conversationID]: messages },
    }));
  },

  addMessage(conversationID, message) {
    set((state) => {
      const existing = state.messages[conversationID] ?? [];
      if (existing.some((item) => item.id === message.id)) return state;
      return {
        messages: {
          ...state.messages,
          [conversationID]: [...existing, message],
        },
      };
    });
  },

  clear() {
    set({ conversationByFriend: {}, messages: {} });
  },
}));
