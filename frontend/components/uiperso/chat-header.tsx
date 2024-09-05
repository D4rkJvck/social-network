import { User } from '@/models/user.model';
import React from 'react';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Button } from '../ui/button';

type ChatHeaderProps = Pick<User, 'profilePicture' | 'nickname' | 'firstname'>

const ChatHeader: React.FC<ChatHeaderProps> = ({
    profilePicture,
    nickname,
    firstname,
}: ChatHeaderProps) => <>
        <Avatar>
            <AvatarImage src={profilePicture} alt={nickname || firstname} />
            <AvatarFallback>{firstname.charAt(0).toUpperCase()}</AvatarFallback>
        </Avatar>
        <span>{nickname || firstname}</span>
        <Button onClick={() => {
            // TODO: setIsChatboxOpen(!isChatboxOpen)
        }}>X</Button>
    </>

export default ChatHeader
