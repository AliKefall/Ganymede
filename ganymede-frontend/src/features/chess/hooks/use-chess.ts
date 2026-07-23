"use client";

import { useMemo } from "react";

import { Chess } from "chess.js";

import { chessActions } from "../store/actions";

import {
  useBoardState,
  useCaptured,
  useClocks,
  useGameID,
  useMoves,
  usePlayers,
  useResult,
  useStatus,
} from "../store/selectors";

export function useChess() {
  const board = useBoardState();

  const gameID = useGameID();

  const status = useStatus();

  const result = useResult();

  const players = usePlayers();

  const clocks = useClocks();

  const moves = useMoves();

  const captured = useCaptured();

  /**
   * FEN her değiştiğinde yeni oyun oluşturulur.
   * FEN bizim tek source of truth'umuzdur.
   */
  const chess = useMemo(() => {
    return new Chess(board.fen);
  }, [board.fen]);

  const isCheck = chess.inCheck();

  const isCheckmate = chess.isCheckmate();

  const isDraw = chess.isDraw();

  const isGameOver = chess.isGameOver();

  const legalMoves = useMemo(() => {
    return chess.moves({
      verbose: true,
    });
  }, [chess]);

  return {
    chess,

    gameID,

    status,

    result,

    players,

    clocks,

    moves,

    captured,

    ...board,

    isCheck,

    isCheckmate,

    isDraw,

    isGameOver,

    legalMoves,

    reset: chessActions.reset,

    replaceGame: chessActions.replaceGame,

    loadPosition: chessActions.loadPosition,

    setStatus: chessActions.setStatus,

    setResult: chessActions.setResult,

    setPlayers: chessActions.setPlayers,

    setClocks: chessActions.setClocks,

    addMove: chessActions.addMove,

    clearMoves: chessActions.clearMoves,

    setPromotion: chessActions.setPromotion,
  };
}
