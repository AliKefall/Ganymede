"use client";

import { useState } from "react";

import { UserPlus } from "lucide-react";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

import { useSendFriendRequest } from "../hooks/use-send-friend-request";

export function AddFriendDialog() {
  const [open, setOpen] = useState(false);
  const [username, setUsername] = useState("");

  const sendFriendRequest = useSendFriendRequest();

  async function handleSubmit(e: React.SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();

    const value = username.trim();

    if (!value) return;

    try {
      await sendFriendRequest.mutateAsync({
        username: value,
      });

      setUsername("");
      setOpen(false);
    } catch {
      /*
        Hata gösterimi mutation içinde yapılabilir
        (toast, alert vs.)
      */
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <UserPlus className="mr-2 h-4 w-4" />
          Add Friend
        </Button>
      </DialogTrigger>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Friend</DialogTitle>

          <DialogDescription>
            Enter your friend`s username to send a friend request.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            autoFocus
          />

          <DialogFooter>
            <Button
              type="submit"
              disabled={sendFriendRequest.isPending || username.trim() === ""}
            >
              {sendFriendRequest.isPending ? "Sending..." : "Send Request"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
