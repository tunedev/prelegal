<script>
  import { parseSSEBlock } from '$lib/sse.js';

  let { onFormData } = $props();

  let messages = $state([
    {
      role: 'assistant',
      content: "Hi! I'll help you put together a Mutual NDA. Let's start with you — what's your name, title, and company?",
    },
  ]);
  let input = $state('');
  let loading = $state(false);
  let extracting = $state(false);

  async function sendMessage() {
    const text = input.trim();
    if (!text || loading) return;

    messages.push({ role: 'user', content: text });
    input = '';
    loading = true;
    extracting = false;

    const history = messages.map((m) => ({ role: m.role, content: m.content }));
    messages.push({ role: 'assistant', content: '' });
    const assistantIndex = messages.length - 1;

    try {
      const res = await fetch('/api/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ messages: history }),
      });
      if (!res.ok) {
        throw new Error(`Chat request failed with status ${res.status}`);
      }

      const reader = res.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const blocks = buffer.split('\n\n');
        buffer = blocks.pop();

        for (const block of blocks) {
          if (!block.trim()) continue;
          const { event, data } = parseSSEBlock(block);
          if (data === null) continue;

          if (event === 'message') {
            messages[assistantIndex].content += JSON.parse(data);
          } else if (event === 'replyDone') {
            extracting = true;
          } else if (event === 'formData') {
            onFormData(JSON.parse(data));
          } else if (event === 'error') {
            showError(assistantIndex);
          }
        }
      }
    } catch {
      showError(assistantIndex);
    } finally {
      loading = false;
      extracting = false;
    }
  }

  function showError(assistantIndex) {
    const prefix = messages[assistantIndex].content ? messages[assistantIndex].content + '\n\n' : '';
    messages[assistantIndex].content = prefix + 'Sorry, something went wrong. Please try again.';
  }

  function handleKeydown(event) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      sendMessage();
    }
  }
</script>

<div class="flex flex-col h-full">
  <div class="flex-1 overflow-y-auto space-y-3 p-4">
    {#each messages as message, i}
      <div class="flex {message.role === 'user' ? 'justify-end' : 'justify-start'}">
        <div
          class="max-w-[85%] rounded-lg px-3 py-2 text-sm whitespace-pre-wrap {message.role === 'user'
            ? 'bg-secondary text-white'
            : 'bg-gray-100 text-gray-900'}"
        >
          {#if message.role === 'assistant' && message.content === '' && loading && i === messages.length - 1}
            <span class="inline-flex gap-1 py-1" aria-label="Thinking">
              <span class="thinking-dot"></span>
              <span class="thinking-dot"></span>
              <span class="thinking-dot"></span>
            </span>
          {:else}
            {message.content}
          {/if}
        </div>
      </div>
    {/each}
    {#if extracting}
      <p class="text-xs text-graytext px-1">Updating document details…</p>
    {/if}
  </div>

  <div class="border-t p-3 flex items-end gap-2">
    <textarea
      bind:value={input}
      onkeydown={handleKeydown}
      placeholder="Type your message…"
      rows="2"
      disabled={loading}
      class="flex-1 rounded-md border border-gray-300 px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-primary disabled:bg-gray-50"
    ></textarea>
    <button
      onclick={sendMessage}
      disabled={loading}
      class="px-4 py-2 bg-secondary text-white text-sm font-medium rounded-md hover:opacity-90 transition-opacity disabled:opacity-50"
    >
      Send
    </button>
  </div>
</div>

<style>
  .thinking-dot {
    width: 6px;
    height: 6px;
    border-radius: 9999px;
    background: currentColor;
    animation: thinking-bounce 1.2s infinite ease-in-out;
  }
  .thinking-dot:nth-child(2) {
    animation-delay: 0.15s;
  }
  .thinking-dot:nth-child(3) {
    animation-delay: 0.3s;
  }

  @keyframes thinking-bounce {
    0%,
    80%,
    100% {
      opacity: 0.3;
      transform: translateY(0);
    }
    40% {
      opacity: 1;
      transform: translateY(-3px);
    }
  }
</style>
