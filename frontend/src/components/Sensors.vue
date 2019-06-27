<template>
  <div id="sensors">
    <table>
      <thead>
        <tr>
          <th>Sensor ID</th>
          <th>Sensor Value</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in sensors" :key="s.ID">
          <td>{{ s.ID }}</td>
          <td>{{ s.Value }}</td>
          <td>
            <button>New value</button>
          </td>
        </tr>
      </tbody>
    </table>
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
      fetch('http://localhost:5000/api/sensor/')
        .then(response => response.json())
        .then(json => {
          this.sensors = json
        })
    },
    data() {
      return {
        sensors: []
      }
    },
    methods: {

    },
    computed: {

    }
  }
</script>

<style scoped></style>
