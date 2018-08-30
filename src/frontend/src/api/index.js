import axios from 'axios'

export default {
  getEndpoints(callback) {
    axios.get("http://localhost:8088/v1/endpoints")
      .then(response => {
        callback(response)
      })
  }
  
  
}