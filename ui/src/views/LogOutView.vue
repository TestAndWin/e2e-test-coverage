<template>
  <div class=""></div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
onMounted(() => {
  sessionStorage.clear();
  localStorage.removeItem('roles_backup'); // Clear the backup roles from localStorage
  deleteAllCookies();

  // Reload page to update the menu
  location.assign('/');
});

const deleteAllCookies = () => {
  const cookies = document.cookie.split(';');
  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i];
    const name = cookie.split('=')[0];
    document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
  }
};
</script>
