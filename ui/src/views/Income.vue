<template>
  <div>
    <LogControl v-on:rebuild="parseResponses();"></LogControl>
    <Chart :chart="chart"></Chart>
  </div>
</template>

<script>
import Chart from '@/components/Chart.vue';
import LogControl from '@/components/LogControl.vue';
import axios from 'axios';

export default {
  components: {
    Chart,
    LogControl
  },
  data() {
    return {
      raw: [],
      exclude: [],
      combine: {},
      ships: [],
      chart: {
        uuid: 'ship-income-analytics',
        traces: [],
        layout: {
          dragmode: 'pan',
          xaxis: { title: 'hours played', range: [60, 100] },
          yaxis: { title: 'Credits', fixedrange: true }
        }
      }
    };
  },
  methods: {
    fetchLog: function() {
      return axios.get(`http://localhost:8100/api/log`);
    },
    fetchExcluded: function() {
      return axios.get(`http://localhost:8100/api/excludes`);
    },
    fetchCombined: function() {
      return axios.get(`http://localhost:8100/api/combines`);
    },
    patchExcluded: function() {
      var self = this;
      Object.keys(this.combine).forEach(value => {
        self.combine[value].forEach(v => {
          self.exclude.push(v);
          self.ships.splice(self.ships.indexOf(v), 1);
        });
      });
    },
    parseResponses: function() {
      var self = this;
      axios
        .all([this.fetchLog(), this.fetchExcluded(), this.fetchCombined()])
        .then(
          axios.spread(function(log, excl, comb) {
            self.raw = log.data.sort(function(a, b) {
              return a.Time - b.Time;
            });
            self.exclude = excl.data;
            self.combine = comb.data;
          })
        )
        .then(response => {
          self.patchExcluded();
        })
        .then(response => {
          this.ships = [
            ...new Set(
              this.raw
                .map(value => {
                  return value.Ship;
                })
                .filter(val => val !== '')
            )
          ].filter(val => !this.exclude.includes(val));
          this.chart.traces = [];
          var latestTime = this.raw[this.raw.length - 1].Time;
          this.chart.layout.xaxis.range[0] = latestTime - 20;
          this.chart.layout.xaxis.range[1] = latestTime;
          this.ships.forEach(currShip => {
            var extras = [];
            if (this.exclude.includes(currShip)) {
              return;
            }
            extras.push(currShip);
            if (Object.keys(this.combine).includes(currShip.toLowerCase())) {
              this.combine[currShip.toLowerCase()].forEach(val => extras.push(val));
            }
            var trades = this.raw
              .filter(val => extras.includes(val.Ship))
              .sort(function(a, b) {
                return a.Time - b.Time;
              });
            var trace = {
              type: 'scatter',
              x: trades.map(value => {
                return Math.round(value.Time * 100) / 100;
              }),
              y: [],
              name: currShip
            };
            trades
              .map(value => value.Money)
              .reduce(function(accumulator, currentValue) {
                trace.y.push((Math.round(accumulator + currentValue) / 100) * 100);
                return accumulator + currentValue;
              }, 0);
            this.chart.traces.push(trace);
          });
        });
    }
  },
  mounted: function() {
    this.parseResponses();
  }
};
</script>

<style>
</style>
