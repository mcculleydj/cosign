<template>
  <v-container class="px-0">
    <v-row class="px-3">
      <v-col class="large-text">
        This is a responsive bar chart that considers a data access pattern
        beginning with a particular member. Once the user chooses a member, the
        top cospsonsors from another party are displayed in descending order. By
        clicking the bar the user can route to the cosponsor view where details
        for each cosponsored bill are listed. This provides a method of quickly
        identifying who is most likely to work with a given member and what
        kinds of leglislation they tend to cosponsor together.
      </v-col>
    </v-row>
    <v-row class="px-3">
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
    <div v-show="member" id="top-cosponsors-svg-container" />
  </v-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import { fromEvent } from 'rxjs'
import { filter, debounceTime, tap } from 'rxjs/operators'

import Bar from '@/d3/bar'
import { sleep } from '@/common/functions'

const numCosponsors = 10
const margins = {
  top: 30,
  bottom: 30,
  left: 30,
  right: 30,
}

function transformData(member, memberMap) {
  return Object.entries(member.counts)
    .map(([id, count]) => ({ value: count, ...memberMap[id] }))
    .sort((m1, m2) => m2.value - m1.value)
    .slice(0, numCosponsors)
}

// TODO: would be cool if you could also get a sense of what subjects / policy areas
//       most often appear in that set of bills

export default {
  props: ['tab'],

  computed: {
    ...mapGetters({
      memberItems: 'memberItems',
      memberMap: 'memberMap',
    }),
  },

  data: () => ({
    member: null,
    barChart: null,
  }),

  subscriptions() {
    const resize$ = fromEvent(window, 'resize').pipe(
      filter(() => this.tab === 1 && this.member),
      debounceTime(500),
      tap(() => {
        this.barChart.update(
          transformData(this.member, this.memberMap),
          'name',
          this.onClick,
        )
      }),
    )

    return { resize$ }
  },

  created() {
    if (!this.memberItems.length) {
      this.dispatchGetMembers()
    }
  },

  mounted() {
    this.barChart = new Bar('top-cosponsors-svg-container', margins)
  },

  methods: {
    ...mapActions({
      dispatchGetMembers: 'getMembers',
    }),

    onClick(cosponsor) {
      this.$router.push({
        name: 'cosponsors',
        query: { sponsorId: this.member.id, cosponsorId: cosponsor.id },
      })
    },

    async waitForDOM() {
      const container = document.getElementById('top-cosponsors-svg-container')
      while (!container || !container.clientWidth || !container.clientHeight) {
        await sleep(10)
      }
    },
  },

  watch: {
    async member() {
      if (this.member) {
        await this.waitForDOM()
        this.barChart.update(
          transformData(this.member, this.memberMap),
          'name',
          this.onClick,
        )
      } else {
        this.barChart.update([], 'name', this.onClick)
      }
    },
  },
}
</script>

<style scoped>
#top-cosponsors-svg-container {
  display: inline-block;
  position: relative;
  width: 100%;
  vertical-align: top;
  height: 600px;
}

.large-text {
  font-size: 1.1rem;
}
</style>
