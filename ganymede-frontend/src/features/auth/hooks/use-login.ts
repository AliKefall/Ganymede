import { useMutation } from "@tanstack/react-query";
import { login } from "../api/auth.service";
import { useAuthStore } from "../auth-store";
import { websocketManager } from "@/lib/websocket";

export function useLogin() {
  const setSession = useAuthStore((state) => state.setSession);

  return useMutation({
    mutationFn: login,
    onSuccess: (data) => {
      setSession(data.access_token, data.user);
      websocketManager.connect(data.access_token)
    },
  });
}
