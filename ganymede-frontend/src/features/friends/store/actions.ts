import { useFriendsStore } from "./store";
import type { Friend, FriendRequest } from "./types";

export const friendsActions = {
  setFriends(friends: Friend[]) {
    useFriendsStore.setState({ friends });
  },

  setIncomingRequest(requests: FriendRequest[]) {
    useFriendsStore.setState({
      incomingRequests: requests,
    });
  },

  setOutgoingRequests(requests: FriendRequest[]) {
    useFriendsStore.setState({
      outgoingRequests: requests,
    });
  },

  addFriend(friend: Friend) {
    useFriendsStore.setState((state) => ({
      friends: [...state.friends, friend],
    }));
  },

  removeFriend(friendID: string) {
    useFriendsStore.setState((state) => ({
      friends: state.friends.filter((friend) => friend.id !== friendID),
    }));
  },

  setFriendOnline(friendID: string) {
    useFriendsStore.setState((state) => ({
      friends: state.friends.map((friend) =>
        friend.id === friendID ? { ...friend, online: true } : friend,
      ),
    }));
  },

  setFriendOffline(friendID: string) {
    useFriendsStore.setState((state) => ({
      friends: state.friends.map((friend) =>
        friend.id === friendID ? { ...friend, online: false } : friend,
      ),
    }));
  },

  addIncomingRequest(request: FriendRequest) {
    useFriendsStore.setState((state) => ({
      incomingRequests: [request, ...state.incomingRequests],
    }));
  },

  removeIncomingRequest(userID: string) {
    useFriendsStore.setState((state) => ({
      incomingRequests: state.incomingRequests.filter(
        (request) => request.id !== userID,
      ),
    }));
  },

  addOutgoingRequest(request: FriendRequest) {
    useFriendsStore.setState((state) => ({
      outgoingRequests: [request, ...state.outgoingRequests],
    }));
  },

  removeOutgoingRequest(userID: string) {
    useFriendsStore.setState((state) => ({
      outgoingRequests: state.outgoingRequests.filter(
        (request) => request.id !== userID,
      ),
    }));
  },

  clear() {
    useFriendsStore.setState({
      friends: [],
      incomingRequests: [],
      outgoingRequests: [],
    });
  },
};
