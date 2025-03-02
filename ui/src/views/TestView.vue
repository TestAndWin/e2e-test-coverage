<template>
  <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">
          A comprehensive overview of all tests conducted on the product, including their status over the last 28 days.
          This also includes tests that have not been assigned to a specific area or feature. Only the most recent test
          for each suite name will be displayed.
        </p>
      </div>
    </div>
  </div>
  <div class="">
    <TestCoverage :productId="1" :selectedComponentProp="selectedComponent" />
  </div>
</template>

<script setup lang="ts">
import TestCoverage from '@/components/TestCoverage.vue';
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const selectedComponent = ref<string | null>(null);

onMounted(() => {
  document.title = 'e2e coverage - Tests';
  
  // Check for component query parameter
  const componentParam = route.query.component;
  if (componentParam && typeof componentParam === 'string') {
    selectedComponent.value = componentParam;
  }
});
</script>
