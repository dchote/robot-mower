<template>
    <v-layout align-end justify-center row fill-height>
      <v-card width="700px">
        <v-flex xs12 hidden-md-and-up pl-4 pt-4 pr-4>
          <v-slider v-model="drive.speed" inverse-label label="Speed"></v-slider>
          <v-slider v-model="cutter.speed" inverse-label label="Cutter"></v-slider>
        </v-flex>
        <v-bottom-nav :value="true" :active.sync="movement">
            <v-flex md3 hidden-sm-and-down mr-1 ml-4 mt-2 pt-1>
              <v-slider v-model="drive.speed" inverse-label label="Speed"></v-slider>
            </v-flex>
            <v-layout md3>
              <v-flex xs3>
                <v-btn color="teal" value="left" flat>
                  <span>Left</span>
                  <v-icon>arrow_back</v-icon>
                </v-btn>
              </v-flex>
              <v-flex xs3>
                <v-btn color="teal" value="forward" flat>
                  <span>Forward</span>
                  <v-icon>arrow_upward</v-icon>
                </v-btn>
              </v-flex>
              <v-flex xs3>
                <v-btn color="teal" value="backward" flat>
                  <span>Backward</span>
                  <v-icon>arrow_downward</v-icon>
                </v-btn>
              </v-flex>
              <v-flex xs3>
                <v-btn color="teal" value="right" flat>
                  <span>Right</span>
                  <v-icon>arrow_forward</v-icon>
                </v-btn>
              </v-flex>
            </v-layout>
            <v-flex md3 hidden-sm-and-down ml-1 mr-4 mt-2 pt-1>
              <v-slider v-model="cutter.speed" label="Cutter"></v-slider>
            </v-flex>
        </v-bottom-nav>
      </v-card>
    </v-layout>
</template>

<script>
  import { mapState } from 'vuex'
  
  export default {
    name: 'ControlPage',
    data () {
      return {
        dialog: false,
        
        movement: null,
        speed: 100,
        cutter: 75
      }
    },
    computed: mapState({
      drive: state => state.mower.drive,
      cutter: state => state.mower.cutter,
    })
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .v-bottom-nav {
    opacity: 0.9;
  }
</style>