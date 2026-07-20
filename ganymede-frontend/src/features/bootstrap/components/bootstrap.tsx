"use client";

import { useEffect } from "react";

import { useAuthStore } from "@/features/auth/auth-store";
import { useFriendRequests } from "@/features/friends/hooks/use-friend-requests";
import { useFriends } from "@/features/friends/hooks/use-friends";
import { dispatchWebSocketEvent } from "@/lib/dispatcher";
import { websocketManager } from "@/lib/websocket";

export function Bootstrap() {
  const accessToken = useAuthStore((state) => state.accessToken);

  useFriends();
  useFriendRequests();

  useEffect(() => {
    if (!accessToken) return;

    websocketManager.connect(accessToken);
    return websocketManager.subscribe(dispatchWebSocketEvent);
  }, [accessToken]);

  return null;
}
