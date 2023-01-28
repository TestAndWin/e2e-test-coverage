<template>
  <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">Number of test cases and their status of the last 28 days including exploratory testing impression per product area and feature.</p>
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
