"use client";

import { Button } from "@/components/ui/button";

import { useAcceptFriendRequest } from "../hooks/use-accept-friend-request";
import { useRejectFriendRequest } from "../hooks/use-reject-friend-request";

import { useIncomingRequests } from "../store/selectors";

import { FriendCard } from "./friend.card";

export function FriendRequests() {
  const requests = useIncomingRequests();

  const acceptMutation = useAcceptFriendRequest();
  const rejectMutation = useRejectFriendRequest();

  if (requests.length === 0) {
    return null;
  }

  return (
    <section className="space-y-3">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold">Friend Requests</h3>

        <span className="text-xs text-muted-foreground">{requests.length}</span>
      </div>

      <div className="space-y-2">
        {requests.map((request) => {
          const friend = {
            id: request.id,
            username: request.username,
            online: false,
            inGame: false,
          };

          return (
            <FriendCard
              key={request.id}
              friend={friend}
              rightSlot={
                <div className="flex gap-2">
                  <Button
                    size="sm"
                    disabled={
                      acceptMutation.isPending || rejectMutation.isPending
                    }
                    onClick={(e) => {
                      e.stopPropagation();

                      acceptMutation.mutate({
                        username: request.username,
                      });
                    }}
                  >
                    Accept
                  </Button>

                  <Button
                    size="sm"
                    variant="outline"
                    disabled={
                      acceptMutation.isPending || rejectMutation.isPending
                    }
                    onClick={(e) => {
                      e.stopPropagation();

                      rejectMutation.mutate({
                        username: request.username,
                      });
                    }}
                  >
                    Reject
                  </Button>
                </div>
              }
            />
          );
        })}
      </div>
    </section>
  );
}
