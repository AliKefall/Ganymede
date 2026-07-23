"use client";

import { useChessStore } from "./chess-store";

export const useGameID = () =>
  useChessStore((state) => state.gameID);

export const useBoardState = () =>
  useChessStore((state) => ({
    fen: state.fen,
    turn: state.turn,
    lastMove: state.lastMove,
  }));

export const usePlayers = () =>
  useChessStore((state) => ({
    white: state.white,
    black: state.black,
  }));

export const useStatus = () =>
  useChessStore((state) => state.status);

export const useResult = () =>
  useChessStore((state) => state.result);

export const useMoves = () =>
  useChessStore((state) => state.moves);

export const useCaptured = () =>
  useChessStore((state) => state.captured);

export const usePromotion = () =>
  useChessStore((state) => state.promotion);

export const useClocks = () =>
  useChessStore((state) => state.clocks);

export const useFlipped = () =>
  useChessStore((state) => state.flipped);
