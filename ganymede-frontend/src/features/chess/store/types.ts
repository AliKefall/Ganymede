export type PlayerColor = "white" | "black";

export type Piece =
  | "p"
  | "n"
  | "b"
  | "r"
  | "q"
  | "k";

export type PromotionPiece =
  | "q"
  | "r"
  | "b"
  | "n";

export type GameStatus =
  | "idle"
  | "waiting"
  | "starting"
  | "playing"
  | "finished"
  | "aborted";

export type GameResult =
  | "white_win"
  | "black_win"
  | "draw"
  | null;

export interface ChessPlayer {
  id: string;
  username: string;
  rating: number;
  color: PlayerColor;
  connected: boolean;
}

export interface ChessClock {
  white: number;
  black: number;
  increment: number;
}

export interface LastMove {
  from: string;
  to: string;
}

export interface Move {
  id: string;

  from: string;
  to: string;

  piece: Piece;

  san: string;

  fen: string;

  createdAt: string;

  promotion?: PromotionPiece;

  captured?: Piece;

  check?: boolean;

  checkmate?: boolean;
}

export interface CapturedPieces {
  white: Piece[];
  black: Piece[];
}

export interface PromotionState {
  isOpen: boolean;

  from: string;

  to: string;

  color: PlayerColor;
}

export interface ChessGameState {
  gameID: string | null;

  status: GameStatus;

  result: GameResult;

  fen: string;

  turn: PlayerColor;

  lastMove?: LastMove;

  white: ChessPlayer | null;

  black: ChessPlayer | null;

  clocks: ChessClock;

  moves: Move[];

  captured: CapturedPieces;

  promotion: PromotionState | null;

  flipped: boolean;
}
