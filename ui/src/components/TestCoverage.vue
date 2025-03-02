<template>
  <div v-if="loading" class="info spinner-border" role="status">
    <span class="visually-hidden">Loading...</span>
  </div>

  <div class="container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
      <button type="button" class="btn-close pointer" aria-label="Close" @click="closeAlert()"></button>
    </div>

    <div v-if="tests.length > 0" class="test shadow p-2 mb-2 rounded">
      <span><b>Filter</b>:</span> &nbsp;
      <span @click="filter(sun)"><span v-if="filterCriteria == sun">&gt;</span><i class="bi bi-sun pointer"></i></span
      >&nbsp;
      <span @click="filter(sun + 0.001)"
        ><span v-if="filterCriteria == sun + 0.001">&gt;</span><i class="bi bi-cloud-sun pointer"></i
      ></span>
      &nbsp;
      <span @click="filter(sunCloud + 0.001)"
        ><span v-if="filterCriteria == sunCloud + 0.001">&gt;</span><i class="bi bi-cloud pointer"></i
      ></span>
      &nbsp;
      <span @click="filter(cloud + 0.001)"
        ><span v-if="filterCriteria == cloud + 0.001">&gt;</span><i class="bi bi-cloud-rain pointer"></i
      ></span>
      &nbsp;
      <span @click="filter(cloudRain + 0.001)"
        ><span v-if="filterCriteria == cloudRain + 0.001">&gt;</span><i class="bi bi-lightning pointer"></i></span
      >&nbsp; <span>or worse</span>&nbsp;<b>OR</b>&nbsp;<span v-if="filterLastFailed">&gt;</span
      ><span class="pointer" @click="switchFilterLastFailed()">failed on last run</span>
      <span>&nbsp;<b>AND</b>&nbsp;Component: </span>
      <select class="component-select" v-model="selectedComponent">
        <option v-for="c in components" :key="c" :value="c">
          {{ c }}
        </option>
      </select>
    </div>

    <div v-for="test in tests" :key="test['id']">
      <template
        v-if="
          ((!filterLastFailed && (test.percent ?? 0) >= filterCriteria) ||
            (filterLastFailed && (test.failures ?? 0) > 0)) &&
          (selectedComponent == 'all' || selectedComponent == test['component'])
        "
      >
        <div class="test shadow p-2 mb-2 rounded">
          <div class="row">
            <div class="col-5">
              <h6
                @click="showTestRuns(test.component ?? '', test.suite ?? '', test['file-name'] ?? '')"
                class="pointer"
              >
                {{ test['suite'] }}
              </h6>
            </div>
            <div class="col-5">
              <TestResult :test="test" />
            </div>
            <div class="col">
              <i v-if="(test.percent ?? 0) > sun && (test.percent ?? 0) <= sunCloud" class="bi bi-cloud-sun"></i>
              <i v-if="(test.percent ?? 0) > sunCloud && (test.percent ?? 0) <= cloud" class="bi bi-cloud"></i>
              <i v-if="(test.percent ?? 0) >= cloud && (test.percent ?? 0) <= cloudRain" class="bi bi-cloud-rain"></i>
              <i v-if="(test.percent ?? 0) > cloudRain" class="bi bi-lightning"></i>
              &nbsp;&nbsp;
              <i
                @click="deleteTests(test.component ?? '', test.suite ?? '', test['file-name'] ?? '')"
                class="bi bi-trash pointer"
              />
            </div>
          </div>
          <div class="row">
            <div class="col">
              <span class="test-suite d-flex justify-content-between">Component: {{ test['component'] }}</span>
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
          <div class="row" v-if="test['url']">
            <div class="col">
              <span class="test-suite d-flex justify-content-between"
                ><a v-bind:href="test['url']" target="_blank">Test Report</a></span
              >
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- Modal to show test results with a graph -->
    <div
      class="modal fade"
      :id="'showTestRuns_' + featureId"
      tabindex="-1"
      aria-labelledby="showTestRunsLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog modal-xl modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="showTestRunsLabel">{{ suite }} - {{ file }}</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <!--
          <div class="row bg-light">
            <div class="col-1">&nbsp;</div>
            <div class="col-3">Date</div>
            <div class="col">Total</div>
            <div class="col">Passes</div>
            <div class="col">Failures</div>
            <div class="col">Pending</div>
            <div class="col">Skipped</div>
          </div>
          <div v-for="(tr, index) in testRuns" :key="index" class="row">
              <div class="col-1">&nbsp;</div>
              <div class="col-3">{{ tr['test-run'] }}</div>
              <div class="col">{{ tr['total'] }}</div>
              <div class="col">{{ tr['passes'] }}</div>
              <div class="col">{{ tr['failures'] }}</div>
              <div class="col">{{ tr['pending'] }}</div>
              <div class="col">{{ tr['skipped'] }}</div>
          </div>
          -->
          <Line v-if="!loading" :data="chartData" :chart-options="chartOptions" />
          <div class="col">(only the latest 10 are displayed)</div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { Modal } from 'bootstrap';
import TestResult from '@/components/TestResult.vue';
import http from '@/common-http';
import type { Test } from '@/types';

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js';
import { Line } from 'vue-chartjs';
ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend);

const props = defineProps({
  productId: Number,
  featureId: Number,
  selectedComponentProp: String
});

const loading = ref(true);
const error = ref('');
const sun = ref(0);
const sunCloud = ref(0.15);
const cloud = ref(0.3);
const cloudRain = ref(0.5);

const chartData = ref({
  labels: [],
  datasets: [
    /*{
      label: 'Total',
      backgroundColor: 'blue',
      borderColor: 'blue',
      data: [],
    },*/
    {
      label: 'Pass',
      backgroundColor: 'green',
      borderColor: 'green',
      data: []
    },
    {
      label: 'Fail',
      backgroundColor: 'red',
      borderColor: 'red',
      data: []
    },
    {
      label: 'Pending',
      backgroundColor: 'orange',
      borderColor: 'orange',
      data: []
    },
    {
      label: 'Skipped',
      backgroundColor: 'yellow',
      borderColor: 'yellow',
      data: []
    }
  ]
});

const chartOptions = ref({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      type: 'linear',
      beginAtZero: true,
      min: 0,
      ticks: {
        stepSize: 1,
        precision: 0
      }
    }
  }
});

const tests = ref<Test[]>([]);
const getTestsForProduct = async () => {
  loading.value = true;
  try {
    const response = await http.get(`/api/v1/coverage/products/${props.productId}/tests`);

    // Extract data from StandardResponse format
    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      tests.value = response.data.data;
    } else {
      tests.value = [];
    }

    // Ensure tests.value is always an array
    if (!Array.isArray(tests.value)) {
      tests.value = [];
    }

    if (tests.value.length > 0) {
      calculatePercentage();
      setComponentsList();
    } else {
      components.value = ['all']; // Default component list
    }
  } catch (err) {
    error.value = `Error loading tests: ${err}`;
    tests.value = []; // Initialize as empty array on error
  }
  loading.value = false;
};

const getTestsForFeature = async () => {
  loading.value = true;
  try {
    const response = await http.get(`/api/v1/coverage/features/${props.featureId}/tests`);

    // Extract data from StandardResponse format
    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      tests.value = response.data.data;
    } else {
      tests.value = [];
    }

    // Ensure tests.value is always an array
    if (!Array.isArray(tests.value)) {
      tests.value = [];
    }

    if (tests.value.length > 0) {
      calculatePercentage();
      setComponentsList();
    } else {
      components.value = ['all']; // Default component list
    }
  } catch (err) {
    error.value = `Error loading tests: ${err}`;
    tests.value = []; // Initialize as empty array on error
  }
  loading.value = false;
};

// Calculate the percentage of the failed tests
const calculatePercentage = () => {
  if (!Array.isArray(tests.value)) {
    return;
  }

  for (const test of tests.value) {
    if (!test) continue;

    // Ensure values are present and valid numbers
    const failedRuns = typeof test['failed-test-runs'] === 'number' ? test['failed-test-runs'] : 0;
    const totalRuns = typeof test['total-test-runs'] === 'number' ? test['total-test-runs'] : 1;

    // Avoid division by zero
    if (totalRuns === 0) {
      test.percent = 0;
    } else {
      test.percent = failedRuns / totalRuns;
    }

    // Ensure other required properties exist with defaults
    test.component = test.component || 'Unknown';
    test.suite = test.suite || '';
    test['file-name'] = test['file-name'] || '';
    test['test-run'] = test['test-run'] || '';
    test.failures = test.failures || 0;
    test.passes = test.passes || 0;
    test.pending = test.pending || 0;
    test.skipped = test.skipped || 0;
    test.total = test.total || 0;
  }
};

const components = ref<string[]>([]);
const selectedComponent = ref('all');
// Get list of components from the test results
const setComponentsList = () => {
  if (!Array.isArray(tests.value)) {
    components.value = ['all'];
    return;
  }

  try {
    // Store current selection to restore it if it exists in the new list
    const currentSelection = selectedComponent.value;

    // Extract component names from tests, filter out any undefined/null values
    const componentSet = new Set(tests.value.filter((item) => item?.component).map((item) => item.component ?? ''));

    // Add 'all' as the first option
    components.value = ['all', ...componentSet];

    // Apply selectedComponentProp if available and exists in the list
    if (props.selectedComponentProp && componentSet.has(props.selectedComponentProp)) {
      selectedComponent.value = props.selectedComponentProp;
    }
    // Otherwise restore previous selection if it still exists in the list
    else if (currentSelection !== 'all' && componentSet.has(currentSelection)) {
      selectedComponent.value = currentSelection;
    }
  } catch (err) {
    components.value = ['all']; // Default value
  }
};

const testRuns = ref([]);
const suite = ref('');
const file = ref('');
const showTestRuns = async (c: string, s: string, f: string) => {
  loading.value = true;
  suite.value = s;
  file.value = f;

  try {
    const response = await http.get(`/api/v1/tests?component=${c}&suite=${s}&file-name=${f}`);

    let testRunsData = [];

    // Extract data from StandardResponse format
    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      testRunsData = response.data.data;
    } else {
      testRunsData = [];
    }

    // We use only the latest 10
    testRuns.value = testRunsData.slice(0, 10);

    // Prepare chart data
    chartData.value.labels = [];
    chartData.value.datasets.forEach((dataset) => {
      dataset.data = [];
    });

    if (testRuns.value.length > 0) {
      for (let i = 0; i < testRuns.value.length; i++) {
        // different order (reverse chronological)
        const r = testRuns.value.length - 1 - i;
        const run = testRuns.value[r];

        if (!run) continue;

        // Format the date for the label
        chartData.value.labels[i] = run['test-run'] || `Run ${i + 1}`;

        // Set data points with safe access
        // chartData.value.datasets[0].data[i] = run['total'] || 0;
        chartData.value.datasets[0].data[i] = run['passes'] || 0;
        chartData.value.datasets[1].data[i] = run['failures'] || 0;
        chartData.value.datasets[2].data[i] = run['pending'] || 0;
        chartData.value.datasets[3].data[i] = run['skipped'] || 0;
      }
    }
  } catch (err) {
    error.value = `Error loading test history: ${err}`;
    testRuns.value = [];
  }

  loading.value = false;
  new Modal('#showTestRuns_' + props.featureId).show();
};

const deleteTests = async (c: string, s: string, f: string) => {
  loading.value = true;

  try {
    await http.delete(`/api/v1/tests?component=${c}&suite=${s}&file-name=${f}`);

    // Refresh the page to show updated data
    loading.value = false;
    location.reload();
  } catch (err) {
    error.value = `Error deleting tests: ${err}`;
    loading.value = false;
  }
};

const closeAlert = () => {
  error.value = '';
};

const filterLastFailed = ref(false);
const switchFilterLastFailed = () => {
  filterLastFailed.value = !filterLastFailed.value;
  filterCriteria.value = -1;
};
const filterCriteria = ref(0);
const filter = (f: number) => {
  filterCriteria.value = f;
};

// Watch for changes in the selectedComponentProp
watch(
  () => props.selectedComponentProp,
  (newValue) => {
    if (newValue) {
      // Check if this component exists in our list
      const componentExists = components.value.includes(newValue);
      if (componentExists) {
        selectedComponent.value = newValue;
      }
      // If component doesn't exist in list yet, it will be handled in setComponentsList when data loads
    }
  },
  { immediate: true }
);

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
