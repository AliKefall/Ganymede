"use client";

import { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";

import { getFriendRequests } from "../api/get-friend-requests";
import { friendsActions } from "../store/actions";

export function useFriendRequests() {
  const query = useQuery({
    queryKey: ["friend-requests"],
    queryFn: getFriendRequests,
  });

  useEffect(() => {
    if (!query.data) return;

    friendsActions.setIncomingRequest(query.data.incoming);

    friendsActions.setOutgoingRequests(query.data.outgoing);
  }, [query.data]);

  return query;
}
