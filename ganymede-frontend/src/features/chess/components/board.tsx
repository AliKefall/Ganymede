"use client";

import { Chessboard } from "react-chessboard";

import { useChess } from "../hooks/use-chess";

export function Board() {
  const {
    fen,
    flipped,
    lastMove,
    makeMove,
  } = useChess();

  return (
    <div className="w-full max-w-[720px] aspect-square">
      <Chessboard
        options={{
          position: fen,

          boardOrientation: flipped ? "black" : "white",

          animationDurationInMs: 180,

          allowDragging: true,

          showNotation: true,

          onPieceDrop: ({ sourceSquare, targetSquare }) => {
            if (!targetSquare) {
              return false;
            }

            return !!makeMove(sourceSquare, targetSquare);
          },

          squareStyles:
            lastMove === undefined
              ? {}
              : {
                  [lastMove.from]: {
                    backgroundColor: "rgba(255,255,0,.35)",
                  },

                  [lastMove.to]: {
                    backgroundColor: "rgba(255,255,0,.35)",
                  },
                },
        }}
      />
    </div>
  );
}
