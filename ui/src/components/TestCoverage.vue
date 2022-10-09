<template>
  <div class="">
    <div v-if="loading" variant="info" class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <div v-for="(test, index) in tests" :key="test['id']" class="test shadow p-2 mb-2 rounded">
      <div :id="`test-${index}`">
        <div class="row">
          <div class="col-3">
            <h6 class="test-suite d-flex justify-content-between" @click="showTestRuns(test['suite'], test['file-name'])">
              {{ test["suite"] }}
            </h6>
          </div>
          <div class="col-1">
            <span class="result total">{{ test["total"] }}</span>&nbsp;
          </div>
          <div class="col-1">
            <span class="result passes">{{ test["passes"] }}</span>&nbsp;
          </div>
          <div class="col-1">
            <span class="result failures">{{ test["failures"] }}</span>&nbsp;
          </div>
          <div class="col-1">
            <span class="result pending">{{ test["pending"] }}</span>&nbsp;
          </div>
          <div class="col-1">
            <span class="result skipped">{{ test["skipped"] }}</span>
          </div>
        </div>
        <div class="row">
          <div class="col">
            <span class="test-suite d-flex justify-content-between">File: {{ test["file-name"] }}</span>
          </div>
        </div>
        <div class="row">
          <div class="col">
            <span class="test-suite d-flex justify-content-between">Test run: {{ test["test-run"] }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal to all tests -->
  <div class="modal modal-fullscreen-sm-down fade" :id="'showTestRuns_'+featureId" tabindex="-1" aria-labelledby="showTestRunsLabel" aria-hidden="true">
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
          <div class="col-3">{{tr['test-run']}}</div>
          <div class="col">{{tr['total']}}</div>
          <div class="col">{{tr['passes']}}</div>
          <div class="col">{{tr['failures']}}</div>
          <div class="col">{{tr['pending']}}</div>
          <div class="col">{{tr['skipped']}}</div>
        </div>
        <!--
        <Line ref="chart" :chart-data="chartData" :chart-options="chartOptions" />
        -->
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { fetchData } from "./ApiHelper";
import { Modal } from "bootstrap";

//import { Line } from "vue-chartjs";
//import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, LineController, LineElement, PointElement } from "chart.js";
//ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, LineController, LineElement, PointElement);

export default defineComponent({
  name: "TestCoverage",
  emits: ["showAlert"],
  props: {
    featureId: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      loading: true,
      tests: [],
      suite: "",
      file: "",
      testRuns: []
    };
  },
  methods: {
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
        })
        .catch((err) => {
          this.$emit("showAlert", err);
        });
      this.loading = false;
      new Modal("#showTestRuns_" + this.featureId).show();
    },
  },
  mounted() {
    this.getTestsForFeature();
  },
  components: { /*Line*/ },
});
</script>

<style scoped>
@import "../assets/styles.css";
</style>

<!--

   chartData: {
        labels: [],
        datasets: [
          {
            label: "Total",
            backgroundColor: "blue",
            data: [],
          },
          {
            label: "Pass",
            backgroundColor: "green",
            data: [],
          },
          {
            label: "Fail",
            backgroundColor: "red",
            data: [],
          },
          {
            label: "Pending",
            backgroundColor: "orange",
            data: [],
          },
          {
            label: "Skipped",
            backgroundColor: "yellow",
            data: [],
          },
        ],
      },
      chartOptions: {
        responsive: true,
        maintainAspectRatio: false,
      },
-->