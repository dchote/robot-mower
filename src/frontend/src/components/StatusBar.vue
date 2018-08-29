<template>
  <div class="statsContainer">
    <v-layout align-center justify-end row fill-height>
      <div class="stat grey lighten-4 elevation-2 text-xs-center">
        <h5>Battery Voltage:</h5>
        {{ voltage }}
      </div>
      <div class="stat grey lighten-4 elevation-2 text-xs-center">
        <h5>Compass:</h5>
        {{ compass }}
      </div>
      <div class="stat grey lighten-4 elevation-2 text-xs-center">
        <h5>GPS:</h5>
        {{ gps }}
      </div>
    </v-layout>
  </div>
</template>

<script>
  export default {
    name: 'StatusBar',
    data () {
      return {
        voltage: '24.1',
        compass: 'NE',
        gps: '40.780715, -78.007729'
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .statsContainer {
    position: fixed;
    height: 100px;
    top: 54px;
    right: 0;
  }
  .stat {
    min-height: 50px;
    min-width: 80px;
    margin: 10px;
    padding: 10px;

  }
</style>