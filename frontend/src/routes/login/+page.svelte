<script lang="ts">
  import { supabase } from '$lib/supabase';

  let email = '';
  let password = '';

  async function login() {
    const { data, error } = await supabase.auth.signInWithPassword({ email, password });
    if (error) {
      alert(error.message);
      return;
    }

    const token = data.session?.access_token;
    if (token) {
      // Skicka token till din Go-backend
      const res = await fetch("http://localhost:8080/api/hello", {
        headers: { Authorization: `Bearer ${token}` }
      });
      alert(await res.text());
    }
  }
</script>

<input bind:value={email} placeholder="Email" />
<input type="password" bind:value={password} placeholder="Password" />
<button on:click={login}>Logga in</button>
