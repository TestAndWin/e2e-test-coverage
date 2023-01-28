<template>
  <div class="product container">
    <div v-if="error" class="alert alert-danger">
      <span>{{ error }}</span>
      <button type="button" class="btn-close" aria-label="Close" @click="closeAlert()"></button>
    </div>

    <div v-if="loading" class="info spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>

    <h4 class="">Areas</h4>
    <div class="area shadow p-2 mb-4 rounded">
      <div v-for="area in areas" :key="area['id']">
        <div :id="`area-${area['id']}`" class="">
          <div class="row">
            <div class="col">
              <h5 class="area-name justify-content-between pointer" @click="showFeatures(area['id'])">
                {{ area['name'] }}
                &nbsp;
                <a @click="showChangeAreaModal(area['id'], area['name'])"><i class="bi bi-pencil pointer"></i></a>&nbsp;
                <a @click="removeArea(area['id'])"><i class="bi bi-trash pointer"></i></a>
              </h5>
            </div>
          </div>
          <div v-if="areaToggle[area['id']]">
            <div v-for="feature in features[area['id']]" :key="feature['id']" class="feature shadow p-2 mb-2 rounded">
              <div :id="`feature-${feature['id']}`" class="">
                <div class="row">
                  <div class="col">
                    <h6 class="feature-name justify-content-between">
                      {{ feature['name'] }}
                      [{{ feature['business-value'] }}]
                      <a v-if="feature['documentation']" v-bind:href="feature['documentation']" target="_blank">
                        <i class="bi bi-file-text pointer" style="color: #2c3e50"></i>
                      </a>
                      &nbsp;
                      <a v-if="feature['url']" v-bind:href="feature['url']" target="_blank">
                        <i class="bi bi-box-arrow-up-right pointer" style="color: #2c3e50"></i>
                      </a>
                      &nbsp;
                      <a @click="showUpdateFeatureModal(feature['id'], feature['name'], feature['documentation'], feature['url'], feature['business-value'])">
                        <i class="bi bi-pencil pointer"></i>
                      </a>
                      &nbsp;
                      <a @click="removeFeature(area['id'], feature['id'])"> <i class="bi bi-trash pointer"></i></a>
                    </h6>
                  </div>
                </div>
              </div>
            </div>

            <div class="input-group mb-3">
              <input type="text" class="form-control" placeholder="Feature Name" aria-label="Feature Name" aria-describedby="button-add-area" v-model="newFeature[area['id']]" />
              <button class="btn btn-outline-secondary bi bi-plus-lg pointer" type="button" id="button-add-feature" @click="addFeature(area['id'])" />
            </div>
          </div>
        </div>
        <hr />
      </div>

      <div v-if="products.length > 0" class="input-group mb-3">
        <input type="text" class="form-control" placeholder="Area Name" aria-label="Area Name" aria-describedby="button-add-area" v-model="newArea" />
        <button class="btn btn-outline-secondary bi bi-plus-lg pointer" type="button" id="button-add-area" @click="addArea()"></button>
      </div>

      <div v-if="products.length == 0 && !error" class="input-group mb-3">
        <input type="text" class="form-control" placeholder="Enter Product Name" aria-label="Enter Product Name" aria-describedby="button-add-product" v-model="newProduct" />
        <button class="btn btn-outline-secondary bi bi-plus-lg pointer" type="button" id="button-add-product" @click="addProduct()" />
      </div>
    </div>
  </div>

  <!-- Modal to change the name of an area -->
  <div class="modal fade" id="changeAreaName" tabindex="-1" aria-labelledby="changeAreaNameLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="changeAreaNameLabel">Update</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="changeAreaName">
            <div class="mb-3">
              <label>Name</label><input type="text" class="form-control" id="newAreaName" v-model="newAreaName" />
              <label>Please be aware that the name is used to "match" the test results.</label>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Close</button>
              <button type="submit" class="btn btn-primary pointer" data-bs-dismiss="modal">Save changes</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal to edit a feature -->
  <div class="modal fade" id="updateFeature" tabindex="-1" aria-labelledby="updateFeatureLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="updateFeatureLabel">Update</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="changeFeature">
            <div class="mb-3">
              <label>Name</label><input type="text" class="form-control" id="newFeatureName" v-model="newFeatureName" />
              <label>Please be aware that the name is used to "match" the test results.</label><br /><br />
              <label>Business Value</label><br />
              <select v-model="featureBusinessValue">
                <option disabled value="">Please select</option>
                <option>low</option>
                <option>medium</option>
                <option>high</option>
              </select>
              <br /><br />
              <label>Link to Documentation</label>
              <input type="text" class="form-control" id="featureDocumentation" v-model="featureDocumentation" /><br />
              <label>URL</label><input type="text" class="form-control" id="featureUrl" v-model="featureUrl" /><br />
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Close</button>
              <button type="submit" class="btn btn-primary pointer" data-bs-dismiss="modal">Save changes</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, ref, onMounted } from 'vue';
import { Modal } from 'bootstrap';
import http from '@/common-http';

const props = defineProps({
  productId: Number,
});

const loading = ref(true);
const featureBusinessValue = ref('');
const featureDocumentation = ref('');
const featureUrl = ref('');
const areaIdToChange = ref(0);
const featureIdToChange = ref(0);
const error = ref('');

// Products
const products = ref([]);
const getProducts = async () => {
  loading.value = true;
  await http
    .get(`/api/v1/products`)
    .then((response) => {
      products.value = response.data;
    })
    .catch((err) => {
      error.value = err + ' | ' + err.response?.data?.error;
    });
  loading.value = false;
};

const newProduct = ref('');
const addProduct = async () => {
  await http.post(`/api/v1/products`, { name: newProduct.value }).catch((err) => {
    error.value = err + ' | ' + err.response?.data?.error;
  });
  newProduct.value = '';
  getProducts();
};

// Areas
const areas = ref([]);
const areaToggle = ref([false]);
const getAreas = async () => {
  loading.value = true;
  await http
    .get(`/api/v1/products/${props.productId}/areas`)
    .then((response) => {
      areas.value = response.data;
      areas.value.forEach((a) => {
        getFeatures(a['id']);
      });
      areaToggle.value = new Array(areas.value.length).fill(false);
    })
    .catch((err) => {
      error.value = err + ' | ' + err.response?.data?.error;
    });
  loading.value = false;
};

const newArea = ref('');
const addArea = async () => {
  await http.post(`/api/v1/areas`, { 'product-id': props.productId, name: newArea.value }).catch((err) => {
    error.value = err + ' | ' + err.response?.data?.error;
  });
  newArea.value = '';
  getAreas();
};

const removeArea = async (areaId: number) => {
  await http.delete(`/api/v1/areas/${areaId}`).catch((err) => {
    error.value = err + ' | ' + err.response?.data?.error;
  });
  getAreas();
};

const newAreaName = ref('');
const changeAreaName = async () => {
  await http
    .put(`/api/v1/areas/${areaIdToChange.value}`, { name: newAreaName.value })
    .catch((err) => {
      error.value = err + ' | ' + err.response?.data?.error;
    })
    .finally(() => {
      newAreaName.value = '';
      areaIdToChange.value = 0;
    });
  getAreas();
};

const showChangeAreaModal = (areaId: number, name: string) => {
  newAreaName.value = name;
  areaIdToChange.value = areaId;
  new Modal('#changeAreaName').show();
};

const showFeatures = (areaId: number) => {
  areaToggle.value[areaId] = !areaToggle.value[areaId];
};

// Features
const features = ref([[]]);
const getFeatures = async (areaId: number) => {
  await http
    .get(`/api/v1/areas/${areaId}/features`)
    .then((response) => {
      features.value[areaId] = response.data;
    })
    .catch((err) => {
      error.value = err + ' | ' + err.response?.data?.error;
    });
};

const newFeature = ref(['']);
const addFeature = async (areaId: number) => {
  await http.post(`/api/v1/features`, { 'area-id': areaId, name: newFeature.value[areaId], documentation: '', url: '', 'business-value': '' }).catch((err) => {
    error.value = err + ' | ' + err.response?.data?.error;
  });
  newFeature.value[areaId] = '';
  getFeatures(areaId);
};

const removeFeature = async (areaId: number, featureId: number) => {
  await http.delete(`/api/v1/features/${featureId}`).catch((err) => {
    error.value = err + ' | ' + err.response?.data?.error;
  });
  getFeatures(areaId);
};

const newFeatureName = ref('');

const showUpdateFeatureModal = (featureId: number, name: string, documentation: string, url: string, businessValue: string) => {
  newFeatureName.value = name;
  featureIdToChange.value = featureId;
  featureBusinessValue.value = businessValue;
  featureDocumentation.value = documentation;
  featureUrl.value = url;
  new Modal('#updateFeature').show();
};

const changeFeature = async () => {
  await http
    .put(`/api/v1/features/${featureIdToChange.value}`, {
      name: newFeatureName.value,
      documentation: featureDocumentation.value,
      url: featureUrl.value,
      'business-value': featureBusinessValue.value,
    })
    .catch((err) => {
      error.value = err + ' | ' + err.response?.data?.error;
    })
    .finally(() => {
      newFeatureName.value = '';
      featureBusinessValue.value = '';
      featureDocumentation.value = '';
      featureUrl.value = '';
      featureIdToChange.value = 0;
    });

  getAreas();
};

const closeAlert = () => {
  error.value = '';
};

onMounted(() => {
  getProducts();
  getAreas();
});
</script>

<style scoped>
@import '../assets/styles.css';
</style>
