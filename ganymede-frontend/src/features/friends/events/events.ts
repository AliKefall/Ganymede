import { friendsActions } from "../store/actions";
import { Friend, FriendRequest } from "../store/types";

export function handleFriendOnline(
    friendID: string,
){
    friendsActions.setFriendOnline(friendID);
}

export function handleFriendOffline(
    friendID: string,
){
    friendsActions.setFriendOffline(friendID);
}

export function handleFriendRequestReceived(
    request: FriendRequest,
){
    friendsActions.addIncomingRequest(request);
}

export function handleFriendRequestAccepted(
    friend: Friend,
){
    friendsActions.removeOutgoingRequest(friend.id);

    friendsActions.addFriend(friend);
}

export function handleFriendRequestRejected(
    userID: string,
){
    friendsActions.removeOutgoingRequest(userID)
}
