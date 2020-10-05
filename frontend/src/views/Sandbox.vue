<template>
  <v-container>
    <div class="brand">Sandbox</div>
    <h1>{{ member ? member.name : 'Null' }}</h1>
    <h1>{{ cosponsor ? cosponsor.name : 'Null' }}</h1>
    <div id="d3-container"></div>
  </v-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import { drawChart } from '@/d3/horizontal-bar'

const params = {
  width: 800,
  height: 600,
  margin: {
    top: 30,
    bottom: 30,
    left: 10,
    right: 30,
  },
  cutoff: 30,
}

export default {
  computed: {
    ...mapGetters(['memberItems', 'memberMap']),

    member() {
      return this.memberItems.length ? this.memberItems[0].value : null
    },
  },

  data: () => ({
    cosponsor: null,
  }),

  async mounted() {
    await this.dispatchGetMembers()
    const raw = { member: this.member, memberMap: this.memberMap }
    drawChart('#d3-container', raw, params, {
      onMouseover: this.onMouseover,
      onMouseleave: this.onMouseleave,
    })
  },

  methods: {
    ...mapActions({
      dispatchGetMembers: 'getMembers',
    }),

    onMouseover(member) {
      this.cosponsor = member
    },

    onMouseleave() {
      this.cosponsor = null
    },
  },
}
</script>

<style>
.active {
  fill: goldenrod;
}
</style>
