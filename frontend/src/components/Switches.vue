<template>
  <div class="row" id="switches">
    <div class="col s3" v-for="sw in switches" :key="sw.ID">
      <div class="card">
        <div class="card-content">
          <center><span class="card-title">{{ sw.ID }}</span></center>
          <div class="card-action">
            <div class="switch">
              <label>
                Off
                <input type="checkbox" v-model="sw.State" :id="sw.ID">
                <span class="lever"></span>
                On
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'switches',
    props: {
    },
      
    mounted() {

    },
    created() {
      this.fetchSwitches();
      this.timer = setInterval(this.fetchSwitches, 5000)
    },
    data: function() {
      return {
        switches: [],
        timer: '',
      }
    },
    methods: {
      fetchSwitches: function() {
        fetch('http://localhost:5000/api/switch/')
          .then(response => response.json())
          .then(json => {
            this.switches = json
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
