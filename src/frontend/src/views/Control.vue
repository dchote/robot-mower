<template>
    <v-layout align-end justify-center row fill-height>
      <v-card height="70px" width="700px" flat>
        <v-bottom-nav absolute color="transparent" :value="true">
            <v-flex xs3 mr-1 ml-4 mt-2>
              <v-slider v-model="speed" inverse-label label="Speed"></v-slider>
            </v-flex>
            <v-layout xs3>
              <v-flex xs3>
                <v-btn color="teal" flat>
                  <span>Left</span>
                  <v-icon>arrow_back</v-icon>
                </v-btn>
              </v-flex>
              <v-flex xs3>
                <v-btn color="teal" flat>
                  <span>Forward</span>
                  <v-icon>arrow_upward</v-icon>
                </v-btn>
              </v-flex>
              <v-flex xs3>
                <v-btn color="teal" flat>
                  <span>Backward</span>
                  <v-icon>arrow_downward</v-icon>
                </v-btn>
              </v-flex>
              <v-flex xs3>
                <v-btn color="teal" flat>
                  <span>Right</span>
                  <v-icon>arrow_forward</v-icon>
                </v-btn>
              </v-flex>
            </v-layout>
            <v-flex xs3 ml-1 mr-4 mt-2>
              <v-slider v-model="cutter" label="Cutter"></v-slider>
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
        speed: 100,
        cutter: 75
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>