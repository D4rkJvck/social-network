import React from 'react';
import ChatBubble from './chat-bubble';
import { Message } from '@/models/message';
import { User } from '@/models/user.model';
import ChatHeader from './chat-header';
import ChatForm from './chat-form';

// const router = useRouter()

interface ChatboxProps {
    receiver: User;
    messages: Message[];
    opened: boolean;
}

const Chatbox: React.FC<ChatboxProps> = ({
    receiver,
    messages,
    opened,
}: ChatboxProps) => !opened
        ? null
        : <>
            <ChatHeader
                profilePicture={receiver.profilePicture}
                nickname={receiver.nickname}
                firstname={receiver.firstname}
            />
            <div>
                {messages?.map(message => (
                    <ChatBubble
                        key={message.id}
                        sender={message.sender}
                        content={message.content}
                        sentAt={message.sentAt}
                    />
                ))}
            </div>

            <ChatForm />
        </>

export default Chatbox
