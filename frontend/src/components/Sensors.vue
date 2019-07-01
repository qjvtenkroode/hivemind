<template>
  <div class="row" id="sensors">
    <div class="collection">
        <a href="#!" class="collection-item" v-for="s in sensors" :key="s.ID">
          {{ s.Name }}
          <span class="badge" :data-badge-caption="s.Unit">{{ s.Value }}</span>
        </a>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'sensors',
    props: {
    },
      
    mounted() {

    },
    created() {
      this.fetchSensors();
      this.timer = setInterval(this.fetchSensors, 5000)
    },
    data() {
      return {
        sensors: [],
        timer: '',
      }
    },
    methods: {
      fetchSensors: function() {
        fetch('http://localhost:5000/api/sensor/')
          .then(response => response.json())
          .then(json => {
            this.sensors = json
          })
      },
      cancelAutoUpdate: function() { clearInterval(this.timer) }
    },
    beforeDestroy() {
      clearInterval(this.timer)
    },
    computed: {

    }
  }
</script>

<style scoped></style>
