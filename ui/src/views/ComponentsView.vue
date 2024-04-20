<template>
  <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">Lists all components with their latest test run date.</p>
      </div>
    </div>
  </div>
  <div class="">
    <ComponentsList />
  </div>
</template>

<script setup lang="ts">
import ComponentsList from '@/components/ComponentsList.vue';
import { onMounted } from 'vue';

onMounted(() => {
  document.title = 'e2e coverage - Component List';
});
</script>
