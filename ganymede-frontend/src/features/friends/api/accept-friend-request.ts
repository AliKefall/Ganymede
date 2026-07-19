import { http } from "@/lib/http";

export interface AcceptFriendRequestBody {
  username: string;
}

export async function acceptFriendRequest(body: AcceptFriendRequestBody) {
  const { data } = await http.post("/friends/requests/accept", body);
  return data;
}
