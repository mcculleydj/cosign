<template>
  <v-container>
    <div class="brand">Visualizations</div>
    <v-row>
      <v-col>
        <v-autocomplete
          label="Member"
          v-model="member"
          :items="memberItems"
          solo
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="auto">
        <v-select
          label="Threshold"
          :items="thresholds"
          v-model="threshold"
          hide-details
          outlined
          dense
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <div id="d3-container" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import { drawChart } from '@/d3/radial-bar'

const params = {
  width: 700,
  height: 700,
  innerRadius: 50,
  displayFontSize: 22,
  bodyFontSize: 10,
}

export default {
  computed: {
    ...mapGetters({
      memberItems: 'memberItems',
      memberMap: 'memberMap',
    }),
  },

  data: () => ({
    member: null,
    cosponsor: null,
    threshold: 40,
    thresholds: [0, 10, 20, 30, 40, 50],
    removeFn: () => {},
  }),

  created() {
    if (!this.memberItems.length) {
      this.dispatchGetMembers()
    }
  },

  methods: {
    ...mapActions({
      dispatchGetMembers: 'getMembers',
    }),

    setCosponsor(cosponsor) {
      this.cosponsor = cosponsor
    },

    async renderChart() {
      this.removeFn()
      await this.$nextTick()
      this.removeFn = drawChart(
        '#d3-container',
        this.member,
        this.memberMap,
        {
          threshold: this.threshold,
          ...params,
        },
        {
          setCosponsor: this.setCosponsor,
        },
      )
    },
  },

  watch: {
    member() {
      this.renderChart()
    },

    threshold() {
      this.renderChart()
    },
  },
}
</script>

<style scoped>
svg {
  border: 1px solid black;
  border-radius: 5px;
  height: 500px;
  width: 100%;
}
</style>

<style>
.active {
  fill: goldenrod;
  opacity: 0.5;
}
</style>
