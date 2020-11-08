<template>
  <v-container>
    <v-row>
      <v-col>
        <v-autocomplete
          label="Policy Area"
          v-model="selectedPolicyAreas"
          :items="policyAreaItems"
          outlined
          hide-details
          multiple
          small-chips
          clearable
          @change="getCells()"
        />
      </v-col>
      <v-col>
        <v-autocomplete
          label="Subject"
          v-model="selectedSubjects"
          :items="subjectItems"
          outlined
          hide-details
          multiple
          small-chips
          clearable
          @change="getCells()"
        />
      </v-col>
    </v-row>
    <h3 v-if="selectedPolicyAreas.length > 0 || selectedSubjects.length > 0">
      Graph
    </h3>
    <div class="view-box">
      <div class="svg-container" />
    </div>
  </v-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  computed: {
    ...mapGetters(['policyAreas', 'subjects']),

    policyAreaItems() {
      return this.policyAreas.map(p => p.policyArea)
    },

    subjectItems() {
      return this.subjects.map(s => s.subject)
    },
  },

  data: () => ({
    selectedSubjects: [],
    selectedPolicyAreas: [],
    loading: false,
    axiosSource: null,
  }),

  created() {
    this.dispatchGetSubjects()
  },

  methods: {
    ...mapActions({
      dispatchGetSubjects: 'getSubjects',
      dispatchGetCells: 'getCells',
    }),

    async getCells() {
      this.loading = true

      const params = {
        policyAreas: this.selectedPolicyAreas.join(','),
        subjects: this.selectedSubjects.join(','),
      }

      try {
        const { cells, source } = await this.dispatchGetCells({
          params,
          previousSource: this.axiosSource,
        })
        this.axiosSource = source
        this.cells = cells || []
      } finally {
        this.loading = false
      }
    },
  },
}
</script>

<style></style>
