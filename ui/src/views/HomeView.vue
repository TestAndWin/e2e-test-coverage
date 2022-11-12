<template>
  <div class="container">
    <div class="card">
      <div class="card-body">
        <h5 class="card-title">e2e test coverage</h5>
        <p class="card-text">
          e2e test coverage gives an overview of the different product areas and how many automated tests exist in each area with their respective test status.
        </p>
        <p>Additionally, the impressions of exploratory tests on the different product areas can be gathered.</p>
        <p>The test results of the last 28 days are taken into account. Following numbers are displayed:</p>
        <p class="rounded area p-2">
          <span class="result total" style="width: 80px;">Total</span> &nbsp; <span class="result passes" style="width: 80px;">Passes</span> &nbsp; <span class="result failures" style="width: 80px;">Failures</span> &nbsp;
          <span class="result pending" style="width: 80px;">Pending</span> &nbsp;
          <span class="result skipped" style="width: 80px;">Skipped</span> &nbsp;
          <span class="result expl-test" style="width: 80px;">Expl. Test</span>
        </p>
        <p class="card-text"><a href="/docs/index.html" class="link-info">API Documentation</a></p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  components: {},
  mounted() {
    document.title = "e2e coverage - Home";
  },
};
</script>
