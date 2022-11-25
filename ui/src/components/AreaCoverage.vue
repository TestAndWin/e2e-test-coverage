<template>
  <div class="container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
      <button type="button" class="btn-close pointer" aria-label="Close" @click="closeAlert()"></button>
    </div>

    <div v-if="loading" variant="info" class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <div v-for="area in areas" :key="area['id']" :id="`area-${area['id']}`" class="area shadow p-2 mb-2 rounded">
      <div class="row">
        <div @click="showFeatures(area['id'])" class="col-5 pointer">
          <h4>{{ area["name"] }}</h4>
        </div>
        <div class="col-5 mb-2">
          <span v-if="area['total'] < 1" class="result failures">
            {{ area["total"] }}
          </span>
          <span v-if="area['total'] > 0" class="result total">
            {{ area["total"] }}
            <i v-if="area['total'] > area['first-total']" class="bi bi-caret-up"></i>
            <i v-if="area['total'] < area['first-total']" class="bi bi-caret-down"></i>
          </span>
          &nbsp;
          <span class="result passes">
            {{ area["passes"] }}
          </span>
          &nbsp;
          <span class="result failures">
            {{ area["failures"] }}
          </span>
          &nbsp;
          <span class="result pending">
            {{ area["pending"] }}
          </span>
          &nbsp;
          <span class="result skipped">
            {{ area["skipped"] }}
          </span>
        </div>
        <div class="col mb-2">
          <span class="result expl-test pointer" @click="showExplTests(area['id'])"> {{ parseFloat(area["expl-rating"]).toFixed(1) }} ({{ area["expl-tests"] }}) </span>
          &nbsp;
          <span class="result expl-test pointer" @click="logExplTest(area['id'])"> New </span>
        </div>
      </div>
      <FeatureCoverage @show-alert="showAlert" v-if="areaToggle[area['id']]" :areaId="area['id']" />
    </div>
  </div>

  <!-- Modal to add an exploratory test-->
  <div class="modal fade" ref="logExplTest" id="logExplTest" tabindex="-1" aria-labelledby="logExplTestLabel" aria-hidden="true">
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
            <star-rating :show-rating="false" active-color="#2c3e50" v-model:rating="etRating" /><br />
            <label>Test Date</label>
            <input id="etDate" class="form-control" type="date" v-model="etDate" />
          </div>
        </form>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Close</button>
          <button type="button" class="btn btn-primary pointer" data-bs-dismiss="modal" @click="saveExplTest">Save test</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal to show exploratory test-->
  <div class="modal fade" ref="showExplTest" id="showExplTest" tabindex="-1" aria-labelledby="showExplTestLabel" aria-hidden="true">
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
                <strong>{{ new String(et["test-run"]).split("T")[0] }} / </strong>
                <i class="bi" v-for="n in 5" :class="{ 'bi-star-fill': n <= et['rating'], 'bi-star': n > et['rating'] }" :key="n"></i>
              </label>
              <p>{{ et["summary"] }}</p>
            </div>
            <div v-if="explTests.length == 0">No exploratory tests logged.</div>
          </div>
        </form>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { Modal } from "bootstrap";
import FeatureCoverage from "./FeatureCoverage.vue";
import { fetchData } from "./ApiHelper";
import StarRating from "vue-star-rating";

export default defineComponent({
  name: "AreaCoverage",
  props: {
    productId: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      loading: true,
      areas: [],
      areaToggle: [false],
      etAreaId: 0,
      etSummary: "",
      etRating: 0,
      etDate: "",
      explTests: [],
      error: "",
    };
  },
  methods: {
    async getAreas() {
      this.loading = true;
      await fetchData(`${process.env.VUE_APP_API_URL}/coverage/${this.productId}/areas`)
        .then((data) => {
          this.areas = data;
        })
        .catch((err) => {
          this.error = err;
        });
      this.areaToggle = new Array(this.areas.length).fill(false);
      this.loading = false;
    },
    async showFeatures(areaId: number) {
      this.areaToggle[areaId] = !this.areaToggle[areaId];
    },
    async logExplTest(areaId: number) {
      this.etAreaId = areaId;
      new Modal("#logExplTest").show();
    },
    async saveExplTest() {
      await fetchData(`${process.env.VUE_APP_API_URL}/expl-tests`, {
        method: "POST",
        mode: "cors",
        body: JSON.stringify({ "area-id": this.etAreaId, summary: this.etSummary, rating: this.etRating, "test-run": this.etDate + "T00:00:00.000Z" }),
      }).catch((err) => {
        this.error = err;
      });

      this.etDate = new Date().toISOString().split("T")[0];
      this.etSummary = "";
      this.etRating = 0;
      this.getAreas();
    },
    async showExplTests(areaId: number) {
      await fetchData(`${process.env.VUE_APP_API_URL}/expl-tests/area/${areaId}`)
        .then((data) => {
          this.explTests = data;
        })
        .catch((err) => {
          this.error = err;
        });
      new Modal("#showExplTest").show();
    },
    async closeAlert() {
      this.error = "";
    },
    async showAlert(msg: never) {
      this.error = msg;
    },
  },
  mounted() {
    this.getAreas();
    this.etDate = new Date().toISOString().split("T")[0];
  },
  components: { FeatureCoverage, StarRating },
});
</script>

<style scoped>
@import "../assets/styles.css";
</style>
