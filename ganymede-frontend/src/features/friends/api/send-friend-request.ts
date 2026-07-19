import { http } from "@/lib/http"

export interface SendFriendRequestBody{
    username: string
}

export async function sendFriendRequest(
    body: SendFriendRequestBody,
){
    const {data} = await http.post("/friends/requests", body);

    return data
}
