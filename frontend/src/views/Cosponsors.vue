<template>
  <v-container>
    <div class="brand">Cosponsors</div>
    <div class="large-text mt-3">
      <p>
        When you select a member as the sponsor, the list of available
        cosponsors adjusts to include only members affiliated with another
        party. When you select a member as the cosponsor, all of the resolutions
        that these two members have cosponsored together are listed below along
        with their policy area, subjects, and sponsorship breakdown along party
        lines in a D3 visualization meant to remind the user of the House floor.
        Selection of a particular member is made easier through the use of
        autocomplete controls allowing a user to type a partial name and limit
        the number of options to select from. Finally, the selected members are
        reflected in the URL state making this view refresh tolerant, link
        sharable, and giving the rest of the application the ability to route to
        this view with any two members already selected simply by setting URL
        query params.
      </p>
    </div>
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
          <p>
            Cosponsored Legislation ({{
              cell.bills.length === 1 ? '1 Bill' : `${cell.bills.length} Bills`
            }})
          </p>
          <Bill
            v-for="bill in cell.bills"
            :key="`bill-number-${bill.number}`"
            :bill="bill"
          />
        </v-col>
      </v-row>
      <v-row v-else-if="cell.bills.length === 0">
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
    sponsorInitialized: true,
    cosponsorInitialized: true,
  }),

  async created() {
    await this.dispatchGetMembers()

    if (this.$route.query.sponsorId) {
      this.sponsor = this.sponsorItems.find(
        s => s.value.id === +this.$route.query.sponsorId,
      ).value
      this.sponsorInitialized = false
    }

    if (this.$route.query.cosponsorId) {
      this.cosponsor = this.cosponsorItems.find(
        c => c.value.id === +this.$route.query.cosponsorId,
      ).value
      this.cosponsorInitialized = false
    }
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
    sponsor(sponsor) {
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

      if (this.sponsorInitialized) {
        if (sponsor) {
          this.$router.replace({
            name: 'cosponsors',
            query: {
              ...this.$route.query,
              sponsorId: sponsor.id,
            },
          })
        } else {
          this.$router.replace({
            name: 'cosponsors',
            query: {
              ...this.$route.query,
              sponsorId: undefined,
            },
          })
        }
      }

      this.sponsorInitialized = true
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

      if (this.cosponsorInitialized) {
        if (cosponsor) {
          this.$router.replace({
            name: 'cosponsors',
            query: {
              ...this.$route.query,
              cosponsorId: cosponsor.id,
            },
          })
        } else {
          this.$router.replace({
            name: 'cosponsors',
            query: {
              ...this.$route.query,
              cosponsorId: undefined,
            },
          })
        }
      }

      this.cosponsorInitialized = true
    },
  },
}
</script>

<style scoped>
.large-text {
  font-size: 1.1rem;
}
</style>
