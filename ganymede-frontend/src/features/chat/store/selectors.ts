import { useChatStore } from "./chat-store";
import { ChatMessage } from "./types";

const EMPTY_MESSAGES: ChatMessage[] = [];

export const useIsChatOpen = () => useChatStore((state) => state.isOpen);

export const useSelectedFriend = () =>
  useChatStore((state) => state.selectedFriend);

export const useSelectedConversationID = () =>
  useChatStore((state) => {
    if (!state.selectedFriend) return undefined;

    return state.conversationByFriend[state.selectedFriend.id];
  });

export const useSelectedConversation = () =>
  useChatStore((state) => {
    if (!state.selectedFriend) {
      return EMPTY_MESSAGES;
    }

    const conversationID = state.conversationByFriend[state.selectedFriend.id];

    if (!conversationID) {
      return EMPTY_MESSAGES;
    }

    return state.messages[conversationID] ?? EMPTY_MESSAGES;
  });
