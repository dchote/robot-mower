// initial state
const state = {
  battery: {
    status: null,
    voltage: '24.3',
    current: '1.4'
  },
  
  compass: {
    status: null,
    bearing: 'NE'
  },
  
  gps: {
    status: null,
    coordinates: '40.780715, -78.007729'
  },
  
  
  drive: {
    speed: 100
  },
  
  cutter: {
    speed: 0
  },
}

// getters
const getters = {}

// actions
const actions = {}

// mutations
const mutations = {}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}