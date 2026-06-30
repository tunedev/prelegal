const SESSION_KEY = 'prelegal_logged_in';

export function isLoggedIn() {
  return localStorage.getItem(SESSION_KEY) === 'true';
}

export function login() {
  localStorage.setItem(SESSION_KEY, 'true');
}

export function logout() {
  localStorage.removeItem(SESSION_KEY);
}

export function shouldRedirectToLogin(pathname, loggedIn) {
  return pathname !== '/login' && !loggedIn;
}
