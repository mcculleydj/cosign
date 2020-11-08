<template>
  <v-container>
    <div class="brand mb-1">Visualizations</div>
    <v-tabs v-model="tab">
      <v-tab>Top Cosponsors</v-tab>
      <v-tab>Subject Graphs</v-tab>
    </v-tabs>
    <v-tabs-items v-model="tab">
      <v-tab-item>
        <TopCosponsors />
      </v-tab-item>
      <v-tab-item>
        <SubjectGraphs />
      </v-tab-item>
    </v-tabs-items>
  </v-container>
</template>

<script>
import TopCosponsors from '@/components/TopCosponsors'
import SubjectGraphs from '@/components/SubjectGraphs'

export default {
  components: {
    TopCosponsors,
    SubjectGraphs,
  },

  data: () => ({
    tab: 0,
  }),
}
</script>
