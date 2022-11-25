<template>
  <div v-if="loading" variant="info" class="spinner-border" role="status">
    <span class="visually-hidden">Loading...</span>
  </div>

  <div class="container">
    <div v-for="(test, index) in tests" :key="test['id']" class="test shadow p-2 mb-2 rounded">
      <div :id="`test-${index}`" class="row">
        <div class="col-5">
          <h6 @click="showTestRuns(test['suite'], test['file-name'])" class="pointer">
            {{ test["suite"] }}
          </h6>
        </div>
        <div class="col-5">
          <span class="result total">
            {{ test["total"] }}
            <i v-if="test['total'] > test['first-total']" class="bi bi-caret-up"></i>
            <i v-if="test['total'] < test['first-total']" class="bi bi-caret-down"></i>
          </span> &nbsp; 
          <span class="result passes">{{ test["passes"] }}</span> &nbsp;
          <span class="result failures">{{ test["failures"] }}</span> &nbsp; <span class="result pending">{{ test["pending"] }}</span> &nbsp;
          <span class="result skipped">{{ test["skipped"] }}</span>
        </div>
        <div class="col">&nbsp;</div>
      </div>
      <div class="row">
        <div class="col">
          <span class="test-suite d-flex justify-content-between">File: {{ test["file-name"] }}</span>
        </div>
      </div>
      <div class="row">
        <div class="col">
          <span class="test-suite d-flex justify-content-between">Test run: {{ test["test-run"] }}</span>
          <span v-if="test['area-id'] == 0"><i>Not assigned to an area/feature</i></span>
        </div>
      </div>
    </div>

    <!-- Modal to all tests -->
    <div class="modal modal-fullscreen-sm-down fade" :id="'showTestRuns_' + featureId" tabindex="-1" aria-labelledby="showTestRunsLabel" aria-hidden="true">
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
            <div class="col-3">{{ tr["test-run"] }}</div>
            <div class="col">{{ tr["total"] }}</div>
            <div class="col">{{ tr["passes"] }}</div>
            <div class="col">{{ tr["failures"] }}</div>
            <div class="col">{{ tr["pending"] }}</div>
            <div class="col">{{ tr["skipped"] }}</div>
          </div>
          <Line ref="chart" :chart-data="chartData" :chart-options="chartOptions" />
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { fetchData } from "./ApiHelper";
import { Modal } from "bootstrap";

import { Line } from "vue-chartjs";
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, LineController, LineElement, PointElement } from "chart.js";
ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, LineController, LineElement, PointElement);

export default defineComponent({
  name: "TestCoverage",
  emits: ["showAlert"],
  props: {
    featureId: {
      type: Number,
      required: false,
    },
    productId: {
      type: Number,
      required: false,
    },
  },
  data() {
    return {
      loading: true,
      tests: [],
      suite: "",
      file: "",
      testRuns: [],
      chartData: {
        labels: [],
        datasets: [
          {
            label: "Total",
            backgroundColor: "blue",
            borderColor: "blue",
            data: [],
          },
          {
            label: "Pass",
            backgroundColor: "green",
            borderColor: "green",
            data: [],
          },
          {
            label: "Fail",
            backgroundColor: "red",
            borderColor: "red",
            data: [],
          },
          {
            label: "Pending",
            backgroundColor: "orange",
            borderColor: "orange",
            data: [],
          },
          {
            label: "Skipped",
            backgroundColor: "yellow",
            borderColor: "yellow",
            data: [],
          },
        ],
      },
      chartOptions: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: {
            type: "linear",
            min: 0,
          },
        },
      },
    };
  },
  methods: {
    async getTestsForProduct() {
      this.loading = true;
      await fetchData(`${process.env.VUE_APP_API_URL}/coverage/products/${this.productId}/tests`)
        .then((data) => {
          this.tests = data;
        })
        .catch((err) => {
          this.$emit("showAlert", err);
        });
      this.loading = false;
    },
    async getTestsForFeature() {
      this.loading = true;
      await fetchData(`${process.env.VUE_APP_API_URL}/coverage/features/${this.featureId}/tests`)
        .then((data) => {
          this.tests = data;
        })
        .catch((err) => {
          this.$emit("showAlert", err);
        });
      this.loading = false;
    },
    async showTestRuns(suite: string, file: string) {
      this.loading = true;
      this.suite = suite;
      this.file = file;
      await fetchData(`${process.env.VUE_APP_API_URL}/tests?suite=${suite}&file-name=${file}`)
        .then((data) => {
          this.testRuns = data;
          this.chartData.labels = [];
          for (let i = 0; i < data.length; i++) {
            const r = data.length - 1 - i;
            this.chartData.labels[i] = data[r]["test-run"];
            this.chartData.datasets[0].data[i] = data[r]["total"];
            this.chartData.datasets[1].data[i] = data[r]["passes"];
            this.chartData.datasets[2].data[i] = data[r]["failures"];
            this.chartData.datasets[3].data[i] = data[r]["pending"];
            this.chartData.datasets[4].data[i] = data[r]["skipped"];
          }
        })
        .catch((err) => {
          this.$emit("showAlert", err);
        });
      this.loading = false;
      new Modal("#showTestRuns_" + this.featureId).show();
    },
  },
  mounted() {
    if (this.featureId) {
      this.getTestsForFeature();
    }
    if (this.productId) {
      this.getTestsForProduct();
    }
  },
  // eslint-disable-next-line
  components: { Line },
});
</script>

<style scoped>
@import "../assets/styles.css";
</style>
