<template>
  <div class="container">
    <div v-if="loading" variant="info" class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <div v-for="feature in features" :key="feature['id']" class="feature shadow p-2 mb-2 rounded">
      <div :id="`feature-${feature['id']}`" class="row">
        <div class="col-5">
          <h5 @click="showTests(feature['id'])">{{ feature["name"] }}</h5>
        </div>
        <div class="col-5">
          <span v-if="feature['total'] < 1" class="result failures">{{ feature["total"] }}</span>
          <span v-if="feature['total'] > 0" class="result total">{{ feature["total"] }}</span>
          &nbsp;
          <span class="result passes">{{ feature["passes"] }}</span> &nbsp; <span class="result failures">{{ feature["failures"] }}</span> &nbsp;
          <span class="result pending">{{ feature["pending"] }}</span> &nbsp;
          <span class="result skipped">{{ feature["skipped"] }}</span>
        </div>
        <div class="col">&nbsp;</div>
      </div>
      <TestCoverage @show-alert="showAlert" v-if="featureToggle[feature['id']]" :featureId="feature['id']" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import TestCoverage from "@/components/TestCoverage.vue";
import { fetchData } from "./ApiHelper";

export default defineComponent({
  name: "FeatureCoverage",
  emits: ["showAlert"],
  props: {
    areaId: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      loading: true,
      features: [],
      featureToggle: [false],
    };
  },
  methods: {
    async getFeatures() {
      this.loading = true;
      await fetchData(`${process.env.VUE_APP_API_URL}/coverage/areas/${this.areaId}/features`)
        .then((data) => {
          this.features = data;
        })
        .catch((err) => {
          this.$emit("showAlert", err);
        });
      this.featureToggle = new Array(this.features.length).fill(false);
      this.loading = false;
    },
    showTests(featureId: number) {
      this.featureToggle[featureId] = !this.featureToggle[featureId];
    },
    showAlert(msg: never) {
      this.$emit("showAlert", msg);
    },
  },
  mounted() {
    this.getFeatures();
  },
  components: { TestCoverage },
});
</script>

<style scoped>
@import "../assets/styles.css";
</style>
