import { beforeEach, describe, expect, it } from 'vitest';
import { isLoggedIn, login, logout, shouldRedirectToLogin } from './auth.js';

describe('auth', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('reports logged out when no session is stored', () => {
    expect(isLoggedIn()).toBe(false);
  });

  it('reports logged in after login() is called', () => {
    login();
    expect(isLoggedIn()).toBe(true);
  });

  it('reports logged out after logout() is called', () => {
    login();
    logout();
    expect(isLoggedIn()).toBe(false);
  });
});

describe('shouldRedirectToLogin', () => {
  it('redirects when not logged in and not already on /login', () => {
    expect(shouldRedirectToLogin('/', false)).toBe(true);
  });

  it('does not redirect when already on /login', () => {
    expect(shouldRedirectToLogin('/login', false)).toBe(false);
  });

  it('does not redirect when logged in', () => {
    expect(shouldRedirectToLogin('/', true)).toBe(false);
  });
});
