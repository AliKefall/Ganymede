"use client";

import React from "react";

import { Friend } from "../store/types";

interface FriendCardProps {
  friend: Friend;
  onClick?: (friend: Friend) => void;
  rightSlot?: React.ReactNode;
}

export function FriendCard({
  friend,
  onClick,
  rightSlot,
}: FriendCardProps) {
  return (
    <button
      type="button"
      onClick={() => onClick?.(friend)}
      className="
        flex
        w-full
        items-center
        justify-between
        rounded-lg
        border
        px-4
        py-3
        transition-colors
        hover:bg-muted
        focus:outline-none
        focus:ring-2
        focus:ring-primary/40
      "
    >
      <div className="flex items-center gap-3">
        <div
          className={`h-3 w-3 rounded-full ${
            friend.online ? "bg-green-500" : "bg-zinc-500"
          }`}
        />

        <div className="flex flex-col items-start">
          <span className="font-medium">
            {friend.username}
          </span>

          <span className="text-xs text-muted-foreground">
            {friend.online ? "Online" : "Offline"}
          </span>
        </div>
      </div>

      {rightSlot}
    </button>
  );
}
