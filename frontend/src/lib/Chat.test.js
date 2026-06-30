import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import Chat from './Chat.svelte';

function fakeStream(chunks) {
  const encoder = new TextEncoder();
  let i = 0;
  return new ReadableStream({
    pull(controller) {
      if (i < chunks.length) {
        controller.enqueue(encoder.encode(chunks[i++]));
      } else {
        controller.close();
      }
    },
  });
}

const sampleFormData = {
  party1: { name: 'Alice', title: '', company: '', address: '' },
  party2: { name: '', title: '', company: '', address: '' },
  effectiveDate: '',
  mndaTermType: 'expires',
  mndaTermYears: 1,
  confidentialityTermType: 'years',
  confidentialityTermYears: 3,
  purpose: '',
  governingLaw: '',
  jurisdiction: '',
  modifications: '',
};

beforeEach(() => {
  vi.restoreAllMocks();
});

describe('Chat', () => {
  it('shows an initial greeting message', () => {
    render(Chat, { props: { onFormData: vi.fn() } });
    expect(screen.getByText(/Mutual NDA/i)).toBeInTheDocument();
  });

  it('sends the conversation to /api/chat, streams the reply, and reports extracted form data', async () => {
    const onFormData = vi.fn();
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      body: fakeStream([
        `event: message\ndata: ${JSON.stringify('Hi there')}\n\n`,
        `event: formData\ndata: ${JSON.stringify(sampleFormData)}\n\n`,
      ]),
    });

    render(Chat, { props: { onFormData } });

    const input = screen.getByPlaceholderText(/type your message/i);
    await fireEvent.input(input, { target: { value: "I'm Alice" } });
    await fireEvent.click(screen.getByRole('button', { name: /send/i }));

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/chat',
      expect.objectContaining({
        method: 'POST',
        body: expect.stringContaining("I'm Alice"),
      })
    );

    await waitFor(() => {
      expect(screen.getByText('Hi there')).toBeInTheDocument();
    });

    await waitFor(() => {
      expect(onFormData).toHaveBeenCalledWith(sampleFormData);
    });
  });

  it('shows an error message and re-enables input when the backend sends an error event', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      body: fakeStream([`event: error\ndata: ${JSON.stringify('upstream failed')}\n\n`]),
    });

    render(Chat, { props: { onFormData: vi.fn() } });

    await fireEvent.input(screen.getByPlaceholderText(/type your message/i), { target: { value: 'hi' } });
    await fireEvent.click(screen.getByRole('button', { name: /send/i }));

    await waitFor(() => {
      expect(screen.getByText(/something went wrong/i)).toBeInTheDocument();
    });
    expect(screen.getByRole('button', { name: /send/i })).not.toBeDisabled();
  });

  it('shows an error message when the network request itself fails', async () => {
    global.fetch = vi.fn().mockRejectedValue(new Error('network down'));

    render(Chat, { props: { onFormData: vi.fn() } });

    await fireEvent.input(screen.getByPlaceholderText(/type your message/i), { target: { value: 'hi' } });
    await fireEvent.click(screen.getByRole('button', { name: /send/i }));

    await waitFor(() => {
      expect(screen.getByText(/something went wrong/i)).toBeInTheDocument();
    });
    expect(screen.getByRole('button', { name: /send/i })).not.toBeDisabled();
  });

  it('shows an error message when the response status is not ok', async () => {
    global.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });

    render(Chat, { props: { onFormData: vi.fn() } });

    await fireEvent.input(screen.getByPlaceholderText(/type your message/i), { target: { value: 'hi' } });
    await fireEvent.click(screen.getByRole('button', { name: /send/i }));

    await waitFor(() => {
      expect(screen.getByText(/something went wrong/i)).toBeInTheDocument();
    });
  });
});
