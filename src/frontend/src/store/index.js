import Vue from 'vue'
import Vuex from 'vuex'

import mower from './modules/mower'

Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

// A Vuex instance is created by combining the state, mutations, actions,
// and getters.
export default new Vuex.Store({
  modules: {
    mower,
  },
  strict: debug,
  plugins: []
})