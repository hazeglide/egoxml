<template>
  <div class="row pt-1">
    <div class="col-4">
      <b-button-group>
        <b-button @click="parse()" :class="{ disabled : parsebtn.disabled }">Parse</b-button>
        <b-button @click="rebuild()" :class="{ disabled : rebuildbtn.disabled }">Rebuild</b-button>
      </b-button-group>
    </div>
    <div class="col-8">
      <div v-for="(log,index) in logs" :key="index">{{ log }}</div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      parsebtn: { disabled: false },
      rebuildbtn: { disabled: false },
      logs: []
    };
  },
  methods: {
    parse: function() {
      this.parsebtn.disabled = true;
      var self = this;
      axios
        .get('http://localhost:8100/api/parse')
        .then(response => {
          self.parsebtn.disabled = false;
          self.logs.push('parse successful');
        })
        .catch(error => {
          console.log(error)
          self.parsebtn.disabled = false;
          self.logs.push('parse failed');
        });
    },
    rebuild: function() {
      this.$emit('rebuild');
    }
  }
};
</script>

<style>
</style>
