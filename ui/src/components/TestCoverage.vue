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
          ((!filterLastFailed && test['percent'] >= filterCriteria) || (filterLastFailed && test['failures'] > 0)) &&
          (selectedComponent == 'all' || selectedComponent == test['component'])
        "
      >
        <div class="test shadow p-2 mb-2 rounded">
          <div class="row">
            <div class="col-5">
              <h6 @click="showTestRuns(test['component'], test['suite'], test['file-name'])" class="pointer">
                {{ test['suite'] }}
              </h6>
            </div>
            <div class="col-5">
              <TestResult :test="test" />
            </div>
            <div class="col">
              <i v-if="test['percent'] == sun" class="bi bi-sun"></i>
              <i v-if="test['percent'] > sun && test['percent'] <= sunCloud" class="bi bi-cloud-sun"></i>
              <i v-if="test['percent'] > sunCloud && test['percent'] <= cloud" class="bi bi-cloud"></i>
              <i v-if="test['percent'] >= cloud && test['percent'] <= cloudRain" class="bi bi-cloud-rain"></i>
              <i v-if="test['percent'] > cloudRain" class="bi bi-lightning"></i>
              &nbsp;&nbsp;
              <i
                @click="deleteTests(test['component'], test['suite'], test['file-name'])"
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
import { ref, onMounted } from 'vue';
import { Modal } from 'bootstrap';
import TestResult from '@/components/TestResult.vue';
import http from '@/common-http';

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
  featureId: Number
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

const tests = ref([]);
const getTestsForProduct = async () => {
  loading.value = true;
  await http
    .get(`/api/v1/coverage/products/${props.productId}/tests`)
    .then((response) => {
      if (response.data) {
        tests.value = response.data;
        calculatePercentage();
        setComponentsList();
      }
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
      if (response.data) {
        tests.value = response.data;
        calculatePercentage();
        setComponentsList();
      }
    })
    .catch((err) => {
      error.value = err;
    });

  loading.value = false;
};

// Calculate the percentage of the failed tests
const calculatePercentage = () => {
  for (const test of tests.value) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (test as any)['percent'] = test['failed-test-runs'] / test['total-test-runs'];
  }
};

const components = ref<string[]>([]);
const selectedComponent = ref('all');
// Get list of components from the test results
const setComponentsList = () => {
  components.value = ['all', ...new Set(tests.value.map((item) => item['component']))];
};

const testRuns = ref([]);
const suite = ref('');
const file = ref('');
const showTestRuns = async (c: string, s: string, f: string) => {
  loading.value = true;
  suite.value = s;
  file.value = f;

  await http
    .get(`/api/v1/tests?component=${c}&suite=${s}&file-name=${f}`)
    .then((response) => {
      // We use only the lastest 10
      testRuns.value = response.data.slice(0, 10);

      chartData.value.labels = [];
      for (let i = 0; i < testRuns.value.length; i++) {
        // different order
        const r = testRuns.value.length - 1 - i;
        chartData.value.labels[i] = testRuns.value[r]['test-run'];
        //chartData.value.datasets[0].data[i] = testRuns.value[r]['total'];
        chartData.value.datasets[0].data[i] = testRuns.value[r]['passes'];
        chartData.value.datasets[1].data[i] = testRuns.value[r]['failures'];
        chartData.value.datasets[2].data[i] = testRuns.value[r]['pending'];
        chartData.value.datasets[3].data[i] = testRuns.value[r]['skipped'];
      }
    })
    .catch((err) => {
      error.value = err;
    });
  loading.value = false;
  new Modal('#showTestRuns_' + props.featureId).show();
};

const deleteTests = async (c: string, s: string, f: string) => {
  loading.value = true;

  await http.delete(`/api/v1/tests?component=${c}&suite=${s}&file-name=${f}`).catch((err) => {
    error.value = err;
  });
  loading.value = false;
  location.reload();
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
