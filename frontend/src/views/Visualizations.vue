<template>
  <v-container>
    <div class="brand mb-1">Visualizations</div>
    <v-tabs v-model="tab">
      <v-tab>Top Cosponsors</v-tab>
      <v-tab>Another Visualization</v-tab>
    </v-tabs>
    <v-tabs-items v-model="tab">
      <v-tab-item>
        <TopCosponsors />
      </v-tab-item>
      <v-tab-item>
        Some Other Visualization
      </v-tab-item>
    </v-tabs-items>
  </v-container>
</template>

<script>
import TopCosponsors from '@/components/TopCosponsors'

export default {
  components: {
    TopCosponsors,
  },

  data: () => ({
    tab: 0,
  }),
}
</script>
