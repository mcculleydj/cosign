<template>
  <v-container>
    <div class="brand">Cosponsors</div>
    <v-row>
      <v-col>
        <v-autocomplete
          :items="sponsorItems"
          label="Sponser"
          spellcheck="false"
          clearable
          v-model="sponsor"
          solo
          hide-details
        />
      </v-col>
      <v-col>
        <v-autocomplete
          :items="cosponsorItems"
          label="Cosponsor"
          spellcheck="false"
          clearable
          v-model="cosponsor"
          solo
          hide-details
        />
      </v-col>
    </v-row>
    <template v-if="cell">
      <v-divider class="mt-3 mb-2" />
      <v-row v-if="cell.bills.length">
        <v-col>
          <p>Cosponsored Legislation</p>
          <Bill
            v-for="bill in cell.bills"
            :key="`bill-number-${bill.number}`"
            :bill="bill"
          />
        </v-col>
      </v-row>
      <v-row v-else-if="cell.bills.length === 0">
        <v-divider class="mt-3 mb-1" />
        <v-col>
          <i>No cosponsored legislation.</i>
        </v-col>
      </v-row>
    </template>
    <template v-else-if="loading">
      <v-divider class="mt-3 mb-2" />
      <v-row v-for="i in 5" :key="`loader-${i}`">
        <v-col>
          <v-skeleton-loader type="card-heading, list-item-three-line" />
        </v-col>
      </v-row>
    </template>
  </v-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import Bill from '@/components/Bill.vue'

export default {
  components: {
    Bill,
  },

  computed: {
    ...mapGetters(['memberItems', 'cell']),

    multiPartyMembers() {
      return this.memberItems.filter(m => m.value.parties.length > 1)
    },

    singlePartyMembers() {
      return this.memberItems.filter(m => m.value.parties.length === 1)
    },

    sponsorItems() {
      let items

      if (!this.cosponsor) {
        items = this.memberItems
      } else if (
        this.multiPartyMembers.map(m => m.value).includes(this.cosponsor)
      ) {
        items = this.memberItems.filter(m => m.value.id !== this.cosponsor.id)
      } else {
        items = this.singlePartyMembers
          .filter(
            m =>
              m.value.id !== this.cosponsor.id &&
              m.value.parties[0] !== this.cosponsor.parties[0],
          )
          .concat(this.multiPartyMembers)
      }

      items.sort((m1, m2) => {
        if (m1.text < m2.text) {
          return -1
        }
        return 1
      })

      return items
    },

    cosponsorItems() {
      let items

      if (!this.sponsor) {
        items = this.memberItems
      } else if (
        this.multiPartyMembers.map(m => m.value).includes(this.sponsor)
      ) {
        items = this.memberItems.filter(m => m.value.id !== this.sponsor.id)
      } else {
        items = this.singlePartyMembers
          .filter(
            m =>
              m.value.id !== this.sponsor.id &&
              m.value.parties[0] !== this.sponsor.parties[0],
          )
          .concat(this.multiPartyMembers)
      }

      items.sort((m1, m2) => {
        if (m1.text < m2.text) {
          return -1
        }
        return 1
      })

      return items
    },
  },

  data: () => ({
    sponsor: null,
    cosponsor: null,
    loading: false,
  }),

  created() {
    this.dispatchGetMembers()
  },

  beforeDestroy() {
    this.dispatchSetCell(null)
  },

  methods: {
    ...mapActions({
      dispatchGetMembers: 'getMembers',
      dispatchGetCell: 'getCell',
      dispatchSetCell: 'setCell',
    }),
  },

  watch: {
    async sponsor(sponsor) {
      if (sponsor && this.cosponsor) {
        let position = `${sponsor.id}_${this.cosponsor.id}`
        if (sponsor.id > this.cosponsor.id) {
          position = `${this.cosponsor.id}_${sponsor.id}`
        }
        this.loading = true
        this.dispatchGetCell(position).finally(() => {
          this.loading = false
        })
      } else {
        this.dispatchSetCell(null)
      }
    },

    cosponsor(cosponsor) {
      if (cosponsor && this.sponsor) {
        let position = `${cosponsor.id}_${this.sponsor.id}`
        if (cosponsor.id > this.sponsor.id) {
          position = `${this.sponsor.id}_${cosponsor.id}`
        }
        this.loading = true
        this.dispatchGetCell(position).finally(() => {
          this.loading = false
        })
      } else {
        this.dispatchSetCell(null)
      }
    },
  },
}
</script>
