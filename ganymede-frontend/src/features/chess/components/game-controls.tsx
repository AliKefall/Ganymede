"use client";

import {
  Flag,
  Handshake,
  RotateCcw,
  Copy,
  ArrowLeftRight,
} from "lucide-react";

import { Button } from "@/components/ui/button";

export function GameControls() {
  return (
    <div className="flex flex-col gap-3 rounded-xl border bg-card p-4">
      <h2 className="text-sm font-semibold">Game</h2>

      <Button
        variant="outline"
        className="justify-start"
        onClick={() => {
          // chessActions.flipBoard()
        }}
      >
        <ArrowLeftRight className="mr-2 h-4 w-4" />
        Flip Board
      </Button>

      <Button
        variant="outline"
        className="justify-start"
        onClick={() => {
          // chessActions.copyPGN()
        }}
      >
        <Copy className="mr-2 h-4 w-4" />
        Copy PGN
      </Button>

      <Button
        variant="outline"
        className="justify-start"
        onClick={() => {
          // chessActions.offerDraw()
        }}
      >
        <Handshake className="mr-2 h-4 w-4" />
        Offer Draw
      </Button>

      <Button
        variant="outline"
        className="justify-start"
        onClick={() => {
          // chessActions.requestRematch()
        }}
      >
        <RotateCcw className="mr-2 h-4 w-4" />
        Rematch
      </Button>

      <Button
        variant="destructive"
        className="justify-start"
        onClick={() => {
          // chessActions.resign()
        }}
      >
        <Flag className="mr-2 h-4 w-4" />
        Resign
      </Button>
    </div>
  );
}
