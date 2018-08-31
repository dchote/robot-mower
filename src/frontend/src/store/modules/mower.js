// initial state
const state = {
  battery: {
    status: null,
    voltage: 'unknown',
    current: 'unknown'
  },
  
  compass: {
    status: null,
    bearing: 'unknown'
  },
  
  gps: {
    status: null,
    coordinates: 'unknown'
  },
  
  
  drive: {
    speed: 100,
    direction: null
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
const mutations = {
  setMowerDriveSpeed(state, value) {
    state.drive.speed = value
  },
  setDirection(state, value) {
    state.drive.direction = value
  },
  setMowerCutterSpeed(state, value) {
    state.cutter.speed = value
  },
  setMowerState(state, event) {
    state.battery = event.battery
    state.compass = event.compass
    state.gps = event.gps
    state.drive = event.drive
    state.cutter = event.cutter
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}