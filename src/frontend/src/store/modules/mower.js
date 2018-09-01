// initial state
const state = {
  platform: {
    hostname: null,
    operating_system: null,
    platform: null,
    load_average: {
      load1: null,
      load5: null,
      load15: null
    },
    memory: {
      total: null,
      available: null
    },
    disk: {
      total: null,
      free: null
    }
  },
  
  battery: {
    status: null,
    voltage_nominal: null,
    voltage_warn: null,
    voltage: null,
    current: null
  },
  
  compass: {
    status: null,
    bearing: null
  },
  
  gps: {
    status: null,
    coordinates: null
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
    state.platform = event.platform
    state.battery = event.battery
    state.compass = event.compass
    state.gps = event.gps
    state.drive = event.drive
    state.cutter = event.cutter
    
    console.log(event)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}