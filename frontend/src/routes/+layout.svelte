<script>
  import '../app.css';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { isLoggedIn, logout, shouldRedirectToLogin } from '$lib/auth.js';

  let { children } = $props();

  $effect(() => {
    if (shouldRedirectToLogin(page.url.pathname, isLoggedIn())) {
      goto('/login');
    }
  });

  function handleLogout() {
    logout();
    goto('/login');
  }
</script>

{#if page.url.pathname !== '/login' && isLoggedIn()}
  <div class="flex justify-end px-4 py-2">
    <button onclick={handleLogout} class="text-xs text-graytext hover:text-navy transition-colors">
      Log out
    </button>
  </div>
{/if}

{@render children()}
