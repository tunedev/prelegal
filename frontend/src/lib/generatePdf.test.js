import { describe, expect, it, vi } from 'vitest';
import { generatePdf } from './generatePdf.js';

describe('generatePdf', () => {
  it('calls window.print()', () => {
    const print = vi.spyOn(window, 'print').mockImplementation(() => {});
    generatePdf();
    expect(print).toHaveBeenCalledOnce();
    print.mockRestore();
  });
});
