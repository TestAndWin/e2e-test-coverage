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
              <h5 class="area-name justify-content-between pointer" @click="showFeatures(area['id'] ?? 0)">
                {{ area['name'] }}
                &nbsp;
                <a @click="showChangeAreaModal(area['id'] ?? 0, area['name'] ?? '')"
                  ><i class="bi bi-pencil pointer"></i></a
                >&nbsp;
                <a @click="removeArea(area['id'] ?? 0)"> <i class="bi bi-trash pointer"></i></a>
              </h5>
            </div>
          </div>
          <div v-if="areaToggle[area['id'] ?? 0]">
            <div
              v-for="feature in features[area['id'] ?? 0] || []"
              :key="feature['id']"
              class="feature shadow p-2 mb-2 rounded"
            >
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
                      <a
                        @click="
                          showUpdateFeatureModal(
                            feature['id'],
                            feature['name'],
                            feature['documentation'],
                            feature['url'],
                            feature['business-value']
                          )
                        "
                      >
                        <i class="bi bi-pencil pointer"></i>
                      </a>
                      &nbsp;
                      <a @click="removeFeature(area['id'] ?? 0, feature['id'] ?? 0)">
                        <i class="bi bi-trash pointer"></i
                      ></a>
                    </h6>
                  </div>
                </div>
              </div>
            </div>

            <div class="input-group mb-3">
              <input
                type="text"
                class="form-control"
                placeholder="Feature Name"
                aria-label="Feature Name"
                aria-describedby="button-add-area"
                v-model="newFeature[area['id'] ?? 0]"
              />
              <button
                class="btn btn-outline-secondary bi bi-plus-lg pointer"
                type="button"
                id="button-add-feature"
                @click="addFeature(area['id'] ?? 0)"
              />
            </div>
          </div>
        </div>
        <hr />
      </div>

      <div v-if="products.length > 0" class="input-group mb-3">
        <input
          type="text"
          class="form-control"
          placeholder="Area Name"
          aria-label="Area Name"
          aria-describedby="button-add-area"
          v-model="newArea"
          @keyup.enter="addArea()"
        />
        <button
          class="btn btn-outline-secondary bi bi-plus-lg pointer"
          type="button"
          id="button-add-area"
          @click="addArea()"
        ></button>
      </div>

      <div v-if="products.length == 0 && !error" class="input-group mb-3">
        <input
          type="text"
          class="form-control"
          placeholder="Enter Product Name"
          aria-label="Enter Product Name"
          aria-describedby="button-add-product"
          v-model="newProduct"
          @keyup.enter="addProduct()"
        />
        <button
          class="btn btn-outline-secondary bi bi-plus-lg pointer"
          type="button"
          id="button-add-product"
          @click="addProduct()"
        />
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
          <form>
            <div class="mb-3">
              <label>Name</label><input type="text" class="form-control" id="newAreaName" v-model="newAreaName" />
              <label>Please be aware that the name is used to "match" the test results.</label>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Close</button>
              <button type="submit" class="btn btn-primary pointer" data-bs-dismiss="modal" @click="changeAreaName">
                Save changes
              </button>
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
          <form>
            <div class="mb-3">
              <label>Name</label
              ><input type="text" class="form-control" id="newFeatureName" v-model="newFeatureName" />
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
              <button type="submit" class="btn btn-primary pointer" data-bs-dismiss="modal" @click="changeFeature">
                Save changes
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Modal } from 'bootstrap';
import http from '@/common-http';
import type { Product, Area } from '@/types';

const props = defineProps({
  productId: Number
});

const loading = ref(true);
const featureBusinessValue = ref('');
const featureDocumentation = ref('');
const featureUrl = ref('');
const areaIdToChange = ref(0);
const featureIdToChange = ref(0);
const error = ref('');

// Products
const products = ref<Product[]>([]);
const getProducts = async () => {
  loading.value = true;
  try {
    console.log('Fetching products');
    const response = await http.get(`/api/v1/products`);
    console.log('Products API response:', response.data);

    // Extract data from StandardResponse format
    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      products.value = response.data.data;
    } else {
      console.warn('Unexpected products response format:', response.data);
      products.value = [];
    }

    console.log('Processed products:', products.value);
  } catch (err) {
    console.error('Error fetching products:', err);
    error.value = `Error loading products: ${err}`;
    products.value = []; // Initialize as empty array on error
  }
  loading.value = false;
};

const newProduct = ref('');
const addProduct = async () => {
  try {
    console.log(`Adding product: ${newProduct.value}`);

    // Check if we have a valid name
    if (!newProduct.value || newProduct.value.trim() === '') {
      error.value = 'Product name cannot be empty';
      return;
    }

    // Make the API call
    const response = await http.post(`/api/v1/products`, {
      name: newProduct.value
    });

    console.log('Add product response:', response.data);

    // Clear input
    newProduct.value = '';

    // Refresh product list
    await getProducts();
  } catch (err) {
    console.error('Error adding product:', err);
    error.value = `Error adding product: ${err}`;
  }
};

// Areas
const areas = ref<Area[]>([]);
const areaToggle = ref([false]);
const getAreas = async () => {
  loading.value = true;
  await http
    .get(`/api/v1/products/${props.productId}/areas`)
    .then((response) => {
      // Check the API response structure
      console.log('Areas API response:', response.data);

      // Extract data from StandardResponse format
      if (response.data && response.data.data && Array.isArray(response.data.data)) {
        areas.value = response.data.data;
      } else {
        console.warn('Unexpected areas response format:', response.data);
        areas.value = [];
      }

      // Safely iterate over areas if it's an array
      if (Array.isArray(areas.value)) {
        areas.value.forEach((a) => {
          getFeatures(a['id'] ?? 0);
        });
        areaToggle.value = new Array(areas.value.length).fill(false);
      } else {
        console.error('areas.value is not an array:', areas.value);
        areas.value = []; // Ensure it's always an array
      }
    })
    .catch((err) => {
      console.error('Error fetching areas:', err);
      error.value = err + ' | ' + err.response?.data?.error;
      areas.value = []; // Initialize as empty array on error
    });
  loading.value = false;
};

const newArea = ref('');
const addArea = async () => {
  try {
    console.log(`Adding area to product ${props.productId}: ${newArea.value}`);

    // Check if we have a valid name
    if (!newArea.value || newArea.value.trim() === '') {
      error.value = 'Area name cannot be empty';
      return;
    }

    // Make the API call
    const response = await http.post(`/api/v1/areas`, {
      'product-id': props.productId,
      name: newArea.value
    });

    console.log('Add area response:', response.data);

    // Clear input
    newArea.value = '';

    // Refresh area list
    await getAreas();
  } catch (err) {
    console.error('Error adding area:', err);
    error.value = `Error adding area: ${err}`;
  }
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
  try {
    console.log(`Fetching features for area ID ${areaId}`);
    const response = await http.get(`/api/v1/areas/${areaId}/features`);
    console.log(`Features API response for area ${areaId}:`, response.data);

    // Extract data from StandardResponse format
    let featureData = [];

    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      featureData = response.data.data;
    } else {
      console.warn(`Unexpected features response format for area ${areaId}:`, response.data);
      featureData = []; // Empty array as fallback
    }

    // Update the features array
    features.value[areaId] = featureData;
    console.log(`Features for area ${areaId} set to:`, features.value[areaId]);
  } catch (err) {
    console.error(`Error fetching features for area ${areaId}:`, err);
    error.value = `Error loading features: ${err}`;
    features.value[areaId] = []; // Empty array on error
  }
};

const newFeature = ref(['']);
const addFeature = async (areaId: number) => {
  try {
    console.log(`Adding feature to area ${areaId}: ${newFeature.value[areaId]}`);
    const response = await http.post(`/api/v1/features`, {
      'area-id': areaId,
      name: newFeature.value[areaId],
      documentation: '',
      url: '',
      'business-value': ''
    });

    console.log(`Feature added response:`, response.data);
    newFeature.value[areaId] = '';

    // Refresh features for this area
    await getFeatures(areaId);
  } catch (err) {
    console.error(`Error adding feature to area ${areaId}:`, err);
    error.value = `Error adding feature: ${err}`;
  }
};

const removeFeature = async (areaId: number, featureId: number) => {
  try {
    console.log(`Removing feature ${featureId} from area ${areaId}`);
    await http.delete(`/api/v1/features/${featureId}`);

    // Refresh features for this area
    await getFeatures(areaId);
  } catch (err) {
    console.error(`Error removing feature ${featureId} from area ${areaId}:`, err);
    error.value = `Error removing feature: ${err}`;
  }
};

const newFeatureName = ref('');

const showUpdateFeatureModal = (
  featureId: number,
  name: string,
  documentation: string,
  url: string,
  businessValue: string
) => {
  newFeatureName.value = name;
  featureIdToChange.value = featureId;
  featureBusinessValue.value = businessValue;
  featureDocumentation.value = documentation;
  featureUrl.value = url;
  new Modal('#updateFeature').show();
};

const changeFeature = async () => {
  try {
    console.log(`Updating feature ${featureIdToChange.value} with name ${newFeatureName.value}`);
    const response = await http.put(`/api/v1/features/${featureIdToChange.value}`, {
      name: newFeatureName.value,
      documentation: featureDocumentation.value,
      url: featureUrl.value,
      'business-value': featureBusinessValue.value
    });

    console.log('Feature update response:', response.data);

    // Reset form fields
    newFeatureName.value = '';
    featureBusinessValue.value = '';
    featureDocumentation.value = '';
    featureUrl.value = '';
    featureIdToChange.value = 0;

    // Refresh areas and features
    await getAreas();
  } catch (err) {
    console.error(`Error updating feature ${featureIdToChange.value}:`, err);
    error.value = `Error updating feature: ${err}`;
  }
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
