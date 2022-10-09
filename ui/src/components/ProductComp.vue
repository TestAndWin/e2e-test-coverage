<template>
  <div class="product container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
      <button type="button" class="btn-close" aria-label="Close" @click="closeAlert()"></button>
    </div>

    <div v-if="loading" variant="info" class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    <div class="area shadow p-2 mb-4 rounded">
      <div v-for="area in areas" :key="area['id']">
        <div :id="`area-${area['id']}`" class="">
          <div class="row">
            <div class="col">
              <h5 class="area-name justify-content-between" @click="showFeatures(area['id'])" >
                {{ area["name"] }}
                &nbsp;
                <i class="bi bi-pencil" @click="showChangeAreaModal(area['id'], area['name'])"></i>
                &nbsp;
                <i class="bi bi-trash" @click="removeArea(area['id'])"></i>
              </h5>
            </div>
          </div>
          <div v-if="areaToggle[area['id']]" >
            <div v-for="feature in features[area['id']]" :key="feature['id']" class="feature shadow p-2 mb-2 rounded">
              <div :id="`feature-${feature['id']}`" class="">
                <div class="row">
                  <div class="col">
                    <h6 class="feature-name justify-content-between">
                      {{ feature["name"] }}
                      &nbsp;
                      <i class="bi bi-pencil" @click="showChangeFeatureModal(feature['id'], feature['name'])"></i>
                      &nbsp;
                      <i class="bi bi-trash" @click="removeFeature(area['id'], feature['id'])"></i>
                    </h6>
                  </div>
                </div>
              </div>
            </div>

            <div class="input-group mb-3">
              <input type="text" class="form-control" placeholder="Feature Name" aria-label="Feature Name" aria-describedby="button-add-area" v-model="newFeature[area['id']]" />
              <button class="btn btn-outline-secondary bi bi-plus-lg" type="button" id="button-add-feature" @click="addFeature(area['id'])"></button>
            </div>
          </div>
        </div>
        <hr />
      </div>

      <div class="input-group mb-3">
        <input type="text" class="form-control" placeholder="Area Name" aria-label="Area Name" aria-describedby="button-add-area" v-model="newArea" />
        <button class="btn btn-outline-secondary bi bi-plus-lg" type="button" id="button-add-area" @click="addArea()"></button>
      </div>

      <div v-if="areas.length == 0" class="input-group mb-3">
        <input type="text" class="form-control" placeholder="Enter Product Name" aria-label="Enter Product Name" aria-describedby="button-add-product" v-model="newProduct" />
        <button class="btn btn-outline-secondary bi bi-plus-lg" type="button" id="button-add-product" @click="addProduct()"></button>
      </div>
    </div>
  </div>

  <!-- Modal to change the name of an area -->
  <div class="modal fade" ref="changeAreaName" id="changeAreaName" tabindex="-1" aria-labelledby="changeAreaNameLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="changeAreaNameLabel">Update</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form>
            <div class="mb-3">
              <input type="text" class="form-control" id="newName" v-model="newName" />
            </div>
          </form>
        </div>
        <div class="modal-body">
          <p>Please be aware that the name is used to "match" the test results.</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          <button type="button" class="btn btn-primary" data-bs-dismiss="modal" @click="changeAreaName()">Save changes</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal to change the feature name -->
  <div class="modal fade" ref="changeFeatureName" id="changeFeatureName" tabindex="-1" aria-labelledby="changeFeatureNameLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="changeFeatureNameLabel">Update</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form>
            <div class="mb-3">
              <input type="text" class="form-control" id="newName" v-model="newName" />
            </div>
          </form>
        </div>
        <div class="modal-body">
          <p>Please be aware that the name is used to "match" the test results.</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          <button type="button" class="btn btn-primary" data-bs-dismiss="modal" @click="changeFeatureName()">Save changes</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { Modal } from "bootstrap";
import { fetchData } from "./ApiHelper";

export default defineComponent({
  name: "ProductComp",
  el: "#app",
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
      newArea: "",
      newProduct: "",
      newFeature: [""],
      features: [[]],
      newName: "",
      areaIdToChange: 0,
      featureIdToChange: 0,
      error: "",
    };
  },
  methods: {
    async getAreas() {
      this.loading = true;
      await fetchData(`${process.env.VUE_APP_API_URL}/products/${this.productId}/areas`)
        .then((data) => {
          this.areas = data;
        })
        .catch((err) => {
          this.error = err;
        });
      this.areas.forEach((a) => {
        this.getFeatures(a["id"]);
      });
      this.areaToggle = new Array(this.areas.length).fill(false);
      this.loading = false;
    },
    async getFeatures(areaId: number) {
      await fetchData(`${process.env.VUE_APP_API_URL}/areas/${areaId}/features`)
        .then((data) => {
          this.features[areaId] = data;
        })
        .catch((err) => {
          this.error = err;
        });
    },
    async addProduct() {
      await fetchData(`${process.env.VUE_APP_API_URL}/products`, {
        method: "POST",
        mode: "cors",
        body: JSON.stringify({ name: this.newProduct }),
      }).catch((err) => {
        this.error = err;
      });
      this.newProduct = "";
      this.getAreas();
    },
    async addArea() {
      await fetchData(`${process.env.VUE_APP_API_URL}/areas`, {
        method: "POST",
        mode: "cors",
        body: JSON.stringify({ "product-id": this.productId, name: this.newArea }),
      }).catch((err) => {
        this.error = err;
      });
      this.newArea = "";
      this.getAreas();
    },
    async addFeature(areaId: number) {
      await fetchData(`${process.env.VUE_APP_API_URL}/features`, {
        method: "POST",
        mode: "cors",
        body: JSON.stringify({ "area-id": areaId, name: this.newFeature[areaId], description: "", importance: "" }),
      }).catch((err) => {
        this.error = err;
      });
      this.newFeature[areaId] = "";
      this.getFeatures(areaId);
    },
    async removeArea(areaId: number) {
      await fetchData(`${process.env.VUE_APP_API_URL}/areas/${areaId}`, {
        method: "DELETE",
        mode: "cors",
      }).catch((err) => {
        this.error = err;
      });
      this.getAreas();
    },
    async removeFeature(areaId: number, featureId: number) {
      await fetchData(`${process.env.VUE_APP_API_URL}/features/${featureId}`, {
        method: "DELETE",
        mode: "cors",
      }).catch((err) => {
        this.error = err;
      });
      this.getFeatures(areaId);
    },
    async showChangeAreaModal(areaId: number, name: string) {
      this.newName = name;
      this.areaIdToChange = areaId;
      new Modal("#changeAreaName").show();
    },
    async showChangeFeatureModal(featureId: number, name: string) {
      this.newName = name;
      this.featureIdToChange = featureId;
      new Modal("#changeFeatureName").show();
    },
    async changeAreaName() {
      await fetchData(`${process.env.VUE_APP_API_URL}/areas/${this.areaIdToChange}`, {
        method: "PUT",
        mode: "cors",
        body: JSON.stringify({ name: this.newName }),
      }).catch((err) => {
        this.error = err;
      });
      this.newName = "";
      this.areaIdToChange = 0;
      this.getAreas();
    },
    async changeFeatureName() {
      await fetchData(`${process.env.VUE_APP_API_URL}/features/${this.featureIdToChange}`, {
        method: "PUT",
        mode: "cors",
        body: JSON.stringify({ name: this.newName, description: "", importance: "" }),
      }).catch((err) => {
        this.error = err;
      });
      this.newName = "";
      this.featureIdToChange = 0;
      this.getAreas();
    },
    async closeAlert() {
      this.error = "";
    },
    async showFeatures(areaId: number) {
      this.areaToggle[areaId] = !this.areaToggle[areaId];
    },
  },
  mounted() {
    this.getAreas();
  },
  components: {},
});
</script>

<style scoped>
@import "../assets/styles.css";
</style>
