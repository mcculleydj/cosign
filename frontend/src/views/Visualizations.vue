<template>
  <v-container>
    <div class="brand mb-1">Visualizations</div>
    <v-tabs v-model="tab">
      <v-tab>Top Cosponsors</v-tab>
      <v-tab>Subject Graph</v-tab>
    </v-tabs>
    <v-tabs-items v-model="tab">
      <v-tab-item>
        <TopCosponsors :tab="tab" />
      </v-tab-item>
      <v-tab-item>
        <SubjectGraph :tab="tab" />
      </v-tab-item>
    </v-tabs-items>
  </v-container>
</template>

<script>
import TopCosponsors from '@/components/TopCosponsors'
import SubjectGraph from '@/components/SubjectGraph'

export default {
  components: {
    TopCosponsors,
    SubjectGraph,
  },

  data: () => ({
    tab: 0,
  }),
}
</script>
