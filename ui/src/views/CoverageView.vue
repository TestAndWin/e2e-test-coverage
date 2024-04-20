<template>
  <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">
          The quantity of test cases and their status over the last 28 days, including the impression of exploratory
          testing on each product area and feature.
        </p>
      </div>
    </div>
  </div>
  <div class="">
    <AreaCoverage :productId="1" />
  </div>
</template>

<script setup lang="ts">
import AreaCoverage from '@/components/AreaCoverage.vue';
import { onMounted } from 'vue';

onMounted(() => {
  document.title = 'e2e coverage - Coverage';
});
</script>
