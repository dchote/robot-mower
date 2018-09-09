<template>
  <div class="statsContainer">
    <v-layout align-end justify-end row fill-height wrap>
      <div class="stat hidden-sm-and-down black elevation-2 text-xs-center white--text">
        <h5>Load Avg:</h5>
        <span class="blue-grey--text text--lighten-3">
          {{ platform.load_average.load1 }}, {{ platform.load_average.load5 }}, {{ platform.load_average.load15 }}
        </span>
      </div>
      <div class="stat hidden-sm-and-down black elevation-2 text-xs-center white--text">
        <h5>CPU:</h5>
        
        <v-progress-linear v-for="cpu in cpuCoreUtilization(platform.cpu)"
              :bind="platform.cpu[cpu.key]"
              :key="cpu.id"
              :size="42"
              :width="2"
              :value="cpu.value"
              color="teal"
        >
        </v-progress-linear>
        
      </div>
      <div class="stat hidden-sm-and-down black elevation-2 text-xs-center white--text">
        <h5>Memory Free:</h5>
        <v-progress-linear
              :size="42"
              :width="2"
              :value="percentFree(platform.memory.total, platform.memory.available)"
              color="teal"
        >
        </v-progress-linear>
        <span class="blue-grey--text text--lighten-3">{{ percentFree(platform.memory.total, platform.memory.available) }}%</span>
      </div>
      <div class="stat hidden-sm-and-down black elevation-2 text-xs-center white--text">
        <h5>Disk Free:</h5>
        <v-progress-linear
              :size="42"
              :width="2"
              :value="percentFree(platform.disk.total, platform.disk.free)"
              color="teal"
              dark
        >
        </v-progress-linear>
        <span class="blue-grey--text text--lighten-3">{{ percentFree(platform.disk.total, platform.disk.free) }}%</span>
      </div>
      <div class="stat black elevation-2 text-xs-center white--text">
        <h5>Battery:</h5>
        <v-progress-linear
              :size="42"
              :width="2"
              :value="batteryPercentFree(battery.voltage, battery.voltage_nominal, battery.voltage_warn)"
              color="teal"
              dark
        >
        </v-progress-linear>
        <span class="blue-grey--text text--lighten-3">{{ battery.voltage }}v ({{ batteryPercentFree(battery.voltage, battery.voltage_nominal, battery.voltage_warn) }}%)</span>
      </div>
      <div class="stat black elevation-2 text-xs-center white--text">
        <h5>Current:</h5>
        <span class="blue-grey--text text--lighten-3">{{ battery.current }}A</span>
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
      platform: state => state.mower.platform,
      battery: state => state.mower.battery,
      compass: state => state.mower.compass,
      gps: state => state.mower.gps,
    }),
    methods: {
      percentFree(total, free) {
        if (total > 0 && free > 0) {
          return Math.round((free / total) * 100)
        }
        return 0
      },
      batteryPercentFree(voltage, nominal, warn) {
        if (voltage) {
          var vSpan = nominal - warn
          var vFree = voltage - warn
          return Math.round((vFree / vSpan) * 100)
        }
        return 0
      },
      cpuCoreUtilization(cpu) {
        if (cpu) {
          var x
          var cores = []
          for (x = 1; x <= cpu.count; x++) {
            var key = 'core_' + x
            cores.push({id: x, value: cpu[key], key: key})
          }
          return cores;
        }
        return null
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .v-progress-circular {
    margin: 1rem
  }
      
  .statsContainer {
    position: fixed;
    top: 60px;
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