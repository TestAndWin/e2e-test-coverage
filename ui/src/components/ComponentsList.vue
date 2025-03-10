<template>
  <div class="product container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
    </div>

    <div v-if="loading" class="info spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    <div v-for="c in components" :key="c['name']">
      <div class="test shadow p-2 mb-2 rounded">
        <div class="row">
          <div class="col">
            <h5 class="area-name justify-content-between pointer" @click="navigateToTests(c['name'])">
              {{ c['name'] }}
              &nbsp;
            </h5>
          </div>
          <div class="col">{{ c['test-run'] }}</div>
          <div class="col-5 mb-2">
            <TestResult :test="c" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import http from '@/common-http';
import TestResult from './TestResult.vue';
import type { Component } from '@/types';
import { useRouter } from 'vue-router';

const loading = ref(false);
const error = ref('');
const router = useRouter();

const components = ref<Component[]>([]);

const navigateToTests = (componentName: string | undefined) => {
  if (componentName) {
    router.push({
      name: 'tests',
      query: { component: componentName }
    });
  }
};
const getComponents = async () => {
  loading.value = true;
  error.value = '';

  http
    .get(`/api/v1/coverage/components`, {})
    .then((response) => {
      // Check for different response formats
      if (response.data && response.data.data) {
        // Response is in StandardResponse format with { success, data, ... }
        components.value = response.data.data;
      } else if (Array.isArray(response.data)) {
        // Response is a direct array
        components.value = response.data;
      } else {
        error.value = 'Unexpected response format';
      }
    })
    .catch((err) => {
      error.value = err + ' | ' + err.response?.data?.error;
    })
    .finally(() => {
      loading.value = false;
    });
};

onMounted(() => {
  getComponents();
});
</script>

<style scoped>
@import '../assets/styles.css';
</style>
