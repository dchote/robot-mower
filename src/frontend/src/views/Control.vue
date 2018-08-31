<template>
    <v-layout align-end justify-center row fill-height>
      <v-card width="700px">
        <v-flex xs12 hidden-md-and-up pl-5 pt-4 pr-5>
          <v-slider v-model="driveSpeed" inverse-label label="Speed"></v-slider>
          <v-slider v-model="cutterSpeed" inverse-label label="Cutter"></v-slider>
        </v-flex>
        <v-bottom-nav :value="true" :active.sync="direction">
            <v-flex md3 hidden-sm-and-down mr-1 ml-4 mt-2 pt-1>
              <v-slider v-model="driveSpeed" inverse-label label="Speed"></v-slider>
            </v-flex>
            <v-layout md3>
              <v-flex xs3>
                <v-btn color="teal" value="left" flat v-on:mousedown="requestDirectionStart('left')"  v-on:mouseup="requestDirectionStop('left')">
                  <span>Left</span>
                  <v-icon>arrow_back</v-icon>
                </v-btn>
              </v-flex>
              <v-flex xs3>
                <v-btn color="teal" value="forward" flat v-on:mousedown="requestDirectionStart('forward')"  v-on:mouseup="requestDirectionStop('forward')">
                  <span>Forward</span>
                  <v-icon>arrow_upward</v-icon>
                </v-btn>
              </v-flex>
              <v-flex xs3>
                <v-btn color="teal" value="backward" flat v-on:mousedown="requestDirectionStart('backward')"  v-on:mouseup="requestDirectionStop('backward')">
                  <span>Backward</span>
                  <v-icon>arrow_downward</v-icon>
                </v-btn>
              </v-flex>
              <v-flex xs3>
                <v-btn color="teal" value="right" flat v-on:mousedown="requestDirectionStart('right')"  v-on:mouseup="requestDirectionStop('right')">
                  <span>Right</span>
                  <v-icon>arrow_forward</v-icon>
                </v-btn>
              </v-flex>
            </v-layout>
            <v-flex md3 hidden-sm-and-down ml-1 mr-4 mt-2 pt-1>
              <v-slider v-model="cutterSpeed" label="Cutter"></v-slider>
            </v-flex>
        </v-bottom-nav>
      </v-card>
    </v-layout>
</template>

<script>  
  export default {
    name: 'ControlPage',
    data () {
      return {
        dialog: false,
      }
    },
    computed: {
      driveSpeed: {
        get() {
          return this.$store.state.mower.drive.speed
        },
        set(value) {
          this.$store.commit('mower/setMowerDriveSpeed', value)
          this.$socket.sendObj({'method': 'setMowerDriveSpeed', 'value': value.toString()})
        }
      },
      cutterSpeed: {
        get() {
          return this.$store.state.mower.cutter.speed
        },
        set(value) {
          this.$store.commit('mower/setMowerCutterSpeed', value)
          this.$socket.sendObj({'method': 'setMowerCutterSpeed', 'value': value.toString()})
        }
      },
      direction: {
        get() {
          return this.$store.state.mower.drive.direction
        }
      }
    },
    methods: {
      requestDirectionStart(value) {
        this.$socket.sendObj({'method': 'requestDirectionStart', 'value': value})
      },
      requestDirectionStop(value) {
        this.$socket.sendObj({'method': 'requestDirectionStop', 'value': value})
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .v-bottom-nav {
    opacity: 0.9;
  }
</style>