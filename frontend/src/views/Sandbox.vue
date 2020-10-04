<template>
  <v-container>
    <div class="brand">Sandbox</div>
    <div id="d3-container"></div>
    <pre>{{ cosponsor }}</pre>
  </v-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import { drawChart } from '@/d3/radial-bar'

// const data = [
//   { quarter: '3/31/20', revenue: 75452, net: 2535 },
//   { quarter: '12/31/19', revenue: 87437, net: 3268 },
//   { quarter: '9/30/19', revenue: 69981, net: 2134 },
//   { quarter: '6/30/19', revenue: 63404, net: 2625 },
//   { quarter: '3/31/19', revenue: 59700, net: 3561 },
//   { quarter: '12/31/18', revenue: 72383, net: 3027 },
//   { quarter: '9/30/18', revenue: 56576, net: 2883 },
//   { quarter: '6/30/18', revenue: 52886, net: 2534 },
//   { quarter: '3/31/18', revenue: 51042, net: 1629 },
//   { quarter: '12/31/17', revenue: 60453, net: 1856 },
//   { quarter: '9/30/17', revenue: 43744, net: 256 },
//   { quarter: '6/30/17', revenue: 37955, net: 197 },
//   { quarter: '3/31/17', revenue: 35714, net: 724 },
//   { quarter: '12/31/16', revenue: 43741, net: 749 },
//   { quarter: '9/30/16', revenue: 32714, net: 252 },
//   { quarter: '6/30/16', revenue: 30404, net: 857 },
//   { quarter: '3/31/16', revenue: 29128, net: 513 },
//   { quarter: '12/31/15', revenue: 35746, net: 482 },
//   { quarter: '9/30/15', revenue: 25358, net: 79 },
//   { quarter: '6/30/15', revenue: 23185, net: -57 },
// ]

const params = {
  width: 800,
  height: 800,
  innerRadius: 50,
  displayFontSize: 22,
  bodyFontSize: 20,
  threshold: 40,
}

export default {
  computed: {
    ...mapGetters(['memberItems', 'memberMap']),

    member() {
      return this.memberItems[0].value
    },
  },

  data: () => ({
    cosponsor: null,
  }),

  async mounted() {
    await this.dispatchGetMembers()
    drawChart('#d3-container', this.member, this.memberMap, params, {
      setCosponsor: this.setCosponsor,
    })
  },

  methods: {
    ...mapActions({
      dispatchGetMembers: 'getMembers',
    }),

    setCosponsor(cosponsor) {
      this.cosponsor = cosponsor
    },
  },
}
</script>

<style>
.active {
  fill: goldenrod;
  opacity: 0.5;
}
</style>
