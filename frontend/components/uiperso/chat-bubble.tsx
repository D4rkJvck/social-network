import React from 'react';
import { Message } from "@/models/message"

// Extract needed properties from  the message model.
type ChatBubbleProps = Pick<Message, 'sender' | 'content' | 'sentAt'>

const ChatBubble: React.FC<ChatBubbleProps> = ({
    sender,
    content = 'No Message',
    sentAt = new Date().toISOString().slice(0, 10),
}: ChatBubbleProps) => (
    <fieldset>
        <legend>{sender?.nickname || sender?.firstname}</legend>
        <p>{content}</p>
        <h6>{sentAt}</h6>
    </fieldset>
)

export default ChatBubble
