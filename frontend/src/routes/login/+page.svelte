<script>
  import { goto } from '$app/navigation';
  import { login } from '$lib/auth.js';

  let email = $state('');
  let password = $state('');

  function handleSubmit(event) {
    event.preventDefault();
    login();
    goto('/');
  }
</script>

<div class="flex h-screen items-center justify-center bg-gray-50">
  <form onsubmit={handleSubmit} class="w-full max-w-sm space-y-4 rounded-lg bg-white p-8 shadow-sm border">
    <div>
      <h1 class="text-xl font-semibold text-navy">Prelegal</h1>
      <p class="text-sm text-graytext">Log in to draft your agreements</p>
    </div>

    <div>
      <label for="email" class="block text-xs text-graytext mb-1">Email</label>
      <input
        id="email"
        type="email"
        bind:value={email}
        class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
        placeholder="you@company.com"
        required
      />
    </div>

    <div>
      <label for="password" class="block text-xs text-graytext mb-1">Password</label>
      <input
        id="password"
        type="password"
        bind:value={password}
        class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
        placeholder="••••••••"
        required
      />
    </div>

    <button
      type="submit"
      class="w-full px-5 py-2 bg-secondary text-white text-sm font-medium rounded-md hover:opacity-90 transition-opacity"
    >
      Log In
    </button>
  </form>
</div>
