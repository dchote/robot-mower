<template>
  <div class="statsContainer">
    <v-layout align-center justify-end row fill-height wrap>
      <div class="stat black elevation-2 text-xs-center white--text">
        <h5>Battery:</h5>
        <span class="blue-grey--text text--lighten-3">{{ battery.voltage }}</span>
      </div>
      <div class="stat black elevation-2 text-xs-center white--text">
        <h5>Current:</h5>
        <span class="blue-grey--text text--lighten-3">{{ battery.current }}</span>
      </div>
      <div class="stat black elevation-2 text-xs-center white--text">
        <h5>Compass:</h5>
        <span class="blue-grey--text text--lighten-3">{{ compass.bearing }}</span>
      </div>
      <div class="stat black elevation-2 text-xs-center white--text">
        <h5>GPS:</h5>
        <span class="blue-grey--text text--lighten-3">{{ gps.coordinates }}</span>
      </div>
    </v-layout>
  </div>
</template>

<script>
  import { mapState } from 'vuex'
  
  export default {
    name: 'StatusBar',
    computed: mapState({
      battery: state => state.mower.battery,
      compass: state => state.mower.compass,
      gps: state => state.mower.gps,
    })
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .statsContainer {
    position: fixed;
    height: 100px;
    top: 53px;
    right: 0;
  }
  .stat {
    opacity: 0.8;
    min-height: 50px;
    min-width: 80px;
    margin: 10px;
    padding: 10px;
    border-radius: 5px;
  }
</style>