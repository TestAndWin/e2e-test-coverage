<template>
  <nav class="navbar navbar-expand-lg">
    <div class="container-fluid">
      <a class="navbar-brand" href="/"> <img src="./assets/logo.png" alt="Test and Win - e2e test coverage" width="110" height="50" /></a>
      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarSupportedContent"
        aria-controls="navbarSupportedContent"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarSupportedContent">
        <ul class="navbar-nav me-auto mb-2 mb-lg-0">
          <li class="nav-item">
            <a class="nav-link" href="/product">Product</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="/coverage">Coverage</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="/tests">Tests</a>
          </li>
        </ul>
      </div>
    </div>
  </nav>
  <router-view />
</template>

<style>
@import "./assets/styles.css";
</style>
