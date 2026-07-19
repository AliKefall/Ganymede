"use client";

import { Users } from "lucide-react";
import { useFriendsPanelStore } from "../store/panel-store";

export function FriendsPanelButton(){
    const toggle = useFriendsPanelStore(
        (state) => state.toggle,
    );

    return (
        <button
        onClick={toggle}
        className="
        fixed
        bottom-6
        right-6
        h-14
        w-14
        rounded-full
        shadow-lg
        bg-primary
        text-primary-foreground
        flex
        items-center
        justify-center
        transition-transform
        hover:scale-185
        active:scale-95
        "
        >
        <Users className="h-6 w-6" />
        </button>
    )
}
