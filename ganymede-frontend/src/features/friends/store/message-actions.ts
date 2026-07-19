import { ChatMessage, useMessagesStore } from "./message-store";

export const messageActions = {
  addMessage(friendID: string, message: ChatMessage) {
    useMessagesStore.setState((state) => ({
      messages: {
        ...state.messages,
        [friendID]: [...(state.messages[friendID] ?? []), message],
      },
    }));
  },

  setMessages(friendID: string, messages: ChatMessage[]) {
    useMessagesStore.setState((state) => ({
      messages: {
        ...state.messages,
        [friendID]: messages,
      },
    }));
  },

  clear() {
    useMessagesStore.setState({
      messages: {},
    });
  },
};
