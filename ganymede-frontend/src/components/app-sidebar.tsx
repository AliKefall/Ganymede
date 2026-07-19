"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";

import { Swords, UserCircle, LogOut } from "lucide-react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
} from "@/components/ui/sidebar";

import { cn } from "@/lib/utils";
import { useLogout } from "@/features/auth/hooks/use-logout";
import { Card } from "./ui/card";
import { useAuthStore } from "@/features/auth/auth-store";

const navItems = [
  {
    href: "/matchmaking",
    label: "Chess",
    icon: Swords,
  },
  {
    href: "/profile",
    label: "Profile",
    icon: UserCircle,
  },
];

export function AppSidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const logoutMutation = useLogout();

  async function handleLogout() {
    try {
      await logoutMutation.mutateAsync();

      router.push("/login");
    } catch {
      /*
        store zaten hook içinde temizleniyor
      */

      router.push("/login");
    }
  }

  return (
    <Sidebar>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu className="top-1.5 gap-1">
              {navItems.map(({ href, label, icon: Icon }) => {
                const isActive = pathname
                  ? pathname === href || pathname.startsWith(href + "/")
                  : false;

                return (
                  <SidebarMenuItem key={href}>
                    <SidebarMenuButton
                      asChild
                      className={cn(
                        "transition-colors duration-200 rounded-lg",

                        isActive
                          ? "bg-primary/10 text-primary font-medium"
                          : "hover:bg-muted/70",
                      )}
                    >
                      <Link href={href} className="flex items-center gap-2.5">
                        <Icon className="h-4 w-4" />
                        <span>{label}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                );
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <Card className="p-3">
              <div className="flex items-center gap-3">
                <div className="flex flex-col overflow-hidden">
                  <span className="truncate font-medium">{user?.username}</span>

                  <span className="text-xs text-muted-foreground">
                    View profile
                  </span>
                </div>
              </div>
            </Card>
            <SidebarMenuButton
              onClick={handleLogout}
              disabled={logoutMutation.isPending}
              className="text-destructive
                         hover:bg-destructive/10
                         hover:text-destructive
                         transition-colors
                         duration-200
                         rounded-lg"
            >
              <LogOut className="h-4 w-4" />

              <span>
                {logoutMutation.isPending ? "Logging out..." : "Logout"}
              </span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}
