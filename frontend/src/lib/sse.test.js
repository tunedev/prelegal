import { describe, expect, it } from 'vitest';
import { parseSSEBlock } from './sse.js';

describe('parseSSEBlock', () => {
  it('parses an event and data line', () => {
    expect(parseSSEBlock('event: message\ndata: hello')).toEqual({ event: 'message', data: 'hello' });
  });

  it('defaults to "message" event when no event line is present', () => {
    expect(parseSSEBlock('data: hello')).toEqual({ event: 'message', data: 'hello' });
  });

  it('returns null data when there is no data line', () => {
    expect(parseSSEBlock('event: done')).toEqual({ event: 'done', data: null });
  });
});
