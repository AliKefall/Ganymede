"use client";

import { Chess } from "chess.js";
import { useChessStore } from "./chess-store";

export const chessActions = {
  reset() {
    const chess = new Chess();

    useChessStore.setState({
      chess,
      fen: chess.fen(),
      turn: "white",
      moves: [],
      lastMove: undefined,
      result: null,
      status: "idle",
      captured: {
        white: [],
        black: [],
      },
      promotion: null,
    });
  },

  loadFEN(fen: string) {
    const chess = new Chess();

    chess.load(fen);

    useChessStore.setState({
      chess,
      fen: chess.fen(),
      turn: chess.turn() === "w" ? "white" : "black",
    });
  },

  makeMove(from: string, to: string, promotion?: "q" | "r" | "b" | "n") {
    const state = useChessStore.getState();

    const chess = state.chess;

    const move = chess.move({
      from,
      to,
      promotion,
    });

    if (!move) {
      return false;
    }

    useChessStore.setState({
      fen: chess.fen(),
      turn: chess.turn() === "w" ? "white" : "black",
      lastMove: {
        from,
        to,
      },
    });
    return move;
  },

  applyServerMove(
      from: string,
      to: string,
      promotion?: "q" | "r" | "b" | "n",
  ){
      return this.makeMove(from, to, promotion);
  },

  undo(){
      const chess = useChessStore.getState().chess;

      const move = chess.undo();

      if (!move){
          return false;
      }

      useChessStore.setState({
          fen: chess.fen(),
          turn: chess.turn() === "w" ? "white" : "black",
      })

      return true
  }
};
