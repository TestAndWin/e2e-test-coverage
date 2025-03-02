<template>
  <div class="container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
      <button type="button" class="btn-close pointer" aria-label="Close" @click="closeAlert()"></button>
    </div>

    <div v-if="loading" class="info spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <h4 class="">Coverage</h4>
    <div v-for="area in areas" :key="area['id']" :id="`area-${area['id']}`" class="area shadow p-2 mb-2 rounded">
      <div class="row">
        <div @click="showFeatures(area.id ?? 0)" class="col-5 pointer">
          <h4>{{ area['name'] }}</h4>
        </div>
        <div class="col-5 mb-2">
          <TestResult :test="area" />
        </div>
        <div class="col mb-2">
          <span class="result expl-test pointer" @click="showExplTests(area.id ?? 0)">
            {{ parseFloat(area['expl-rating']?.toString() || '0').toFixed(1) }} ({{ area['expl-tests'] || 0 }})
          </span>
          &nbsp;
          <span class="result expl-test pointer" @click="showLogExplTest(area.id ?? 0)"> New </span>
        </div>
      </div>
      <FeatureCoverage @show-alert="showAlert" v-if="areaToggle[area.id ?? 0]" :areaId="area.id ?? 0" />
    </div>
  </div>

  <!-- Modal to add an exploratory test-->
  <div class="modal fade" id="logExplTest" tabindex="-1" aria-labelledby="logExplTestLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="logExplTestLabel">Log Exploratory Test</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <form>
          <div class="modal-body">
            <label>Summary</label>
            <textarea type="text" class="form-control" id="etSummary" v-model="etSummary" />
            <br />
            <label>Rating</label>
            <br />
            <select v-model="etRating">
              <option disabled value="">Number stars</option>
              <option>0</option>
              <option>1</option>
              <option>2</option>
              <option>3</option>
              <option>4</option>
              <option>5</option>
            </select>
            <br /><br />
            <label>Test Date</label>
            <input id="etDate" class="form-control" type="date" v-model="etDate" />
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Close</button>
            <button type="button" class="btn btn-primary pointer" data-bs-dismiss="modal" @click="saveExplTest">
              Save test
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>

  <!-- Modal to show exploratory test-->
  <div class="modal fade" id="showExplTest" tabindex="-1" aria-labelledby="showExplTestLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="showExplTestLabel">Exploratory Tests</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <form>
          <div class="modal-body">
            <div v-for="et in explTests" :key="et['id']">
              <label>
                <strong>{{ new String(et['test-run']).split('T')[0] }} / </strong>
                <i
                  class="bi"
                  v-for="n in 5"
                  :class="{
                    'bi-star-fill': n <= et['rating'],
                    'bi-star': n > et['rating']
                  }"
                  :key="n"
                ></i>
              </label>
              <p>{{ et['summary'] }}</p>
              <p>Tester: {{ et['tester'] }}</p>
            </div>
            <div v-if="explTests.length == 0">No exploratory tests logged.</div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Close</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Modal } from 'bootstrap';
import FeatureCoverage from './FeatureCoverage.vue';
import TestResult from './TestResult.vue';
import http from '@/common-http';
import type { Area } from '@/types';

const props = defineProps({
  productId: Number
});

const loading = ref(true);
const error = ref('');

const areaToggle = ref([false]);
const areas = ref<Area[]>([]);
const getAreas = async () => {
  loading.value = true;
  try {
    const response = await http.get(`/api/v1/coverage/${props.productId}/areas`);

    // Extract data from StandardResponse format
    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      areas.value = response.data.data;
    } else {
      areas.value = [];
    }

    // Ensure all required properties have default values
    areas.value = areas.value.map((area: Area) => ({
      ...area,
      total: area.total || 0,
      passes: area.passes || 0,
      failures: area.failures || 0,
      pending: area.pending || 0,
      skipped: area.skipped || 0,
      'first-total': area['first-total'] || 0,
      'expl-tests': area['expl-tests'] || 0,
      'expl-rating': area['expl-rating'] || 0
    }));

    areaToggle.value = new Array(areas.value.length).fill(false);
  } catch (err) {
    error.value = `Error loading coverage data: ${err}`;
  }

  loading.value = false;
};

const showFeatures = (areaId: number) => {
  areaToggle.value[areaId] = !areaToggle.value[areaId];
};

const etAreaId = ref(0);
const etSummary = ref('');
const etRating = ref(0);
const etDate = ref('');
const explTests = ref([]);
const saveExplTest = async () => {
  await http
    .post(`/api/v1/expl-tests`, {
      'area-id': etAreaId.value,
      summary: etSummary.value,
      rating: Number(etRating.value),
      'test-run': etDate.value + 'T00:00:00.000Z'
    })
    .catch((err) => {
      error.value = err + ' | ' + err.response.data.error;
    })
    .finally(() => {
      etDate.value = new Date().toISOString().split('T')[0];
      etSummary.value = '';
      etRating.value = 0;
    });
  getAreas();
};

const showLogExplTest = (areaId: number) => {
  etAreaId.value = areaId;
  etDate.value = new Date().toISOString().split('T')[0];
  new Modal('#logExplTest').show();
};

const showExplTests = async (areaId: number) => {
  try {
    const response = await http.get(`/api/v1/expl-tests/area/${areaId}`);
    
    // Handle standardized response format
    if (response.data && response.data.data) {
      explTests.value = response.data.data;
    } else if (Array.isArray(response.data)) {
      explTests.value = response.data;
    } else {
      explTests.value = [];
    }
    
    new Modal('#showExplTest').show();
  } catch (err: any) {
    if (err.response && err.response.status === 404) {
      // For 404, just show empty tests rather than error
      explTests.value = [];
      new Modal('#showExplTest').show();
    } else {
      const errorMsg = err.response?.data?.error || err.message || String(err);
      error.value = `Error: ${errorMsg}`;
    }
  }
};

const closeAlert = () => {
  error.value = '';
};

const showAlert = (msg: never) => {
  error.value = msg;
};

onMounted(() => {
  getAreas();
});
</script>

<style scoped>
@import '../assets/styles.css';
</style>
