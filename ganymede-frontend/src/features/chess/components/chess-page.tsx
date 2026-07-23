"use client";

import { Board } from "./board";
import { ChessClock } from "./chess-clock";
import { GameControls } from "./game-controls";
import { MoveHistory } from "./move-history";
import { PlayerCard } from "./player-card";
import { PromotionModal } from "./promotion-modal";

export function ChessPage() {
  return (
    <>
      <PromotionModal />

      <div className="mx-auto flex h-full w-full max-w-7xl gap-8 p-8">
        <section className="flex flex-1 flex-col items-center justify-center gap-6">
          <PlayerCard color="black" />

          <Board />

          <PlayerCard color="white" />
        </section>

        <aside className="flex w-[340px] flex-col gap-4">
          <ChessClock color="black" />

          <MoveHistory />

          <ChessClock color="white" />

          <GameControls />
        </aside>
      </div>
    </>
  );
}
