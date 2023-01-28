<template>
  <div v-if="loading" class="info spinner-border" role="status">
    <span class="visually-hidden">Loading...</span>
  </div>

  <div class="container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
      <button type="button" class="btn-close pointer" aria-label="Close" @click="closeAlert()"></button>
    </div>

    <div v-for="test in tests" :key="test['id']" class="test shadow p-2 mb-2 rounded">
      <div class="row">
        <div class="col-5">
          <h6 @click="showTestRuns(test['suite'], test['file-name'])" class="pointer" :set="(percentage = test['failed-test-runs'] / test['total-test-runs'])">
            {{ test['suite'] }}
          </h6>
        </div>
        <div class="col-5">
          <span class="result total">
            {{ test['total'] }}
            <i v-if="test['total'] > test['first-total']" class="bi bi-caret-up"></i>
            <i v-if="test['total'] < test['first-total']" class="bi bi-caret-down"></i>
          </span>
          &nbsp; <span class="result passes">{{ test['passes'] }}</span> &nbsp; <span class="result failures">{{ test['failures'] }}</span> &nbsp;
          <span class="result pending">{{ test['pending'] }}</span> &nbsp;
          <span class="result skipped">{{ test['skipped'] }}</span>
        </div>
        <div class="col">
          <i v-if="percentage == 0" class="bi bi-sun"></i>
          <i v-if="percentage > 0 && percentage <= 0.15" class="bi bi-cloud-sun"></i>
          <i v-if="percentage > 0.15 && percentage <= 0.3" class="bi bi-cloud"></i>
          <i v-if="percentage >= 0.3 && percentage <= 0.5" class="bi bi-cloud-rain"></i>
          <i v-if="percentage > 0.5" class="bi bi-lightning"></i>
        </div>
      </div>
      <div class="row">
        <div class="col">
          <span class="test-suite d-flex justify-content-between">File: {{ test['file-name'] }}</span>
        </div>
      </div>
      <div class="row">
        <div class="col">
          <span class="test-suite d-flex justify-content-between">Test run: {{ test['test-run'] }}</span>
          <span v-if="test['area-id'] == 0"><i>Not assigned to an area/feature</i></span>
        </div>
      </div>
    </div>

    <!-- Modal to show all tests with a graph -->
    <div class="modal fade" :id="'showTestRuns_' + featureId" tabindex="-1" aria-labelledby="showTestRunsLabel" aria-hidden="true">
      <div class="modal-dialog modal-fullscreen">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="showTestRunsLabel">{{ suite }} - {{ file }}</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="row bg-light">
            <div class="col-1">&nbsp;</div>
            <div class="col-3">Date</div>
            <div class="col">Total</div>
            <div class="col">Passes</div>
            <div class="col">Failures</div>
            <div class="col">Pending</div>
            <div class="col">Skipped</div>
          </div>
          <div v-for="tr in testRuns" :key="tr['id']" class="row">
            <div class="col-1">&nbsp;</div>
            <div class="col-3">{{ tr['test-run'] }}</div>
            <div class="col">{{ tr['total'] }}</div>
            <div class="col">{{ tr['passes'] }}</div>
            <div class="col">{{ tr['failures'] }}</div>
            <div class="col">{{ tr['pending'] }}</div>
            <div class="col">{{ tr['skipped'] }}</div>
          </div>
          <Line v-if="!loading" :data="chartData" :chart-options="chartOptions" />
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Modal } from 'bootstrap';
import http from '@/common-http';

import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend } from 'chart.js';
import { Line } from 'vue-chartjs';
ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend);

const props = defineProps({
  productId: Number,
  featureId: Number,
});

const loading = ref(true);
const error = ref('');
const percentage = ref(0.0);

const chartData = ref({
  labels: [],
  datasets: [
    {
      label: 'Total',
      backgroundColor: 'blue',
      borderColor: 'blue',
      data: [],
    },
    {
      label: 'Pass',
      backgroundColor: 'green',
      borderColor: 'green',
      data: [],
    },
    {
      label: 'Fail',
      backgroundColor: 'red',
      borderColor: 'red',
      data: [],
    },
    {
      label: 'Pending',
      backgroundColor: 'orange',
      borderColor: 'orange',
      data: [],
    },
    {
      label: 'Skipped',
      backgroundColor: 'yellow',
      borderColor: 'yellow',
      data: [],
    },
  ],
});

const chartOptions = ref({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      type: 'linear',
      min: 0,
    },
  },
});

const tests = ref([]);
const getTestsForProduct = async () => {
  loading.value = true;
  await http
    .get(`/api/v1/coverage/products/${props.productId}/tests`)
    .then((response) => {
      tests.value = response.data;
    })
    .catch((err) => {
      error.value = err;
    });
  loading.value = false;
};

const getTestsForFeature = async () => {
  loading.value = true;
  await http
    .get(`/api/v1/coverage/features/${props.featureId}/tests`)
    .then((response) => {
      tests.value = response.data;
    })
    .catch((err) => {
      error.value = err;
    });

  loading.value = false;
};

const testRuns = ref([]);
const suite = ref('');
const file = ref('');
const showTestRuns = async (s: string, f: string) => {
  loading.value = true;
  suite.value = s;
  file.value = f;

  await http
    .get(`/api/v1/tests?suite=${s}&file-name=${f}`)
    .then((response) => {
      testRuns.value = response.data;

      chartData.value.labels = [];
      for (let i = 0; i < testRuns.value.length; i++) {
        // different order
        const r = testRuns.value.length - 1 - i;
        chartData.value.labels[i] = testRuns.value[r]['test-run'];
        chartData.value.datasets[0].data[i] = testRuns.value[r]['total'];
        chartData.value.datasets[1].data[i] = testRuns.value[r]['passes'];
        chartData.value.datasets[2].data[i] = testRuns.value[r]['failures'];
        chartData.value.datasets[3].data[i] = testRuns.value[r]['pending'];
        chartData.value.datasets[4].data[i] = testRuns.value[r]['skipped'];
      }
    })
    .catch((err) => {
      error.value = err;
    });
  loading.value = false;
  new Modal('#showTestRuns_' + props.featureId).show();
};

const closeAlert = () => {
  error.value = '';
};

onMounted(() => {
  if (props.featureId) {
    getTestsForFeature();
  }
  if (props.productId) {
    getTestsForProduct();
  }
});
</script>

<style scoped>
@import '../assets/styles.css';
</style>
