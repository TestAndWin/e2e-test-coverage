<template>
  <div class="container mb-3">
    <div class="card">
      <div class="card-body">
        <p class="card-text">Number of test cases and their status of the last 28 days including exploratory testing impression.</p>
      </div>
    </div>
  </div>
  <div class="">
    <AreaCoverage :productId="1" />
  </div>
</template>

<script lang="ts">
import AreaCoverage from "@/components/AreaCoverage.vue";

export default {
  components: {
    AreaCoverage,
  },
  mounted() {
    document.title = "e2e coverage - Coverage";
  },
};
</script>
