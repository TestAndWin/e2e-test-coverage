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
          <div v-if="areaToggleMap.get(area['id'] ?? 0)">
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
          <form @submit.prevent>
            <div class="mb-3">
              <label>Name</label><input type="text" class="form-control" id="newAreaName" v-model="newAreaName" />
              <label>Please be aware that the name is used to "match" the test results.</label>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary pointer" data-bs-dismiss="modal">Close</button>
              <button type="button" class="btn btn-primary pointer" data-bs-dismiss="modal" @click="changeAreaName">
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
          <form @submit.prevent>
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
              <button type="button" class="btn btn-primary pointer" @click="changeFeatureAndCloseModal">
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
    const response = await http.get(`/api/v1/products`);

    // Extract data from StandardResponse format
    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      products.value = response.data.data;
    } else {
      products.value = [];
    }
  } catch (err) {
    error.value = `Error loading products: ${err}`;
    products.value = []; // Initialize as empty array on error
  }
  loading.value = false;
};

const newProduct = ref('');
const addProduct = async () => {
  try {
    // Check if we have a valid name
    if (!newProduct.value || newProduct.value.trim() === '') {
      error.value = 'Product name cannot be empty';
      return;
    }

    // Make the API call
    await http.post(`/api/v1/products`, {
      name: newProduct.value
    });

    // Clear input
    newProduct.value = '';

    // Refresh product list
    await getProducts();
  } catch (err) {
    error.value = `Error adding product: ${err}`;
  }
};

// Areas
const areas = ref<Area[]>([]);
// Use a Map instead of an array to store toggle states by area ID
const areaToggleMap = ref(new Map<number, boolean>());
const getAreas = async () => {
  loading.value = true;

  try {
    const response = await http.get(`/api/v1/products/${props.productId}/areas`);

    // Extract data from StandardResponse format
    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      areas.value = response.data.data;
    } else {
      areas.value = [];
    }

    // Clean up stale entries from areaToggleMap
    // Get a set of current area IDs
    const currentAreaIds = new Set(areas.value.map((a) => a['id'] ?? 0));

    // Remove any areas from the toggle map that no longer exist
    for (const areaId of areaToggleMap.value.keys()) {
      if (!currentAreaIds.has(areaId)) {
        areaToggleMap.value.delete(areaId);
      }
    }

    // Safely iterate over areas if it's an array
    if (Array.isArray(areas.value)) {
      areas.value.forEach((a) => {
        const areaId = a['id'] ?? 0;
        getFeatures(areaId);

        // Initialize toggle state for new areas if not already set
        if (!areaToggleMap.value.has(areaId)) {
          areaToggleMap.value.set(areaId, false);
        }
      });
    } else {
      areas.value = []; // Ensure it's always an array
    }
  } catch (err: any) {
    const errorMsg = err.response?.data?.error || err.message || String(err);
    error.value = `Error loading areas: ${errorMsg}`;
    areas.value = []; // Initialize as empty array on error
  }

  loading.value = false;
};

const newArea = ref('');
const addArea = async () => {
  try {
    // Check if we have a valid name
    if (!newArea.value || newArea.value.trim() === '') {
      error.value = 'Area name cannot be empty';
      return;
    }

    // Make the API call
    await http.post(`/api/v1/areas`, {
      'product-id': props.productId,
      name: newArea.value
    });

    // Clear input
    newArea.value = '';

    // Refresh area list
    await getAreas();
  } catch (err) {
    error.value = `Error adding area: ${err}`;
  }
};

const removeArea = async (areaId: number) => {
  try {
    // Make the API call to delete from database
    await http.delete(`/api/v1/areas/${areaId}`);

    // On successful deletion, reload the current page to ensure all components are properly reset
    // This forces a complete refresh of the UI state, clearing any stale references
    window.location.reload();
  } catch (err: any) {
    // Only show error for the actual deletion operation
    const errorMsg = err.response?.data?.error || err.message || String(err);
    error.value = `Error deleting area: ${errorMsg}`;
  }
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
  // Get current state (default to false if not set)
  const currentState = areaToggleMap.value.get(areaId) || false;
  // Toggle the state
  areaToggleMap.value.set(areaId, !currentState);
};

// Features
const features = ref([[]]);
const getFeatures = async (areaId: number) => {
  try {
    const response = await http.get(`/api/v1/areas/${areaId}/features`);

    // Extract data from StandardResponse format
    let featureData = [];

    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      featureData = response.data.data;
    } else {
      featureData = []; // Empty array as fallback
    }

    // Update the features array
    features.value[areaId] = featureData;
  } catch (err: any) {
    // Don't show error for 404 responses after deletion
    if (err.response && err.response.status === 404) {
      // Area was likely deleted, just clear the features
      features.value[areaId] = [];
    } else {
      const errorMsg = err.response?.data?.error || err.message || String(err);
      error.value = `Error loading features: ${errorMsg}`;
      features.value[areaId] = []; // Empty array on error
    }
  }
};

const newFeature = ref(['']);
const addFeature = async (areaId: number) => {
  try {
    await http.post(`/api/v1/features`, {
      'area-id': areaId,
      name: newFeature.value[areaId],
      documentation: '',
      url: '',
      'business-value': ''
    });

    newFeature.value[areaId] = '';

    // Refresh features for this area
    await getFeatures(areaId);
  } catch (err) {
    error.value = `Error adding feature: ${err}`;
  }
};

const removeFeature = async (areaId: number, featureId: number) => {
  try {
    await http.delete(`/api/v1/features/${featureId}`);

    // Refresh features for this area
    await getFeatures(areaId);
  } catch (err) {
    error.value = `Error removing feature: ${err}`;
  }
};

const newFeatureName = ref('');
// Need to track which area a feature belongs to
const featureAreaId = ref(0);

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

  // Find which area this feature belongs to
  for (const areaId in features.value) {
    if (!features.value[areaId]) continue;

    const featureIndex = features.value[areaId].findIndex((f: any) => f.id === featureId);
    if (featureIndex !== -1) {
      featureAreaId.value = parseInt(areaId);
      break;
    }
  }

  // Store the modal instance so we can close it programmatically
  updateFeatureModal = new Modal('#updateFeature');
  updateFeatureModal.show();
};

const changeFeature = async () => {
  try {
    const response = await http.put(`/api/v1/features/${featureIdToChange.value}`, {
      name: newFeatureName.value,
      documentation: featureDocumentation.value,
      url: featureUrl.value,
      'business-value': featureBusinessValue.value
    });

    // Only refresh features for the specific area that contains this feature
    if (featureAreaId.value > 0) {
      // Ensure the area stays open
      if (!areaToggleMap.value.has(featureAreaId.value)) {
        areaToggleMap.value.set(featureAreaId.value, true);
      } else if (!areaToggleMap.value.get(featureAreaId.value)) {
        // If area is closed, force it open as the user was working on a feature inside it
        areaToggleMap.value.set(featureAreaId.value, true);
      }

      // Refresh just the features for this area
      await getFeatures(featureAreaId.value);
    }

    // Reset form fields
    newFeatureName.value = '';
    featureBusinessValue.value = '';
    featureDocumentation.value = '';
    featureUrl.value = '';
    featureIdToChange.value = 0;
    featureAreaId.value = 0;
  } catch (err) {
    error.value = `Error updating feature: ${err}`;
  }
};

const closeAlert = () => {
  error.value = '';
};

// Reference to the feature update modal
let updateFeatureModal: any = null;

// Function to handle both changing feature and closing modal
const changeFeatureAndCloseModal = async () => {
  // First ensure the area toggle is explicitly set to open
  if (featureAreaId.value > 0) {
    areaToggleMap.value.set(featureAreaId.value, true);
  }

  // Change the feature
  await changeFeature();

  // Close the modal programmatically
  if (updateFeatureModal) {
    updateFeatureModal.hide();
  }
};

onMounted(() => {
  // Initialize the map
  areaToggleMap.value = new Map<number, boolean>();

  getProducts();
  getAreas();
});
</script>

<style scoped>
@import '../assets/styles.css';
</style>
