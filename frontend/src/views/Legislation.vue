<template>
  <v-container>
    <v-row align="center">
      <v-col class="py-0">
        <div class="brand">Legislation</div>
      </v-col>
      <v-col cols="auto" class="py-0">
        <v-menu offset-y left>
          <template v-slot:activator="{ on }">
            <v-btn outlined color="primary" v-on="on">
              <v-icon>mdi-filter-variant</v-icon>
              Add Filter
            </v-btn>
          </template>
          <v-list dense>
            <v-list-item
              v-if="!bipartisanFilter"
              @click="applyBipartisanFilter()"
            >
              <v-chip color="deep-purple" dark>Bipartisan Sponsorship</v-chip>
            </v-list-item>
            <v-list-item @click="openFilterDialog('Policy Area')">
              <v-chip color="primary">By Policy Area</v-chip>
            </v-list-item>
            <v-list-item @click="openFilterDialog('Subject')">
              <v-chip dark>By Subject</v-chip>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-col>
    </v-row>

    <ActiveFilters
      :bipartisanFilter="bipartisanFilter"
      :policyAreaFilters="policyAreaFilters"
      :subjectFilters="subjectFilters"
      @clear="clearFilters()"
      @remove="removeFilter($event)"
    />

    <v-dialog v-model="showFilterDialog" max-width="500" persistent>
      <FilterDialog
        :filterType="filterType"
        :policyAreaFilters="policyAreaFilters"
        :subjectFilters="subjectFilters"
        @close="showFilterDialog = false"
        @apply="applyFilters($event)"
      />
    </v-dialog>

    <v-row>
      <v-col>
        <v-combobox
          v-model="bill"
          :items="bills"
          item-text="title"
          :loading="loading"
          :search-input.sync="query"
          solo
          :hide-no-data="hideNoData"
          placeholder="Start typing to search"
          prepend-inner-icon="mdi-magnify"
          clearable
          hint="Use * to see all results after any filters are applied."
          persistent-hint
          return-object
          no-filter
        >
          <template v-slot:no-data>
            <v-list class="py-0">
              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title>
                    No matches found for: {{ query }}
                  </v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </template>
        </v-combobox>
        <v-expand-transition>
          <Bill
            v-if="bill && typeof bill === 'object'"
            :bill="bill"
            :isBillView="true"
          />
        </v-expand-transition>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions } from 'vuex'

import FilterDialog from '@/components/FilterDialog'
import ActiveFilters from '@/components/ActiveFilters'
import Bill from '@/components/Bill'

export default {
  components: {
    FilterDialog,
    ActiveFilters,
    Bill,
  },

  computed: {
    hideNoData() {
      return (
        this.loading ||
        !this.query ||
        this.query.length < 3 ||
        (this.bill !== null && this.bill.title === this.query)
      )
    },
  },

  data: () => ({
    bills: [],
    loading: false,
    bill: null,
    query: '',
    bipartisanFilter: false,
    policyAreaFilters: [],
    subjectFilters: [],
    showFilterDialog: false,
    filterType: '',
    timeout: null,
    axiosSource: null,
  }),

  created() {
    this.dispatchGetSubjects()
  },

  beforeDestroy() {
    clearTimeout(this.timeout)
  },

  methods: {
    ...mapActions({
      dispatchGetBillsByTitle: 'getBillsByTitle',
      dispatchGetSubjects: 'getSubjects',
    }),

    clearFilters() {
      this.bipartisanFilter = false
      this.policyAreaFilters = []
      this.subjectFilters = []
    },

    applyBipartisanFilter() {
      this.bipartisanFilter = true
      this.bill = null
      this.bills = []
      this.showFilterDialog = false
    },

    openFilterDialog(filterType) {
      this.filterType = filterType
      this.showFilterDialog = true
    },

    applyFilters(filters) {
      if (this.filterType === 'Policy Area') {
        this.policyAreaFilters.push(...filters)
      } else {
        this.subjectFilters.push(...filters)
      }
      this.bill = null
      this.bills = []
      this.showFilterDialog = false
    },

    removeFilter({ type, name }) {
      if (type === 'bipartisan') {
        this.bipartisanFilter = false
      } else if (type === 'policyArea') {
        this.policyAreaFilters = this.policyAreaFilters.filter(
          f => f.policyArea !== name,
        )
      } else if (type === 'subject') {
        this.subjectFilters = this.subjectFilters.filter(
          f => f.subject !== name,
        )
      }
    },

    search() {
      this.loading = true

      clearTimeout(this.timeout)

      const billNumbers = []
      this.policyAreaFilters.forEach(f => {
        billNumbers.push(...f.billNumbers)
      })
      this.subjectFilters.forEach(f => {
        billNumbers.push(...f.billNumbers)
      })

      const params = {
        query: this.query,
        bipartisan: this.bipartisanFilter,
        billNumbers: billNumbers.join(','),
      }

      this.timeout = setTimeout(async () => {
        try {
          const { bills, source } = await this.dispatchGetBillsByTitle({
            params,
            previousSource: this.axiosSource,
          })
          this.axiosSource = source
          this.bills = bills || []
        } finally {
          this.loading = false
        }
      }, 2000)
    },
  },

  watch: {
    query(query) {
      if (
        query &&
        (query.length >= 3 || query === '*') &&
        (!this.bill || this.bill.title !== query)
      ) {
        this.search(query)
      } else {
        if (this.axiosSource) this.axiosSource.cancel()
        this.bills = []
      }
    },
  },
}
</script>
