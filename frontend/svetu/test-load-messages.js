// Test loadMessages
async function test() {
  const { loadMessages } = await import('./src/lib/i18n/loadMessages.ts');

  const messages = await loadMessages('en', ['marketplace']);

  console.log('Keys in messages:', Object.keys(messages).slice(0, 20));
  console.log('Has marketplace key:', 'marketplace' in messages);

  if (messages.marketplace) {
    console.log(
      'marketplace keys:',
      Object.keys(messages.marketplace).slice(0, 10)
    );
  }
}

test().catch(console.error);
