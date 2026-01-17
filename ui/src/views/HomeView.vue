<template>
  <div class="container">
    <div class="card">
      <div class="card-body">
        <h5 class="card-title">e2e test coverage</h5>
        <p class="card-text">
          e2e test coverage provides a summary of the various product areas and the number of automated tests in each
          area, along with their respective test status.
        </p>
        <p>Furthermore, insights from exploratory testing on the different product areas can be obtained.</p>
        <p>The results of the tests conducted over the last 28 days are considered. Following numbers are displayed:</p>
        <div class="rounded area p-2">
          <p class="">
            <span class="result total" style="width: 80px">Total</span>&nbsp;
            <span class="result passes" style="width: 80px">Passes</span> &nbsp;
            <span class="result failures" style="width: 80px">Failures</span> &nbsp;
            <span class="result pending" style="width: 80px">Pending</span> &nbsp;
            <span class="result skipped" style="width: 80px">Skipped</span> &nbsp;
            <span class="result expl-test" style="width: 80px">Expl. Test</span>
          </p>
          <p>
            <span class="result total" style="width: 80px">
              <i class="bi bi-caret-up"></i>
              <i class="bi bi-caret-down"></i>
            </span>
            indicates whether there has been an increase or decrease in the number of tests over the last 28 days.
          </p>
          <hr />
          <p>The weather icon represents the stability of testing over the last 28 days.</p>
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
          <hr />
          <p>
            $ symbol represents the business value of a feature. &dollar; = low, &dollar;&dollar; = medium,
            &dollar;&dollar;&dollar; = high
          </p>
        </div>
        <p></p>
        <p>&copy; TestAndWin 2026</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';

onMounted(() => {
  document.title = 'e2e coverage - Home';
});
</script>
