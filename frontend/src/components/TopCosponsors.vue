<template>
  <v-container>
    <v-row>
      <v-col>
        <v-autocomplete
          label="Member"
          v-model="member"
          :items="memberItems"
          outlined
          hide-details
          clearable
        />
      </v-col>
    </v-row>

    <template v-if="member">
      <h3>Top 30 Cosponsors from Another Party</h3>
      <h4 style="white-space: pre">{{ memberToString(cosponsor) || ' ' }}</h4>
      <div class="text-center" style="margin-bottom: -10px">
        <small>Number of Cosponsored Bills</small>
      </div>
    </template>
    <div class="svg-container" />
  </v-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

import { drawChart } from '@/d3/horizontal-bar'
import { memberToString } from '@/common/functions'

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
    ...mapGetters({
      memberItems: 'memberItems',
      memberMap: 'memberMap',
    }),
  },

  data: () => ({
    member: null,
    cosponsor: null,
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

    memberToString,

    setCosponsor(cosponsor) {
      this.cosponsor = cosponsor
    },

    async redrawChart() {
      this.removeFn()
      await this.$nextTick()
      this.removeFn = drawChart(
        '.svg-container',
        {
          member: this.member,
          memberMap: this.memberMap,
        },
        params,
        {
          onMouseover: this.onMouseover,
          onMouseleave: this.onMouseleave,
          onMouseup: this.onMouseup,
        },
      )
    },

    onMouseover(cosponsor) {
      this.cosponsor = cosponsor
    },

    onMouseleave() {
      this.cosponsor = null
    },

    onMouseup(cosponsor) {
      this.$router.push({
        name: 'cosponsors',
        query: { sponsorId: this.member.id, cosponsorId: cosponsor.id },
      })
    },
  },

  watch: {
    member(state) {
      if (!state) {
        this.removeFn()
        this.removeFn = () => {}
        return
      }
      this.redrawChart()
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

<style scoped>
.svg-container {
  display: inline-block;
  position: relative;
  width: 100%;
  padding-bottom: 100%;
  vertical-align: top;
  overflow: hidden;
}

.svg-content {
  display: inline-block;
  position: absolute;
  top: 0;
  left: 0;
}
</style>
