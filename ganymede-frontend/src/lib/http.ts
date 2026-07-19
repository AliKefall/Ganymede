import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";

import { useAuthStore } from "@/features/auth/auth-store";

type ApiErrorBody = {
  error?: {
    code?: string;
    message?: string;
  };
};

export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string,
  ) {
    super(message);
  }
}

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export const http = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

const refreshClient = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
});

http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = useAuthStore.getState().accessToken;

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

let isRefreshing = false;

let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (error: unknown) => void;
}> = [];

function processQueue(error: unknown, token: string | null = null) {
  failedQueue.forEach((promise) => {
    if (error) {
      promise.reject(error);
    } else if (token) {
      promise.resolve(token);
    }
  });
  failedQueue = [];
}

http.interceptors.response.use(
  (response) => response,

  async (error: AxiosError<ApiErrorBody>) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({
            resolve: (token: string) => {
              originalRequest.headers.Authorization = `Bearer ${token}`;
              resolve(http(originalRequest));
            },
            reject,
          });
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const response = await refreshClient.post<{
          access_token: string;
        }>("/auth/refresh");
        const newAccessToken = response.data.access_token;
        const currentUser = useAuthStore.getState().user;

        if (!currentUser) {
          throw new Error("No user in auth store");
        }

        useAuthStore.getState().setSession(newAccessToken, currentUser);

        processQueue(null, newAccessToken);

        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;

        return http(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);

        useAuthStore.getState().clearSession();
        window.location.href = "/login"

        return Promise.reject(
          new ApiError(401, "session_expired", "Session expired"),
        );
      } finally {
        isRefreshing = false;
      }
    }

    const body = error.response?.data;

    return Promise.reject(
      new ApiError(
        error.response?.status ?? 0,
        body?.error?.code ?? "request_failed",
        body?.error?.message ?? error.message,
      ),
    );
  },
);
