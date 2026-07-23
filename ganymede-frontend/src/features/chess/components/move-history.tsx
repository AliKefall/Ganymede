"use client";

export function MoveHistory() {
    const moves = [
        "e4",
        "e5"
    ]// Placeholder for now

    const rows = [];

    for (let i = 0; i < moves.length; i += 2){
        rows.push({
            moveNumber: i / 2 +1 ,
            white: moves[i],
            black: moves[i+1] ?? "",
        });
    }

    return (
    <div className="flex flex-1 flex-col rounded-xl border bg-card">
      <div className="border-b px-4 py-3">
        <h2 className="font-semibold">Moves</h2>
      </div>

      <div className="flex-1 overflow-y-auto">
        {rows.length === 0 ? (
          <div className="flex h-full items-center justify-center text-sm text-muted-foreground">
            No moves yet
          </div>
        ) : (
          <div className="divide-y">
            {rows.map((row) => (
              <div
                key={row.moveNumber}
                className="grid grid-cols-[48px_1fr_1fr] items-center px-4 py-2 text-sm"
              >
                <span className="text-muted-foreground">
                  {row.moveNumber}.
                </span>

                <button className="rounded px-2 py-1 text-left transition hover:bg-accent">
                  {row.white}
                </button>

                <button className="rounded px-2 py-1 text-left transition hover:bg-accent">
                  {row.black}
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );

}
