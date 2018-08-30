import axios from 'axios'

const endpointsURL = 'http://robot-mower.local:8088/v1/endpoints'

export default {
  getEndpoints(callback) {
    axios.get(endpointsURL) // this needs to be dynamic
      .then(response => {
        callback(response)
      })
  }
  
  
}