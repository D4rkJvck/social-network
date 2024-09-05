import Picker, { EmojiClickData } from 'emoji-picker-react'

// import React from 'react';
import { Textarea } from '../ui/textarea';
import { Button } from '../ui/button';
import { useState } from 'react';

const ChatForm = () => {
    const [isDisplayed, setIsDisplayed] = useState(false);
    const [inputContent, setInputContent] = useState('');

    const addEmoji = (choice: EmojiClickData) => {
        setInputContent(prev => prev + choice.emoji)
        setIsDisplayed(false)
    }

    return (
        <form onSubmit={() => {
            // TODO: UseWS()
        }}>
            <img
                src='https://icons.getbootstrap.com/assets/icons/emoji-smile.svg'
                onClick={() => setIsDisplayed(val => !val)}
            />
            {isDisplayed && <Picker onEmojiClick={addEmoji} />}
            <Textarea
                onChange={event => setInputContent(event.target.value)}>
                {inputContent}
            </Textarea>
            <Button type="submit">Send</Button>
        </form>
    )
}

export default ChatForm
