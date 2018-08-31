import Vue from 'vue'
import Vuex from 'vuex'

import endpoints from './modules/endpoints'
import mower from './modules/mower'

Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
  plugins: [ ],
  modules: {
    endpoints,
    mower,
  },
  mutations:{
      SOCKET_ONOPEN (state, event)  {
        console.log("ws open", state, event)
      },
      SOCKET_ONCLOSE (state, event)  {
        console.log("ws close", state, event)
      },
      SOCKET_ONERROR (state, event)  {
        console.log("ws error:", state, event)
      },
    },
  strict: debug
})