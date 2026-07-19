import { useMutation } from "@tanstack/react-query";
import { logout } from "../api/auth.service";
import { useAuthStore } from "../auth-store";
import { websocketManager } from "@/lib/websocket";

export function useLogout() {
  const clearSession = useAuthStore((state) => state.clearSession);

  return useMutation({
    mutationFn: logout,
    onSuccess: () => {
        websocketManager.disconnect();
      clearSession();
    },
    onError: () => {
      websocketManager.disconnect()
        clearSession();

    },
  });
}
