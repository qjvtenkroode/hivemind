<template>
  <div id="app" class="container">
      <h1>Hivemind's drones</h1>

      <switches />
      <sensors />
  </div>
</template>

<script>
import Switches from '@/components/Switches.vue'
import Sensors from '@/components/Sensors.vue'

export default {
  name: 'app',
  components: {
    Switches,
    Sensors,
  },
}
</script>

<style>
</style>
