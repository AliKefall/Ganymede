import {
    handleFriendOffline,
    handleFriendOnline,
    handleFriendRequestAccepted,
    handleFriendRequestReceived,
    handleFriendRequestRejected,
} from "@/features/friends/events/events";

import type { Friend, FriendRequest } from "@/features/friends/store/types";

export interface WebSocketMessage<T = unknown> {
    id?: string;
    type: string;
    payload: T;
    created_at?: string;
}

interface FriendPresencePayload {
    user_id: string;
    username: string;
    online: boolean;
}

interface FriendRequestPayload {
    user: FriendRequest;
}

interface FriendAcceptedPayload {
    user: Friend;
}

interface FriendRejectedPayload {
    user: {
        id: string;
    };
}

type EventHandler = (payload: unknown) => void;

const handlers: Record<string, EventHandler> = {
    friend_online: (payload) => {
        const p = payload as FriendPresencePayload;
        handleFriendOnline(p.user_id);
    },

    friend_offline: (payload) => {
        const p = payload as FriendPresencePayload;
        handleFriendOffline(p.user_id);
    },

    friend_request_received: (payload) => {
        const p = payload as FriendRequestPayload;
        handleFriendRequestReceived(p.user);
    },

    friend_request_accepted: (payload) => {
        const p = payload as FriendAcceptedPayload;
        handleFriendRequestAccepted(p.user);
    },

    friend_request_rejected: (payload) => {
        const p = payload as FriendRejectedPayload;
        handleFriendRequestRejected(p.user.id);
    },
};

export function dispatchWebSocketEvent(
    message: WebSocketMessage,
) {
    const handler = handlers[message.type];

    if (!handler) {
        console.warn(
            `Unhandled websocket event: ${message.type}`,
        );
        return;
    }

    handler(message.payload);
}
