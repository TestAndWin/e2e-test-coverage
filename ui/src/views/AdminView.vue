<template>
  <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">Maintaining user.</p>
      </div>
    </div>
  </div>
  <div class="">
    <UserAdmin />
  </div>
</template>

<script setup lang="ts">
import UserAdmin from '@/components/UserAdmin.vue';
import { onMounted } from 'vue';

onMounted(() => {
  document.title = 'e2e coverage - Admin';
});
</script>
