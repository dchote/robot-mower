import Vue from 'vue'
import VueNativeSock from 'vue-native-websocket'
import axios from 'axios'

import store from '../store'


const endpointsURL = 'http://' + location.hostname + ':8088/v1/endpoints'
//const endpointsURL = 'http://robot-mower.local:8088/v1/endpoints'

const fallbackCameraImage = 'https://media.giphy.com/media/3o6vXRxrhj7Ov94Gbu/source.gif'

export default {
  getEndpoints(callback) {
    axios.get(endpointsURL) // this needs to be dynamic
      .then(response => {
        response.isFallbackValues = false
        
        Vue.use(VueNativeSock, response.data.ws, { store: store, format: 'json' })
        
        callback(response)
      }).catch(function (error) {
        console.log('getEndpoints error:', error)
        
        var fallback = {
          data: {
            isFallbackValues: true,
            camera: fallbackCameraImage, 
            ws: ''
          }}
        
        callback(fallback)
      })
  }
  
  
}