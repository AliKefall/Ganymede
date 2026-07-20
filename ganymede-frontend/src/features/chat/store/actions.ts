import { useChatStore } from "./chat-store";
import type { ChatMessage } from "./types";

export const chatActions = {
  open: useChatStore.getState().open,
  close: useChatStore.getState().close,

  setConversation(friendID: string, conversationID: string) {
    useChatStore.getState().setConversation(friendID, conversationID);
  },

  setMessages(conversationID: string, messages: ChatMessage[]) {
    useChatStore.getState().setMessages(conversationID, messages);
  },

  addMessage(conversationID: string, message: ChatMessage) {
    useChatStore.getState().addMessage(conversationID, message);
  },

  clear() {
    useChatStore.getState().clear();
  },
};
