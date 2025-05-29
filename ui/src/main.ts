import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import 'bootstrap/dist/css/bootstrap.css';
import { fetchCurrentUser } from './stores/user';

const app = createApp(App);
app.use(router);

// Load user info on startup to populate roles
fetchCurrentUser();

app.mount('#app');

import 'bootstrap/dist/js/bootstrap.js';
import 'bootstrap-icons/font/bootstrap-icons.css';
