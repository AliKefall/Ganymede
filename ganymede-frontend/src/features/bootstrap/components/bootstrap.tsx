"use client";

import { useFriendRequests } from "@/features/friends/hooks/use-friend-requests";
import { useFriends } from "@/features/friends/hooks/use-friends";

export function Bootstrap(){
    useFriends();
    useFriendRequests();

    return null;
}
