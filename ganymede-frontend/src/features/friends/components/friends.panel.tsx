"use client";

import { X } from "lucide-react";

import { Separator } from "@/components/ui/separator";

import { FriendCard } from "./friend.card";
import { FriendRequests } from "./friend.requests";

import { useFriendsPanelStore } from "../store/panel-store";
import { useFriendsList, useIncomingRequests } from "../store/selectors";
import { AddFriendDialog } from "./add.friend.dialog";
import { chatActions } from "@/features/chat/store/actions";

export function FriendsPanel() {
  const isOpen = useFriendsPanelStore((state) => state.isOpen);

  const close = useFriendsPanelStore((state) => state.close);

  const friends = useFriendsList();
  const requests = useIncomingRequests();

  const onlineFriends = friends.filter((friend) => friend.online);

  const offlineFriends = friends.filter((friend) => !friend.online);

  return (
    <aside
      className={`
        fixed
        right-0
        top-0
        z-50
        flex
        h-screen
        w-80
        flex-col
        border-l
        bg-background
        shadow-xl
        transition-transform
        duration-300
        ${isOpen ? "translate-x-0" : "translate-x-full"}
      `}
    >
      <header className="border-b p-4">
        <div className="flex items-start justify-between">
          <div>
            <h2 className="text-xl font-semibold">Friends</h2>

            <p className="text-sm text-muted-foreground">
              {friends.length} {friends.length === 1 ? "Friend" : "Friends"}
            </p>
          </div>

          <button
            onClick={close}
            className="rounded-md p-2 transition-colors hover:bg-muted"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="mt-4">
          <AddFriendDialog />
        </div>
      </header>

      <main className="flex-1 overflow-y-auto p-4 space-y-6">
        <FriendRequests />

        {requests.length > 0 && <Separator />}

        {friends.length === 0 ? (
          <div className="flex h-full items-center justify-center">
            <p className="text-center text-sm text-muted-foreground">
              You do not have any friends yet.
            </p>
          </div>
        ) : (
          <>
            {onlineFriends.length > 0 && (
              <section className="space-y-2">
                <h3 className="text-sm font-semibold">
                  Online ({onlineFriends.length})
                </h3>

                {onlineFriends.map((friend) => (
                  <FriendCard
                    key={friend.id}
                    friend={friend}
                    onClick={chatActions.open}
                  />
                ))}
              </section>
            )}

            {offlineFriends.length > 0 && (
              <section className="space-y-2">
                <h3 className="text-sm font-semibold text-muted-foreground">
                  Offline ({offlineFriends.length})
                </h3>

                {offlineFriends.map((friend) => (
                  <FriendCard
                    key={friend.id}
                    friend={friend}
                    onClick={chatActions.open}
                  />
                ))}
              </section>
            )}
          </>
        )}
      </main>
    </aside>
  );
}
