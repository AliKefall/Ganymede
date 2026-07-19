import { http } from "@/lib/http";

export interface MessageDTO {
    id: string;
    sender_id: string;
    recipient_id: string;
    content: string;
    created_at: string;
}

export async function getMessages(friendID: string){
    const {data} = await http.get<MessageDTO[]>(
        `/messages/${friendID}`,
    )

    return data
}
