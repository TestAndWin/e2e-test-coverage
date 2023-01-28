<template>
  <div class="container">
    <div v-if="loading" class="info spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <div v-for="feature in features" :key="feature['id']" class="feature shadow p-2 mb-2 rounded">
      <div :id="`feature-${feature['id']}`" class="row">
        <div class="col-5 mb-2">
          <h5 @click="showTests(feature['id'])" class="pointer">{{ feature['name'] }}</h5>
          <span v-if="feature['business-value'] == 'low'">&dollar;</span>
          <span v-if="feature['business-value'] == 'medium'">&dollar;&dollar;</span>
          <span v-if="feature['business-value'] == 'high'">&dollar;&dollar;&dollar;</span>&nbsp;
          <a v-if="feature['documentation']" v-bind:href="feature['documentation']" target="_blank"><i class="bi bi-file-text pointer" style="color: #2c3e50"></i></a>&nbsp;
          <a v-if="feature['url']" v-bind:href="feature['url']" target="_blank"><i class="bi bi-box-arrow-up-right pointer" style="color: #2c3e50"></i></a>
        </div>
        <div class="col-5">
          <span v-if="feature['total'] < 1" class="result failures">{{ feature['total'] }}</span>
          <span v-if="feature['total'] > 0" class="result total">
            {{ feature['total'] }}
            <i v-if="feature['total'] > feature['first-total']" class="bi bi-caret-up"></i>
            <i v-if="feature['total'] < feature['first-total']" class="bi bi-caret-down"></i>
          </span>
          &nbsp;
          <span class="result passes">{{ feature['passes'] }}</span> &nbsp; <span class="result failures">{{ feature['failures'] }}</span> &nbsp;
          <span class="result pending">{{ feature['pending'] }}</span> &nbsp;
          <span class="result skipped">{{ feature['skipped'] }}</span>
        </div>
        <div class="col">&nbsp;</div>
      </div>
      <TestCoverage v-if="featureToggle[feature['id']]" :featureId="feature['id']" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, ref, onMounted } from 'vue';
import TestCoverage from '@/components/TestCoverage.vue';
import http from '@/common-http';

const props = defineProps({
  areaId: Number,
});
const emit = defineEmits(['showAlert']);
const loading = ref(true);

const features = ref([]);
const featureToggle = ref([false]);
const getFeatures = async () => {
  loading.value = true;
  await http
    .get(`/api/v1/coverage/areas/${props.areaId}/features`)
    .then((response) => {
      features.value = response.data;
      featureToggle.value = new Array(features.value.length).fill(false);
    })
    .catch((err) => {
      emit('showAlert', err + ' | ' + err.response.data.error);
    });
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
