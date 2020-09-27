<template>
  <v-card>
    <v-card-title>Add {{ filterType }} Filter(s)</v-card-title>
    <v-card-text>
      <v-autocomplete
        v-model="filters"
        :items="filterType === 'Policy Area' ? policyAreaItems : subjectItems"
        chips
        multiple
        outlined
      />
    </v-card-text>
    <v-card-actions class="justify-end">
      <v-btn text @click="close()">Cancel</v-btn>
      <v-btn outlined color="primary" @click="apply()">Apply</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  props: ['filterType', 'policyAreaFilters', 'subjectFilters'],

  computed: {
    ...mapGetters({
      policyAreas: 'policyAreas',
      subjects: 'subjects',
    }),

    // filter out selected filters

    policyAreaItems() {
      return this.policyAreas
        .map(p => ({
          text: p.policyArea,
          value: p,
        }))
        .filter(f => !this.policyAreaFilters.includes(f.value))
        .sort((f1, f2) => {
          if (f1.text > f2.text) return 1
          return -1
        })
    },

    subjectItems() {
      return this.subjects
        .map(s => ({
          text: s.subject,
          value: s,
        }))
        .filter(f => !this.subjectFilters.includes(f.value))
        .sort((f1, f2) => {
          if (f1.text > f2.text) return 1
          return -1
        })
    },
  },

  data: () => ({
    filters: [],
  }),

  methods: {
    close() {
      this.$emit('close')
      this.filters = []
    },

    apply() {
      this.$emit('apply', this.filters)
      this.filters = []
    },
  },
}
</script>
