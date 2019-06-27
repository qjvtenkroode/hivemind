<template>
  <div id="switches">
    <table>
      <thead>
        <tr>
          <th>Switch ID</th>
          <th>Switch State</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="sw in switches" :key="sw.ID">
          <td>{{ sw.ID }}</td>
          <td>{{ sw.State }}</td>
          <td>
            <button>Toggle</button>
          </td>
        </tr>
      </tbody>
    </table>
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
      fetch('http://localhost:5000/api/switch/')
        .then(response => response.json())
        .then(json => {
          this.switches = json
        })
    },
    data() {
      return {
        switches: []
      }
    },
    methods: {

    },
    computed: {

    }
  }
</script>

<style scoped></style>
