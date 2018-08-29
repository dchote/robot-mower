<template>
  <v-layout align-center justify-center row fill-height>
      <section>
        <v-layout column wrap align-center grey lighten-4 elevation-4>
          <v-flex xs12 sm4 pa-5 class="my-3">
            <div class="text-xs-center">
              <h2 class="headline">Welcome to your Robot Mower</h2>
              <span class="subheading">
                To control your mower, please make a selection from the menu on the left.
              </span>
              <img src="../assets/Mower_Chassis_v35.png" width="70%" class="mt-4">
            </div>
          </v-flex>
        </v-layout>
      </section>
  </v-layout>
</template>

<script>
  export default {
    name: 'Welcome',
    data: () => ({
      
    })
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>