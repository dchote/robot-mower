import api from '../../api'

// initial state
const state = {
  endpoints: {
    camera: "",
    ws: ""
  }
}

// getters
const getters = {}

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
    state.endpoints = endpoints.data
    
    console.log("endpoints:", state.endpoints)
  }
  
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}