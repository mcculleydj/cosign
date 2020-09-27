<template>
  <v-card class="my-3">
    <v-card-title>
      <a :href="bill.link" target="_blank" style="word-break: keep-all">
        HR {{ bill.number }} - {{ bill.title }}
      </a>
    </v-card-title>
    <v-card-text>
      <v-row style="padding: 0px 8px">
        <v-col v-if="bill.policyArea" cols="auto" class="pa-1">
          <v-chip color="primary">Policy Area: {{ bill.policyArea }}</v-chip>
        </v-col>
        <v-col
          v-for="subject in bill.subjects.slice(0, 10)"
          :key="`subject-${bill.number}-${subject}`"
          cols="auto"
          class="pa-1"
        >
          <v-chip>{{ subject }}</v-chip>
        </v-col>
        <v-col v-if="bill.subjects.length > 10" cols="auto" class="pa-1">
          <v-chip color="secondary" @click="showSubjectsDialog = true">
            See all {{ bill.subjects.length }} subjects
          </v-chip>
          <v-dialog v-model="showSubjectsDialog">
            <SubjectsDialog :bill="bill" @close="showSubjectsDialog = false" />
          </v-dialog>
        </v-col>
      </v-row>
      <v-divider class="mt-3" />
      <v-row align="center">
        <v-col>
          <div class="text-right blue--text text--darken-3 display-3">
            {{ bill.numDems }}
          </div>
        </v-col>
        <v-col>
          <div :id="isBillView ? 'svg-seats' : `svg-${bill.number}`" />
        </v-col>
        <v-col>
          <div class="text-left red--text display-3">{{ bill.numReps }}</div>
        </v-col>
      </v-row>
      <template v-if="isBillView">
        <v-divider />
        <v-row>
          <v-col>
            <div
              v-for="dem in demSponsors"
              :key="dem"
              class="text-center blue--text text--darken-3"
            >
              {{ dem }}
            </div>
          </v-col>
          <v-col v-if="otherSponsors.length">
            <div
              v-for="other in otherSponsors"
              :key="other"
              class="text-center green--text"
            >
              {{ other }}
            </div>
          </v-col>
          <v-col>
            <div
              v-for="rep in repSponsors"
              :key="rep"
              class="text-center red--text"
            >
              {{ rep }}
            </div>
          </v-col>
        </v-row>
      </template>
    </v-card-text>
  </v-card>
</template>

<script>
import SubjectsDialog from '@/components/SubjectsDialog'
import { drawChart, clearChart } from '@/d3/seats'

export default {
  components: {
    SubjectsDialog,
  },

  props: ['bill', 'isBillView'],

  computed: {
    numSponsors() {
      return this.bill.sponsors.length + this.bill.cosponsors.length
    },

    demSponsors() {
      if (!this.bill) return []
      return this.bill.sponsors
        .concat(this.bill.cosponsors)
        .filter(s => s.includes(['[D-']))
    },

    repSponsors() {
      if (!this.bill) return []
      return this.bill.sponsors
        .concat(this.bill.cosponsors)
        .filter(s => s.includes(['[R-']))
    },

    otherSponsors() {
      if (!this.bill) return []
      return this.bill.sponsors
        .concat(this.bill.cosponsors)
        .filter(s => s.includes(['[I-']) || s.includes(['[L-']))
    },
  },

  data: () => ({
    showSubjectsDialog: false,
  }),

  mounted() {
    drawChart(
      this.isBillView ? '#svg-seats' : `#svg-${this.bill.number}`,
      this.bill.numDems,
      this.bill.numReps,
      this.bill.numLibs + this.bill.numInds,
    )
  },

  watch: {
    bill() {
      clearChart('#svg-seats')
      drawChart(
        '#svg-seats',
        this.bill.numDems,
        this.bill.numReps,
        this.bill.numLibs + this.bill.numInds,
      )
    },
  },
}
</script>
