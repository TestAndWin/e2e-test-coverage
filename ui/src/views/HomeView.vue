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
        <div class="rounded area p-2">
          <p class="">
            <span class="result total" style="width: 80px">Total</span>&nbsp; <span class="result passes" style="width: 80px">Passes</span> &nbsp;
            <span class="result failures" style="width: 80px">Failures</span> &nbsp; <span class="result pending" style="width: 80px">Pending</span> &nbsp;
            <span class="result skipped" style="width: 80px">Skipped</span> &nbsp;
            <span class="result expl-test" style="width: 80px">Expl. Test</span>
          </p>
          <p>
            <span class="result total" style="width: 80px">
              <i class="bi bi-caret-up"></i>
              <i class="bi bi-caret-down"></i>
            </span>
            indicates if there are more or less tests for the last 28 days.
          </p>
          <hr/>
          <p>The weather symbol shows the test stability for the last 28 days.</p>
          <p>
            <i class="bi bi-sun"></i>
            All test runs without failures
          </p>
          <p>
              <i class="bi bi-cloud-sun"></i>
            &lt;=15% test runs with failures
          </p>
          <p>
              <i class="bi bi-cloud"></i>
            &lt;=30% test runs with failures
          </p>
          <p>
              <i class="bi bi-cloud-rain"></i>
            &lt;= 50% test runs with failures
          </p>
          <p>
              <i class="bi bi-lightning"></i>
            > 50% test runs with failures
          </p>
          <hr/>
          <p>$ symbol shows the business value of a feature. &dollar; = low, &dollar;&dollar; = medium, &dollar;&dollar;&dollar; = high</p>
        </div>
        <p></p>
        <p class="card-text"><a href="/docs/index.html" class="link-info">API Documentation</a></p>
        <p>&copy; TestAndWin 2022</p>
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
