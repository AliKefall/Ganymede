"use client";

import { Crown } from "lucide-react";

interface PlayerCardProps {
  color: "white" | "black";
}

export function PlayerCard({ color }: PlayerCardProps) {
  const player = {
    username: color === "white" ? "You" : "Opponent",
    rating: 1500,
    online: true,
  };

  return (
    <div className="flex w-full max-w-[720px] items-center justify-between rounded-xl border bg-card px-5 py-3">
      <div className="flex items-center gap-4">
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted text-lg font-bold">
          {player.username.charAt(0).toUpperCase()}
        </div>

        <div>
          <div className="flex items-center gap-2">
            <span className="font-semibold">{player.username}</span>

            {player.online && (
              <span className="h-2.5 w-2.5 rounded-full bg-green-500" />
            )}
          </div>

          <span className="text-sm text-muted-foreground">{player.rating}</span>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <Crown className="h-5 w-5 text-amber-500" />

        <div
          className={`h-8 w-8 rounded-full border ${
            color === "white" ? "bg-white" : "border-zinc-700 bg-black"
          }`}
        />
      </div>
    </div>
  );
}
