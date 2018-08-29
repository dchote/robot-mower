<template>
    <v-layout align-end justify-center row fill-height>
      <v-card height="70px" width="400px" flat>
        <v-bottom-nav absolute color="transparent" :value="true">
        <v-btn color="teal" flat>
          <span>Left</span>
          <v-icon>arrow_back</v-icon>
        </v-btn>
          <v-btn color="teal" flat>
            <span>Forward</span>
            <v-icon>arrow_upward</v-icon>
          </v-btn>
          <v-btn color="teal" flat>
            <span>Backward</span>
            <v-icon>arrow_downward</v-icon>
          </v-btn>
          <v-btn color="teal" flat>
            <span>Right</span>
            <v-icon>arrow_forward</v-icon>
          </v-btn>
        </v-bottom-nav>
      </v-card>
    </v-layout>
</template>

<script>
  export default {
    name: 'ControlPage',
    data () {
      return {
      
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>