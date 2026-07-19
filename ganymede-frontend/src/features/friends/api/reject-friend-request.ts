import { http } from "@/lib/http";

export interface RejectFriendRequestBody{
    username: string;
}

export async function rejectFriendRequest(
    body: RejectFriendRequestBody,
){
    const {data} = await http.post("/friends/requests/reject", body)
        return data
}
