export function parseSSEBlock(block) {
  let event = 'message';
  let data = null;

  for (const line of block.split('\n')) {
    if (line.startsWith('event: ')) {
      event = line.slice('event: '.length);
    } else if (line.startsWith('data: ')) {
      data = line.slice('data: '.length);
    }
  }

  return { event, data };
}
