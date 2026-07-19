import { create } from "zustand";
import type { Friend } from "./types";

interface ChatState {
    isOpen: boolean;
    selectedFriend: Friend | null;

    open(friend: Friend): void;
    close(): void;
}

export const useChatStore = create<ChatState>((set) => ({
    isOpen: false,
    selectedFriend: null,

    open(friend) {
        set({
            isOpen: true,
            selectedFriend: friend,
        });
    },

    close() {
        set({
            isOpen: false,
            selectedFriend: null,
        });
    },
}));
