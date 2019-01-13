<template>
  <div class="row pt-3">
    <div class="col-2">
      <select v-model="currentShip" class="ships-select">
        <option
          v-for="(ship,index) in ships"
          :key="ship"
          :value="index"
          :selected=" index == currentShip"
        >{{ ship }}</option>
      </select>
    </div>
    <div class="col-10"></div>
    <div class="col-12">
      <Chart :chart="chart"></Chart>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import Chart from '@/components/Chart.vue';

export default {
  components: {
    Chart
  },
  data() {
    return {
      raw: [],
      exclude: [],
      combine: {},
      ships: [],
      currentShip: -1,
      chart: {
        uuid: 'single-ship-income-analytics',
        traces: [],
        layout: {
          dragmode: 'pan',
          separators: ',.',
          hoverformat: ',.0~df',
          tickformat: ',.0~df',
          title: 'Average Income / Hour',
          xaxis: { title: 'Hour' },
          yaxis: { title: 'Credits', fixedrange: true }
        }
      }
    };
  },
  methods: {
    init: function() {
      var self = this;
      axios
        .all([this.fetchLog(), this.fetchExcluded(), this.fetchCombined()])
        .then(
          axios.spread(function(log, excl, comb) {
            self.raw = log.data;
            self.exclude = excl.data;
            self.combine = comb.data;
          })
        )
        .then(response => {
          self.patchExcluded();
          self.currentShip = 0;
        });
    },
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
    }
  },
  mounted: function() {
    this.init();
  },
  watch: {
    raw: function(value) {
      this.ships = [...new Set(value.map(entry => entry.Ship))].filter(
        val => !this.exclude.includes(val)
      );
    },
    currentShip: function(value) {
      var extras = [];
      var currShip = this.ships[this.currentShip];
      extras.push(currShip);
      if (Object.keys(this.combine).includes(currShip.toLowerCase())) {
        this.combine[currShip.toLowerCase()].forEach(val => extras.push(val));
      }
      var shipLogs = this.raw.filter(val => extras.includes(val.Ship));
      var min = Math.floor(Math.min(...shipLogs.map(val => val.Time)));
      var max = Math.ceil(Math.max(...shipLogs.map(val => val.Time)));
      this.chart.layout.xaxis.range = [min, max];
      var trace = {
        type: 'scatter',
        x: [...new Set(shipLogs.map(val => Math.floor(val.Time)))],
        y: [],
        line: { shape: 'hvh' },
        name: 'average per hour'
      };

      for (var x in trace.x) {
        var i = trace.x[x];
        var tmp = shipLogs.filter(val => val.Time >= i && val.Time < i + 1);
        trace.y.push(
          tmp.map(val => Math.round(val.Money)).reduce((pv, cv) => pv + cv, 0) / tmp.length
        );
      }
      this.chart.traces = [trace];
    }
  }
};
</script>

<style>
</style>
