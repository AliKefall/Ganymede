import { useChatStore } from "./chat-store";
import { useMessagesStore } from "./message-store";
import { useFriendsStore } from "./store";

export const useFriendsList = () => useFriendsStore((state) => state.friends);

export const useOnlineFriends = () =>
  useFriendsStore((state) => state.friends.filter((friend) => friend.online));

export const useOfflineFriends = () =>
  useFriendsStore((state) => state.friends.filter((friend) => !friend.online));

export const useIncomingRequests = () =>
  useFriendsStore((state) => state.incomingRequests);

export const useOutgoingRequests = () =>
  useFriendsStore((state) => state.outgoingRequests);

export const useIncomingRequestCount = () =>
  useFriendsStore((state) => state.incomingRequests.length);

export const useOutgoingRequestCount = () =>
  useFriendsStore((state) => state.outgoingRequests.length);

// Chat Messages

export const useIsChatOpen = () =>
    useChatStore((state) => state.isOpen);

export const useSelectedFriend = () =>
    useChatStore((state) => state.selectedFriend);

export const useConversation = (friendID?: string) =>
    useMessagesStore((state) =>
        friendID ? (state.messages[friendID] ?? []) : [],
    );
