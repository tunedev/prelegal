import { render, screen, fireEvent } from '@testing-library/svelte';
import { beforeEach, describe, expect, it } from 'vitest';
import { isLoggedIn, logout } from '$lib/auth.js';
import { goto } from '$app/navigation';
import LoginPage from './+page.svelte';

describe('login page', () => {
  beforeEach(() => {
    logout();
    goto.mockClear();
  });

  it('renders email and password fields and a submit button', () => {
    render(LoginPage);

    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /log in/i })).toBeInTheDocument();
  });

  it('logs the user in and navigates to / on submit, without validating credentials', async () => {
    render(LoginPage);

    await fireEvent.input(screen.getByLabelText(/email/i), { target: { value: 'someone@example.com' } });
    await fireEvent.input(screen.getByLabelText(/password/i), { target: { value: 'anything' } });
    await fireEvent.click(screen.getByRole('button', { name: /log in/i }));

    expect(isLoggedIn()).toBe(true);
    expect(goto).toHaveBeenCalledWith('/');
  });
});
