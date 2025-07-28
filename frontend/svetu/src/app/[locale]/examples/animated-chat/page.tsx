'use client';

import React, { useState, useRef, useEffect } from 'react';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

const AnimatedChat = () => {
  const [messages, setMessages] = useState<any[]>([
    {
      id: 1,
      text: 'Привет! Товар еще в наличии?',
      sender: 'buyer',
      time: '10:23',
      reactions: ['👍'],
    },
    {
      id: 2,
      text: 'Да, конечно! iPhone в отличном состоянии 📱',
      sender: 'seller',
      time: '10:25',
      reactions: [],
    },
    {
      id: 3,
      text: 'Отлично! Можно фото коробки?',
      sender: 'buyer',
      time: '10:26',
      reactions: [],
    },
  ]);
  const [inputMessage, setInputMessage] = useState('');
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const [selectedEmoji, setSelectedEmoji] = useState<string | null>(null);
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const animatedEmojis = [
    { emoji: '❤️', animation: 'animate-bounce' },
    { emoji: '😂', animation: 'animate-spin' },
    { emoji: '🎉', animation: 'animate-ping' },
    { emoji: '👍', animation: 'animate-pulse' },
    { emoji: '🔥', animation: 'animate-bounce' },
    { emoji: '💯', animation: 'animate-spin' },
    { emoji: '😍', animation: 'animate-pulse' },
    { emoji: '🚀', animation: 'animate-bounce' },
  ];

  const quickReplies = [
    'Да, конечно!',
    'Нет, спасибо',
    'Можно скидку?',
    'Когда можно забрать?',
    'Есть доставка?',
  ];

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const sendMessage = () => {
    if (inputMessage.trim()) {
      const newMessage = {
        id: messages.length + 1,
        text: inputMessage,
        sender: 'buyer',
        time: new Date().toLocaleTimeString('ru-RU', {
          hour: '2-digit',
          minute: '2-digit',
        }),
        reactions: [],
      };
      setMessages([...messages, newMessage]);
      setInputMessage('');

      // Simulate seller typing
      setIsTyping(true);
      setTimeout(() => {
        setIsTyping(false);
        const sellerReply = {
          id: messages.length + 2,
          text: 'Спасибо за ваш интерес! Сейчас отправлю фото 📸',
          sender: 'seller',
          time: new Date().toLocaleTimeString('ru-RU', {
            hour: '2-digit',
            minute: '2-digit',
          }),
          reactions: [],
        };
        setMessages((prev) => [...prev, sellerReply]);
      }, 2000);
    }
  };

  const addReaction = (messageId: number, emoji: string) => {
    setMessages(
      messages.map((msg) =>
        msg.id === messageId
          ? { ...msg, reactions: [...msg.reactions, emoji] }
          : msg
      )
    );
    setSelectedEmoji(emoji);
    setTimeout(() => setSelectedEmoji(null), 1000);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200">
      {/* Header */}
      <div className="navbar bg-base-100 shadow-lg">
        <div className="navbar-start">
          <button className="btn btn-ghost btn-circle">
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 19l-7-7 7-7"
              />
            </svg>
          </button>
        </div>
        <div className="navbar-center">
          <div className="flex items-center gap-3">
            <div className="avatar online">
              <div className="w-10 rounded-full">
                <img
                  src="https://ui-avatars.com/api/?name=Tech+Store&background=6366f1&color=fff"
                  alt="Seller"
                />
              </div>
            </div>
            <div>
              <h1 className="font-bold">TechStore</h1>
              <p className="text-xs text-success">В сети</p>
            </div>
          </div>
        </div>
        <div className="navbar-end">
          <button className="btn btn-ghost btn-circle">
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
              />
            </svg>
          </button>
        </div>
      </div>

      <div className="container mx-auto px-4 max-w-4xl">
        <div className="flex h-[calc(100vh-64px)]">
          {/* Chat Container */}
          <div className="flex-1 flex flex-col">
            {/* Product Info Bar */}
            <AnimatedSection animation="slideDown">
              <div className="bg-base-100 p-4 shadow-md">
                <div className="flex items-center gap-4">
                  <img
                    src="/api/minio/download?fileName=listings/0a47e66f-d8da-459f-a2ba-8e2b85ae0163/38ad29e6-7b07-4bfc-9db2-d965cb6b966f.jpg"
                    alt="Product"
                    className="w-16 h-16 rounded-lg object-cover"
                  />
                  <div className="flex-1">
                    <h3 className="font-bold">iPhone 14 Pro Max 256GB</h3>
                    <p className="text-2xl font-bold text-primary">€899</p>
                  </div>
                  <button className="btn btn-primary btn-sm">Купить</button>
                </div>
              </div>
            </AnimatedSection>

            {/* Messages */}
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
              {messages.map((message, idx) => (
                <AnimatedSection
                  key={message.id}
                  animation={
                    message.sender === 'buyer' ? 'slideRight' : 'slideLeft'
                  }
                  delay={idx * 0.1}
                >
                  <div
                    className={`flex ${message.sender === 'buyer' ? 'justify-end' : 'justify-start'}`}
                  >
                    <div
                      className={`max-w-xs lg:max-w-md ${message.sender === 'buyer' ? 'order-2' : ''}`}
                    >
                      <div
                        className={`relative px-4 py-2 rounded-2xl ${
                          message.sender === 'buyer'
                            ? 'bg-primary text-primary-content rounded-br-none'
                            : 'bg-base-300 rounded-bl-none'
                        }`}
                      >
                        <p>{message.text}</p>
                        <span
                          className={`text-xs ${message.sender === 'buyer' ? 'text-primary-content/70' : 'text-base-content/60'}`}
                        >
                          {message.time}
                        </span>

                        {/* Reactions */}
                        {message.reactions.length > 0 && (
                          <div className="absolute -bottom-3 right-0 flex gap-1">
                            {message.reactions.map(
                              (reaction: string, idx: number) => (
                                <span
                                  key={idx}
                                  className={`text-lg ${selectedEmoji === reaction ? 'animate-bounce' : ''}`}
                                >
                                  {reaction}
                                </span>
                              )
                            )}
                          </div>
                        )}
                      </div>

                      {/* Add Reaction Button */}
                      <button
                        className="btn btn-ghost btn-xs mt-1"
                        onClick={() => {
                          const randomEmoji =
                            animatedEmojis[
                              Math.floor(Math.random() * animatedEmojis.length)
                            ];
                          addReaction(message.id, randomEmoji.emoji);
                        }}
                      >
                        + React
                      </button>
                    </div>
                  </div>
                </AnimatedSection>
              ))}

              {/* Typing Indicator */}
              {isTyping && (
                <AnimatedSection animation="fadeIn">
                  <div className="flex justify-start">
                    <div className="bg-base-300 rounded-2xl rounded-bl-none px-4 py-3">
                      <div className="flex gap-1">
                        <span
                          className="w-2 h-2 bg-base-content/40 rounded-full animate-bounce"
                          style={{ animationDelay: '0ms' }}
                        ></span>
                        <span
                          className="w-2 h-2 bg-base-content/40 rounded-full animate-bounce"
                          style={{ animationDelay: '150ms' }}
                        ></span>
                        <span
                          className="w-2 h-2 bg-base-content/40 rounded-full animate-bounce"
                          style={{ animationDelay: '300ms' }}
                        ></span>
                      </div>
                    </div>
                  </div>
                </AnimatedSection>
              )}

              <div ref={messagesEndRef} />
            </div>

            {/* Quick Replies */}
            <AnimatedSection animation="slideUp">
              <div className="px-4 py-2">
                <div className="flex gap-2 overflow-x-auto pb-2">
                  {quickReplies.map((reply, idx) => (
                    <button
                      key={idx}
                      className="btn btn-sm btn-outline whitespace-nowrap"
                      onClick={() => setInputMessage(reply)}
                    >
                      {reply}
                    </button>
                  ))}
                </div>
              </div>
            </AnimatedSection>

            {/* Input Area */}
            <AnimatedSection animation="slideUp" delay={0.2}>
              <div className="bg-base-100 border-t p-4">
                <div className="flex items-end gap-2">
                  <button className="btn btn-ghost btn-circle">
                    <svg
                      className="w-6 h-6"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"
                      />
                    </svg>
                  </button>

                  <div className="relative flex-1">
                    <input
                      type="text"
                      placeholder="Написать сообщение..."
                      className="input input-bordered w-full pr-10"
                      value={inputMessage}
                      onChange={(e) => setInputMessage(e.target.value)}
                      onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
                    />
                    <button
                      className="btn btn-ghost btn-sm btn-circle absolute right-1 top-1/2 -translate-y-1/2"
                      onClick={() => setShowEmojiPicker(!showEmojiPicker)}
                    >
                      😊
                    </button>
                  </div>

                  <button
                    className="btn btn-primary btn-circle"
                    onClick={sendMessage}
                  >
                    <svg
                      className="w-6 h-6"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
                      />
                    </svg>
                  </button>
                </div>

                {/* Emoji Picker */}
                {showEmojiPicker && (
                  <AnimatedSection animation="slideUp">
                    <div className="bg-base-200 rounded-lg p-4 mt-2">
                      <h4 className="font-semibold mb-2">
                        Анимированные эмодзи
                      </h4>
                      <div className="grid grid-cols-8 gap-2">
                        {animatedEmojis.map((item, idx) => (
                          <button
                            key={idx}
                            className={`text-2xl hover:scale-125 transition-transform ${item.animation}`}
                            onClick={() => {
                              setInputMessage(inputMessage + item.emoji);
                              setShowEmojiPicker(false);
                            }}
                          >
                            {item.emoji}
                          </button>
                        ))}
                      </div>
                    </div>
                  </AnimatedSection>
                )}
              </div>
            </AnimatedSection>
          </div>

          {/* Sidebar Info */}
          <AnimatedSection
            animation="slideRight"
            delay={0.3}
            className="hidden lg:block w-80 p-4"
          >
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title">💬 Фишки чата</h3>
                <div className="space-y-4">
                  <div className="flex items-start gap-3">
                    <span className="text-2xl animate-bounce">🎉</span>
                    <div>
                      <h4 className="font-semibold">Анимированные реакции</h4>
                      <p className="text-sm text-base-content/60">
                        Нажмите на сообщение чтобы добавить реакцию
                      </p>
                    </div>
                  </div>
                  <div className="flex items-start gap-3">
                    <span className="text-2xl animate-pulse">💬</span>
                    <div>
                      <h4 className="font-semibold">Быстрые ответы</h4>
                      <p className="text-sm text-base-content/60">
                        Используйте готовые фразы для быстрого ответа
                      </p>
                    </div>
                  </div>
                  <div className="flex items-start gap-3">
                    <span className="text-2xl animate-spin">😊</span>
                    <div>
                      <h4 className="font-semibold">Живые эмодзи</h4>
                      <p className="text-sm text-base-content/60">
                        Эмодзи с анимацией для ярких эмоций
                      </p>
                    </div>
                  </div>
                  <div className="flex items-start gap-3">
                    <span className="text-2xl">🔔</span>
                    <div>
                      <h4 className="font-semibold">Умные уведомления</h4>
                      <p className="text-sm text-base-content/60">
                        Получайте уведомления о важных сообщениях
                      </p>
                    </div>
                  </div>
                </div>

                <div className="divider"></div>

                <div className="bg-info/10 rounded-lg p-3">
                  <h4 className="font-semibold flex items-center gap-2">
                    <svg
                      className="w-5 h-5 text-info"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                        clipRule="evenodd"
                      />
                    </svg>
                    Безопасность
                  </h4>
                  <p className="text-sm mt-2">
                    Все сообщения шифруются и защищены. Используйте эскроу для
                    безопасных сделок.
                  </p>
                </div>
              </div>
            </div>
          </AnimatedSection>
        </div>
      </div>

      {/* Floating Emoji Animation */}
      {selectedEmoji && (
        <div className="fixed inset-0 pointer-events-none flex items-center justify-center">
          <span className="text-8xl animate-ping">{selectedEmoji}</span>
        </div>
      )}
    </div>
  );
};

export default AnimatedChat;
