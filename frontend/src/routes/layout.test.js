import { render, screen, fireEvent } from '@testing-library/svelte';
import { beforeEach, describe, expect, it } from 'vitest';
import { tick, createRawSnippet } from 'svelte';
import { goto } from '$app/navigation';
import { page } from '$app/state';
import { isLoggedIn, login, logout } from '$lib/auth.js';
import Layout from './+layout.svelte';

function setPath(path) {
  page.url = new URL(`http://localhost${path}`);
}

const childrenSnippet = createRawSnippet(() => ({
  render: () => '<div data-testid="child">child content</div>',
}));

describe('root layout auth guard', () => {
  beforeEach(() => {
    logout();
    goto.mockClear();
    setPath('/');
  });

  it('redirects to /login when not logged in on a protected route', async () => {
    render(Layout, { props: { children: childrenSnippet } });
    await tick();

    expect(goto).toHaveBeenCalledWith('/login');
  });

  it('does not redirect when logged in', async () => {
    login();
    render(Layout, { props: { children: childrenSnippet } });
    await tick();

    expect(goto).not.toHaveBeenCalled();
  });

  it('does not redirect while already on /login', async () => {
    setPath('/login');
    render(Layout, { props: { children: childrenSnippet } });
    await tick();

    expect(goto).not.toHaveBeenCalled();
  });

  it('shows a logout link when logged in, and logs out on click', async () => {
    login();
    render(Layout, { props: { children: childrenSnippet } });
    await tick();

    const logoutBtn = screen.getByRole('button', { name: /log out/i });
    await fireEvent.click(logoutBtn);

    expect(isLoggedIn()).toBe(false);
    expect(goto).toHaveBeenCalledWith('/login');
  });

  it('does not show a logout link when on /login', async () => {
    setPath('/login');
    render(Layout, { props: { children: childrenSnippet } });
    await tick();

    expect(screen.queryByRole('button', { name: /log out/i })).not.toBeInTheDocument();
  });
});
