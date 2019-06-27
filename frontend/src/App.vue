<template>
  <div id="app" class="small-container">
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
  button {
    background: #009435;
    border: 1px solid #009435;
  }

  .small-container {
    max-width: 680px;
  }
</style>
