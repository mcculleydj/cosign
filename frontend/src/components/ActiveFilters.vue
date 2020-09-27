<template>
  <v-card v-if="hasFilters" class="mb-6">
    <v-card-text class="pt-1">
      <v-row align="center">
        <v-col class="py-0">
          Active Filters (applies to the list of results in the drop down when
          searching)
        </v-col>
        <v-col cols="auto" class="py-0 pr-0">
          <v-btn color="error" small text @click="$emit('clear')">
            Clear Filters
          </v-btn>
        </v-col>
      </v-row>
      <v-chip
        v-if="bipartisanFilter"
        color="deep-purple"
        dark
        close
        class="mr-3 mt-3"
        @click:close="$emit('remove', { type: 'bipartisan' })"
      >
        Bipartisan Sponsorship
      </v-chip>
      <v-chip
        v-for="filter in policyAreaFilters"
        :key="`policy-filter-${filter.policyArea}`"
        color="primary"
        dark
        close
        class="mr-3 mt-3"
        @click:close="
          $emit('remove', { type: 'policyArea', name: filter.policyArea })
        "
      >
        Policy Area: {{ filter.policyArea }}
      </v-chip>
      <v-chip
        v-for="filter in subjectFilters"
        :key="`subject-filter-${filter.subject}`"
        dark
        close
        class="mr-3 mt-3"
        @click:close="
          $emit('remove', { type: 'subject', name: filter.subject })
        "
      >
        Subject: {{ filter.subject }}
      </v-chip>
    </v-card-text>
  </v-card>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  props: ['bipartisanFilter', 'policyAreaFilters', 'subjectFilters'],

  computed: {
    ...mapGetters({
      policyAreas: 'policyAreas',
      subjects: 'subjects',
    }),

    hasFilters() {
      return (
        this.bipartisanFilter ||
        this.policyAreaFilters.length > 0 ||
        this.subjectFilters.length > 0
      )
    },
  },
}
</script>
