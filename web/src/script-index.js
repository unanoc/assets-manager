// Collection of dependency libs, to be bundled

import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'

// Import Bootstrap an BootstrapVue CSS files (order is important)
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import * as assetsLib from "../assetsLib/dist/assetsLib";

const octokit = require("@octokit/request");
export const UIkit = require('uikit');
import Icons from 'uikit/dist/js/uikit-icons';

export const assets = assetsLib;
export const octoRequest = octokit.request;

// Make BootstrapVue available throughout your project
Vue.use(BootstrapVue)
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin)

// loads the Icon plugin
UIkit.use(Icons);
// components can be called from the imported UIkit reference
//UIkit.notification('Hello world.');
