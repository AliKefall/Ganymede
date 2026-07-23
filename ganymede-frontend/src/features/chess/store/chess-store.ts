"use client";

import { Chess } from "chess.js";
import { create } from "zustand";

import type {
  CapturedPieces,
  ChessClock,
  ChessGameState,
  ChessPlayer,
  GameResult,
  GameStatus,
  Move,
  PlayerColor,
  PromotionState,
} from "./types";

const INITIAL_FEN = new Chess().fen();

interface ChessStore extends ChessGameState {
  reset(): void;

  replaceGame(game: Partial<ChessGameState>): void;

  loadPosition(fen: string): void;

  clearMoves(): void;

  setGameID(id: string | null): void;

  setStatus(status: GameStatus): void;

  setResult(result: GameResult): void;

  setFen(fen: string): void;

  setTurn(turn: PlayerColor): void;

  setPlayers(
    white: ChessPlayer | null,
    black: ChessPlayer | null,
  ): void;

  setClocks(clocks: ChessClock): void;

  setLastMove(from: string, to: string): void;

  setMoves(moves: Move[]): void;

  addMove(move: Move): void;

  setCaptured(captured: CapturedPieces): void;

  setPromotion(promotion: PromotionState | null): void;

  setFlipped(flipped: boolean): void;
}

export const useChessStore = create<ChessStore>((set) => ({
  gameID: null,

  status: "idle",

  result: null,

  fen: INITIAL_FEN,

  turn: "white",

  white: null,

  black: null,

  clocks: {
    white: 0,
    black: 0,
    increment: 0,
  },

  moves: [],

  captured: {
    white: [],
    black: [],
  },

  promotion: null,

  flipped: false,

  lastMove: undefined,

  reset() {
    set({
      gameID: null,

      status: "idle",

      result: null,

      fen: INITIAL_FEN,

      turn: "white",

      white: null,

      black: null,

      clocks: {
        white: 0,
        black: 0,
        increment: 0,
      },

      moves: [],

      captured: {
        white: [],
        black: [],
      },

      promotion: null,

      flipped: false,

      lastMove: undefined,
    });
  },

  replaceGame(game) {
    set((state) => ({
      ...state,
      ...game,
    }));
  },

  loadPosition(fen) {
    set({ fen });
  },

  clearMoves() {
    set({
      moves: [],
    });
  },

  setGameID(gameID) {
    set({ gameID });
  },

  setStatus(status) {
    set({ status });
  },

  setResult(result) {
    set({ result });
  },

  setFen(fen) {
    set({ fen });
  },

  setTurn(turn) {
    set({ turn });
  },

  setPlayers(white, black) {
    set({
      white,
      black,
    });
  },

  setClocks(clocks) {
    set({ clocks });
  },

  setLastMove(from, to) {
    set({
      lastMove: {
        from,
        to,
      },
    });
  },

  setMoves(moves) {
    set({ moves });
  },

  addMove(move) {
    set((state) => {
      if (state.moves.some((m) => m.id === move.id)) {
        return state;
      }

      return {
        moves: [...state.moves, move],
      };
    });
  },

  setCaptured(captured) {
    set({ captured });
  },

  setPromotion(promotion) {
    set({ promotion });
  },

  setFlipped(flipped) {
    set({ flipped });
  },
}));
