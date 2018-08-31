import api from '../../api'

// initial state
const state = {
  urls: {
    isFallbackValues: false,
    
    camera: '',
    ws: ''
  }
}

// getters
const getters = {
  cameraBackgroundCSS: state => {
    return "url('" + state.urls.camera + "') no-repeat center center fixed"
  }
}

// actions
const actions = {
  getEndpoints({ commit }) {
    api.getEndpoints(endpoints => {
      commit('setEndpoints', endpoints)
    })
  }
}

// mutations
const mutations = {
  setEndpoints(state, endpoints) {
    state.urls = endpoints.data
    
    console.log('endpoints:', state.urls)
  }
  
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}