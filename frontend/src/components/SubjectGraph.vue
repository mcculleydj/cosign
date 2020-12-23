<template>
  <v-container>
    <v-row>
      <v-col class="large-text">
        <p>
          This network graph visualization considers a data access pattern
          beginning with a particular subject or policy area. The user then
          further specifies what kinds of legislation they are interested in by
          selecting individual bills from the list of all corresponding bills
          (set union). Because the focus of this application is on bipartisan
          support for bills, every bill listed will have at least one sponsor
          from each party.
        </p>
        <p>
          Once two or more bills are selected, D3's force layout helps identify
          a bipartisan cohort of members likely to support similar legislation
          by clustering members between the bills they've sponsored. By default,
          I hide all member nodes that have not sponsored at least two of the of
          selected bills to reduce noise. That information may still be useful
          or the user may want to tighten the bounds further leading to the
          inclusion of a select control to set the minimum number of bills a
          member must have sponsored to remain in the graph.
        </p>
      </v-col>
    </v-row>

    <h3 class="my-2">Step 1: Choose relevant policy areas and subjects</h3>

    <v-row>
      <v-col>
        <v-autocomplete
          label="Policy Areas"
          v-model="policyAreaBillNumbers"
          :items="policyAreaItems"
          outlined
          multiple
          small-chips
          deletable-chips
          clearable
        />
        <v-autocomplete
          label="Subjects"
          v-model="subjectBillNumbers"
          :items="subjectItems"
          outlined
          multiple
          hide-details
          small-chips
          deletable-chips
          clearable
        />
      </v-col>
    </v-row>

    <h3 class="my-2">
      Step 2: Choose at least two pieces of relevant legislation
    </h3>

    <v-autocomplete
      v-model="bills"
      :items="billOptions"
      item-text="title"
      :loading="loading"
      :search-input.sync="query"
      outlined
      :hide-no-data="hideNoData"
      placeholder="Start typing to search"
      prepend-inner-icon="mdi-magnify"
      clearable
      small-chips
      deletable-chips
      hint="Use * to see all results matching policy area(s) and subject(s)."
      persistent-hint
      return-object
      no-filter
      multiple
      @change="onChange()"
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
    </v-autocomplete>

    <v-row v-if="selectedBillNumbers.length > 1">
      <v-col class="pb-0">
        <v-select
          v-model="minimumBillCount"
          :items="minimumBillCountItems"
          hide-details
          outlined
          label="Minimum Bills"
          @change="onChange(true)"
          style="max-width: 140px"
        />
      </v-col>
    </v-row>

    <div
      v-show="selectedBillNumbers.length > 1"
      id="subject-graph-svg-container"
      class="border mt-5"
    />
  </v-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import Graph from '@/d3/graph'
import { sleep } from '@/common/functions'

// data transformation
function generateGraph(selectedBills, minimumBillCount) {
  // first determine how many bills each member has sponsored
  // from the subset of bills selected
  const memberToBillCount = {}
  selectedBills.forEach(b => {
    b.sponsors.concat(b.cosponsors).forEach(m => {
      if (memberToBillCount[m]) {
        memberToBillCount[m]++
      } else {
        memberToBillCount[m] = 1
      }
    })
  })

  // now populate nodes and links
  // filtering out any members not meeting minimum bill count threshold
  const nodes = []
  const links = []
  const memberToNodeID = {}
  const nodeIDToNode = {}
  let nodeIDSequence = 1

  selectedBills.forEach(b => {
    const billNodeID = nodeIDSequence++

    // assign the next nodeID and increment sequence
    nodes.push({
      id: billNodeID,
      type: 'bill',
      display: b.title,
    })

    // iterate over all sponsors meeting the threshold
    b.sponsors
      .concat(b.cosponsors)
      .filter(m => memberToBillCount[m] >= minimumBillCount)
      .forEach(m => {
        // set type based on party
        const party = m.split('[')[1][0]
        let type
        if (party === 'D') {
          type = 'democrat'
        } else if (party === 'R') {
          type = 'republican'
        } else {
          type = 'other'
        }

        if (memberToNodeID[m]) {
          // if this member is already in the graph, simply add the link
          links.push({
            source: billNodeID,
            target: memberToNodeID[m],
          })
          nodeIDToNode[memberToNodeID[m]].billCount++
        } else {
          // add the member to the graph along with a link to this bill
          // and memoize the fact that this member is part of the graph now
          const memberNodeID = nodeIDSequence++
          memberToNodeID[m] = memberNodeID
          const node = {
            id: memberNodeID,
            type,
            display: m,
            billCount: 1,
          }
          nodeIDToNode[memberNodeID] = node
          nodes.push(node)
          links.push({
            source: billNodeID,
            target: memberNodeID,
          })
        }
      })
  })

  return { nodes, links }
}

// TODO: make it clear if this is the intersection between subjects
// or the union, maybe even allow a user to control this

export default {
  computed: {
    ...mapGetters(['policyAreas', 'subjects']),

    policyAreaItems() {
      return this.policyAreas.map(p => ({
        text: p.policyArea,
        value: p.billNumbers,
      }))
    },

    subjectItems() {
      return this.subjects.map(s => ({
        text: s.subject,
        value: s.billNumbers,
      }))
    },

    hideNoData() {
      return (
        this.loading ||
        !this.query ||
        this.query.length < 3 ||
        this.bills.length > 0
      )
    },

    minimumBillCountItems() {
      const items = []
      for (let i = 1; i <= this.selectedBillNumbers.length; i++) {
        items.push(i)
      }
      return items
    },
  },

  data: () => ({
    graph: null,
    policyAreaBillNumbers: [],
    subjectBillNumbers: [],
    query: '',
    loading: false,
    bills: [],
    billOptions: [],
    axiosSource: null,
    timeout: null,
    selectedBillNumbers: [],
    minimumBillCount: 2,
  }),

  // TODO: handle resize events using RxJS

  created() {
    this.dispatchGetSubjects()
  },

  mounted() {
    this.graph = new Graph('subject-graph-svg-container')
  },

  beforeDestroy() {
    clearTimeout(this.timeout)
  },

  methods: {
    ...mapActions({
      dispatchGetSubjects: 'getSubjects',
      dispatchGetBillsByTitle: 'getBillsByTitle',
    }),

    search() {
      this.loading = true
      clearTimeout(this.timeout)

      const selectedBillNumbers = this.bills.map(b => b.number)

      const billNumberSet = {}
      this.policyAreaBillNumbers
        .concat(this.subjectBillNumbers)
        .filter(n => !selectedBillNumbers.includes(n))
        .forEach(n => {
          billNumberSet[n] = true
        })

      const params = {
        query: this.query,
        bipartisan: true,
        billNumbers: Object.keys(billNumberSet).join(','),
      }

      this.timeout = setTimeout(async () => {
        try {
          const { bills, source } = await this.dispatchGetBillsByTitle({
            params,
            previousSource: this.axiosSource,
          })
          this.axiosSource = source
          this.billOptions = bills ? bills.concat(this.bills) : this.bills
        } finally {
          this.loading = false
        }
      }, 2000)
    },

    onChange(paramChange) {
      // check for changes to the selected bills
      if (
        paramChange ||
        this.bills.length !== this.selectedBillNumbers.length ||
        !this.selectedBillNumbers.every((n, i) => n === this.bills[i].number)
      ) {
        this.selectedBillNumbers = this.bills.map(b => b.number)

        // supply D3 with new data on change and if number of bills is 2 or more
        if (this.selectedBillNumbers.length > 1) {
          this.updateGraph()
        }
      }
    },

    async waitForDOM() {
      const container = document.getElementById('subject-graph-svg-container')
      while (!container || !container.clientWidth || !container.clientHeight) {
        await sleep(10)
      }
    },

    // TODO: use this
    exploreBill(bill) {
      window.open(bill.link, '_blank')
    },

    async updateGraph() {
      await this.waitForDOM()
      // TODO: add callbacks for mouse events
      this.graph.update(generateGraph(this.bills, this.minimumBillCount))
    },

    clear() {
      this.bills = []
      this.selectedBillNumbers = []

      // TODO: call update to remove nodes and edges
    },
  },

  watch: {
    query(query) {
      if (query && (query.length >= 3 || query === '*')) {
        this.search(query)
      } else {
        if (this.axiosSource) this.axiosSource.cancel()
        this.billOptions = this.bills
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

.large-text {
  font-size: 1.1rem;
}

#subject-graph-svg-container {
  display: inline-block;
  position: relative;
  width: 100%;
  vertical-align: top;
  overflow: hidden;
  height: 500px;
  cursor: grab;
}
</style>
