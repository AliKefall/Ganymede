import { AppSidebar } from "@/components/app-sidebar";

import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";
import { Bootstrap } from "@/features/bootstrap/components/bootstrap";
import { ChatPanel } from "@/features/friends/components/chat.panel";
import { FriendsPanel } from "@/features/friends/components/friends.panel";
import { FriendsPanelButton } from "@/features/friends/components/friends.panel.button";

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SidebarProvider>
      <Bootstrap />
      <AppSidebar />
      <SidebarInset>
        {children}
        <FriendsPanelButton />

        <FriendsPanel />
        <ChatPanel />
      </SidebarInset>
    </SidebarProvider>
  );
}
