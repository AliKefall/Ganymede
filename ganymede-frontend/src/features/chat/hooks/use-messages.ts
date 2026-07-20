"use client";

import { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { getMessages } from "../api/get-messages";
import { chatActions } from "../store/actions";

export function useMessages(conversationID?: string) {
  const query = useQuery({
    queryKey: ["chat-messages", conversationID],
    queryFn: () => getMessages(conversationID!),
    enabled: Boolean(conversationID),
  });

  useEffect(() => {
    if (conversationID && query.data) {
      chatActions.setMessages(conversationID, query.data);
    }
  }, [conversationID, query.data]);

  return query;
}
