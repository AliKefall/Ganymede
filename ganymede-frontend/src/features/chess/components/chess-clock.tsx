"use client";

interface ChessClockProps {
  color: "white" | "black";
}

function formatTime(ms: number) {
  const totalSeconds = Math.max(0, Math.floor(ms / 1000));

  const minutes = Math.floor(totalSeconds / 60);

  const seconds = totalSeconds % 60;

  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
}

export function ChessClock({color}: ChessClockProps) {
    const remaining =
        color === "white"
    ? 10 * 60 * 1000
    : 10 * 60 * 1000

    const active = color === "white";

    return (
    <div
      className={`flex items-center justify-between rounded-xl border px-5 py-4 transition-colors ${
        active
          ? "border-primary bg-primary/10"
          : "bg-card"
      }`}
    >
      <span className="text-sm font-medium capitalize">
        {color}
      </span>

      <span className="font-mono text-3xl font-bold tabular-nums">
        {formatTime(remaining)}
      </span>
    </div>
  );
}
