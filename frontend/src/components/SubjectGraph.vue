<template>
  <v-card flat>
    <v-card-text>
      <h3 class="my-2">Step 1: Choose relevant subjects</h3>
      <v-row>
        <v-col>
          <v-autocomplete
            label="Subject"
            v-model="selectedSubjects"
            :items="subjectItems"
            outlined
            hide-details
            multiple
            small-chips
            deletable-chips
            clearable
            @change="getBills()"
          />
        </v-col>
      </v-row>
      <template v-if="selectedSubjects.length > 0">
        <h3 class="my-2">Step 2: Choose relevant legislation</h3>
        <v-row>
          <v-col class="border ma-3">
            <div
              v-for="bill in concatenatedBills"
              :key="`bill-chip-${bill.number}`"
            >
              <v-tooltip v-if="bill.title.length > 75" bottom>
                <template v-slot:activator="{ on }">
                  <v-chip
                    class="ma-1"
                    v-on="on"
                    :color="bill.selected ? 'primary' : ''"
                    close-icon="mdi-open-in-new"
                    close
                    @click="toggleBill(bill)"
                    @click:close="exploreBill(bill)"
                  >
                    {{ bill.truncatedTitle }} ({{ bill.sponsorCount }})
                  </v-chip>
                </template>
                {{ bill.title }}
              </v-tooltip>
              <v-chip
                v-else
                class="ma-1"
                :color="bill.selected ? 'primary' : ''"
                close-icon="mdi-open-in-new"
                close
                @click="toggleBill(bill)"
                @click:close="exploreBill(bill)"
              >
                {{ bill.truncatedTitle }} ({{ bill.sponsorCount }})
              </v-chip>
            </div>
          </v-col>
        </v-row>
      </template>
      <div v-if="selectedBillNumbers.length > 0" class="view-box">
        <div class="subject-graph-svg-container" />
      </div>
    </v-card-text>
  </v-card>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import { truncate } from '@/common/functions'
import { drawGraph } from '@/d3/graph'

export default {
  computed: {
    ...mapGetters(['subjects']),

    subjectItems() {
      return this.subjects.map(s => s.subject)
    },

    truncatedBills() {
      return this.bills.map(b => ({
        ...b,
        truncatedTitle: truncate(b.title, 75),
        sponsorCount: b.numDems + b.numReps + b.numLibs + b.numInds,
        selected: this.selectedBillNumbers.includes(b.number),
      }))
    },

    selectedBills() {
      return this.truncatedBills.filter(b => b.selected)
    },

    unselectedBills() {
      return this.truncatedBills.filter(b => !b.selected)
    },

    concatenatedBills() {
      return this.selectedBills.concat(this.unselectedBills)
    },
  },

  data: () => ({
    selectedSubjects: [],
    axiosSource: null,
    bills: [],
    billsLoading: false,
    selectedBillNumbers: [],
    removeFn: () => {},
  }),

  created() {
    this.dispatchGetSubjects()
  },

  methods: {
    ...mapActions({
      dispatchGetSubjects: 'getSubjects',
      dispatchGetBillsBySubject: 'getBillsBySubject',
    }),

    exploreBill(bill) {
      window.open(bill.link, '_blank')
    },

    toggleBill(bill) {
      if (this.selectedBillNumbers.includes(bill.number)) {
        this.selectedBillNumbers = this.selectedBillNumbers.filter(
          n => n !== bill.number,
        )
      } else {
        this.selectedBillNumbers.push(bill.number)
      }

      // trigger CD
      this.bills = this.bills.slice()

      // selecting or unselecting a bill should result in a new graph
      this.redrawGraph()
    },

    generateGraph() {
      const nodes = []
      const links = []
      const memberToGraphID = {}
      let graphID = 1
      this.selectedBills.forEach(b => {
        const billGraphID = graphID++
        nodes.push({
          id: billGraphID,
          type: 'bill',
          display: b.title,
        })
        b.sponsors.concat(b.cosponsors).forEach(m => {
          const partyChar = m.split('[')[1][0]
          let type
          if (partyChar === 'D') {
            type = 'democrat'
          } else if (partyChar === 'R') {
            type = 'republican'
          } else {
            type = 'other'
          }
          if (memberToGraphID[m]) {
            nodes.push({
              id: memberToGraphID[m],
              type,
              display: m,
            })
            links.push({
              source: billGraphID,
              target: memberToGraphID[m],
            })
          } else {
            nodes.push({
              id: graphID++,
              type,
              display: m,
            })
            memberToGraphID[m] = graphID - 1
            links.push({
              source: billGraphID,
              target: graphID - 1,
            })
          }
        })
      })
      return { nodes, links }
    },

    async redrawGraph() {
      this.removeFn()
      await this.$nextTick()
      this.removeFn = drawGraph(
        '.subject-graph-svg-container',
        this.generateGraph(),
        {
          onMouseover: this.onMouseover,
          onMouseleave: this.onMouseleave,
          onMouseup: this.onMouseup,
        },
      )
    },

    async getBills() {
      if (this.selectedSubjects.length === 0) {
        return
      }

      this.billsLoading = true
      const params = {
        subjects: this.selectedSubjects.join(','),
        bipartisan: true,
      }

      try {
        const { bills, source } = await this.dispatchGetBillsBySubject({
          params,
          previousSource: this.axiosSource,
        })
        this.axiosSource = source
        this.bills = bills || []
      } finally {
        this.billsLoading = false
      }
    },

    clear() {
      this.bills = []
      this.selectedBillNumbers = []
      this.removeFn()
    },
  },

  watch: {
    selectedSubjects(state) {
      if (!state.length) {
        this.clear()
      }
    },
  },
}
</script>

<style scoped>
.border {
  border: 1px solid rgba(0, 0, 0, 0.38);
  border-radius: 4px;
}
</style>

<style>
.subject-graph-svg-container {
  display: inline-block;
  position: relative;
  width: 100%;
  vertical-align: top;
  overflow: hidden;
}

.svg-content {
  display: inline-block;
  position: absolute;
  top: 0;
  left: 0;
}

.bill-node {
  fill: gray;
}

.dem-node {
  fill: #1565c0;
}

.rep-node {
  fill: #f44336;
}

.oth-node {
  fill: green;
}
</style>
