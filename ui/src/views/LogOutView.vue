<template>
  <div class=""></div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { clearUser } from '@/stores/user';
import http from '@/common-http';
import { useRouter } from 'vue-router';

const router = useRouter();

onMounted(async () => {
  try {
    await http.post('/api/v1/auth/logout');
  } catch (e) {
    // Fehler ignorieren, da Cookies sowieso gelöscht werden
  }
  clearUser();
  deleteAllCookies();

  router.push('/');
});

const deleteAllCookies = () => {
  const cookies = document.cookie.split(';');
  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i];
    if (!cookie) continue;
    const name = cookie.split('=')[0];
    document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
  }
};
</script>
