import React, { useEffect, useState, useCallback } from 'react';
import { useSubscription, useMutation, useQuery } from '@apollo/client/react';
import MESSAGE_CREATED from './graphql/subscriptions/messageCreated';
import CREATE_MESSAGE from './graphql/mutations/createMessage';
import MESSAGES from './graphql/query/messages';

type Message = {
  id: string
  message: string
}
type MessageSubscription = {
  messageCreated: Message
}
type MessagesQuery = {
  messages: Message[]
}

export const Component: React.FC = () => {
  const { data } = useSubscription<MessageSubscription>(MESSAGE_CREATED);
  const [createMessage] = useMutation<Message>(CREATE_MESSAGE);
  const queryResult = useQuery<MessagesQuery>(MESSAGES)
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState('')

  useEffect(() => {
    if (data?.messageCreated?.message) setMessages(m => [...m, data?.messageCreated])
  }, [data?.messageCreated])

  useEffect(() => {
    if (queryResult.data?.messages) setMessages(queryResult.data.messages)
  }, [queryResult.data?.messages])

  const handleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value)
  }, [])

  const handleClick = useCallback(async (e: any) => {
    e.preventDefault()
    await createMessage({ variables: { message: inputValue } })
  }, [inputValue, createMessage])

  return (
    <div style={{ maxWidth: '1180px', minHeight: '100vh', margin: '0 auto', padding: '1rem', display: 'flex', justifyContent: 'center' }}>
      <div style={{ marginTop: '3rem', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
        <div style={{ display: 'flex' }}>
          <input
            placeholder="enter message"
            style={{ width: '400px', marginRight: '0.75rem', padding: '0.5rem' }}
            value={inputValue}
            onChange={handleChange}
          />
          <button
            onClick={handleClick}
            style={{ padding: '0.5rem 1rem', backgroundColor: '#4299E1', color: 'white', border: 'none', borderRadius: '0.25rem' }}
          >
            Submit
          </button>
        </div>

        <div style={{ padding: '1rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
          {messages.map((m) => <div key={m.id}>{m.message}</div>)}
        </div>
      </div>
    </div>
  );
}
