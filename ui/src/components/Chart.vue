<template>
  <div :ref="chart.uuid"></div>
</template>

<script>
import Plotly from 'plotly.js-dist';
export default {
  props: {
    chart: [Object]
  },
  mounted() {
    Plotly.plot(this.$refs[this.chart.uuid], this.chart.traces, this.chart.layout);
  },
  watch: {
    chart: {
      handler: function() {
        Plotly.react(this.$refs[this.chart.uuid], this.chart.traces, this.chart.layout);
        Plotly.relayout(this.$refs[this.chart.uuid], { 'xaxis.range': this.chart.layout.xaxis.range })
      },
      deep: true
    }
  }
};
</script>

<style>
</style>
