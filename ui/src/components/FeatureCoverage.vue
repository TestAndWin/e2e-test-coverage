<template>
  <div class="container">
    <div v-if="loading" class="info spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <div v-for="feature in features" :key="feature['id']" class="feature shadow p-2 mb-2 rounded">
      <div :id="`feature-${feature['id']}`" class="row">
        <div class="col-5 mb-2">
          <h5 @click="showTests(feature['id'] ?? 0)" class="pointer">{{ feature['name'] }}</h5>
          <span v-if="feature['business-value'] == 'low'">&dollar;</span>
          <span v-if="feature['business-value'] == 'medium'">&dollar;&dollar;</span>
          <span v-if="feature['business-value'] == 'high'">&dollar;&dollar;&dollar;</span>&nbsp;
          <a v-if="feature['documentation']" v-bind:href="feature['documentation']" target="_blank"
            ><i class="bi bi-file-text pointer" style="color: #2c3e50"></i></a
          >&nbsp;
          <a v-if="feature['url']" v-bind:href="feature['url']" target="_blank"
            ><i class="bi bi-box-arrow-up-right pointer" style="color: #2c3e50"></i
          ></a>
        </div>
        <div class="col-5">
          <TestResult :test="feature" />
        </div>
        <div class="col">&nbsp;</div>
      </div>
      <TestCoverage v-if="featureToggle[feature['id'] ?? 0]" :featureId="feature['id'] ?? 0" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import TestCoverage from '@/components/TestCoverage.vue';
import TestResult from '@/components/TestResult.vue';
import http from '@/common-http';
import type { Feature } from '@/types';

const props = defineProps({
  areaId: Number
});
const emit = defineEmits(['showAlert']);
const loading = ref(true);

const features = ref<Feature[]>([]);
const featureToggle = ref([false]);
const getFeatures = async () => {
  loading.value = true;
  try {
    const response = await http.get(`/api/v1/coverage/areas/${props.areaId}/features`);

    // Extract data from StandardResponse format
    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      features.value = response.data.data;
    } else {
      features.value = [];
    }

    // Ensure all required properties have default values
    features.value = features.value.map((feature) => ({
      ...feature,
      'business-value': feature['business-value'] || 'low',
      total: feature.total || 0,
      passes: feature.passes || 0,
      failures: feature.failures || 0,
      pending: feature.pending || 0,
      skipped: feature.skipped || 0,
      'first-total': feature['first-total'] || 0,
      documentation: feature.documentation || '',
      url: feature.url || ''
    })) as Feature[];

    featureToggle.value = new Array(features.value.length).fill(false);
  } catch (err) {
    emit('showAlert', `Error loading feature data: ${err}`);
  }
  loading.value = false;
};

const showTests = (featureId: number) => {
  featureToggle.value[featureId] = !featureToggle.value[featureId];
};

onMounted(() => {
  getFeatures();
});
</script>

<style scoped>
@import '../assets/styles.css';
</style>
