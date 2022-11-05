<template>
  <div class="container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
      <button type="button" class="btn-close" aria-label="Close" @click="closeAlert()"></button>
    </div>

    <div v-if="loading" variant="info" class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <div v-for="area in areas" :key="area['id']" :id="`area-${area['id']}`" class="area shadow p-2 mb-2 rounded">
      <div class="row">
        <div @click="showFeatures(area['id'])" class="col-4">
          <h4>{{ area["name"] }}</h4>
        </div>
        <div class="col-1">
          <span v-if="area['total'] < 1" class="result failures">
            {{ area["total"] }}<br />
            <div class="d-none d-lg-block">total</div>
          </span>
          <span v-if="area['total'] > 0" class="result total">
            {{ area["total"] }}<br />
            <div class="d-none d-lg-block">total</div>
          </span>
          &nbsp;
        </div>
        <div class="col-1">
          <span class="result passes">
            {{ area["passes"] }}<br />
            <div class="d-none d-lg-block">passes</div>
          </span>
          &nbsp;
        </div>
        <div class="col-1">
          <span class="result failures">
            {{ area["failures"] }}<br />
            <div class="d-none d-lg-block">failures</div>
          </span>
          &nbsp;
        </div>
        <div class="col-1">
          <span class="result pending">
            {{ area["pending"] }}<br />
            <div class="d-none d-lg-block">pending</div>
          </span>
          &nbsp;
        </div>
        <div class="col-1">
          <span class="result skipped">
            {{ area["skipped"] }}<br />
            <div class="d-none d-lg-block">skipped</div>
          </span>
        </div>
        <div class="col">
          <span class="result expl-test" @click="showExplTests(area['id'])">
            {{ parseFloat(area["expl-rating"]).toFixed(1) }} ({{ area["expl-tests"] }})<br />
            <div class="d-none d-lg-block">Expl. Test</div>
          </span>
          &nbsp;
          <span class="result expl-test" @click="logExplTest(area['id'])">
            Log<br />
            <div class="d-none d-lg-block">new</div>
          </span>
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
            <label>Rating ({{ etRating }} stars)</label>
            <input type="range" class="form-range" min="0" max="5" step="1.0" id="etRating" v-model="etRating" />
            <label>Test Date</label>
            <input id="etDate" class="form-control" type="date" v-model="etDate" />
          </div>
        </form>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          <button type="button" class="btn btn-primary" data-bs-dismiss="modal" @click="saveExplTest">Save test</button>
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
                <strong>{{ new String(et["test-run"]).split("T")[0] }} / {{ et["rating"] }} stars</strong>
              </label>
              <p>{{ et["summary"] }}</p>
            </div>
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
      etRating: "3",
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
      this.etRating = "3";
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
  components: { FeatureCoverage },
});
</script>

<style scoped>
@import "../assets/styles.css";
</style>
