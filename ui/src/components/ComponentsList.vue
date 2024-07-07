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
            <h5 class="area-name justify-content-between">
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

const loading = ref(false);
const error = ref('');

const components = ref([]);
const getComponents = async () => {
  loading.value = true;
  error.value = '';

  http
    .get(`/api/v1/coverage/components`, {})
    .then((response) => {
      if (response.data) {
        components.value = response.data;
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
