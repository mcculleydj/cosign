<template>
  <v-container>
    <v-row>
      <v-col class="large-text">
        <p>
          This network graph visualization considers a data access pattern
          beginning with a particular set of subjects. The user then further
          specifies what kinds of legislation they are interested in by
          selecting individual bills from the list of all bills that match any
          of the selected subjects (set union). Because the focus of this
          application is on bipartisan support for bills, every bill listed will
          have at least one sponsor from each party.
        </p>
        <p>
          Once two or more bills are selected we can generate a useful graph
          using D3's force layout to show which members have sponsored these
          bills. For a single bill, this information should be retrieved in the
          legislation view where every sponsor for the selected bill is listed
          by name, but for multiple bills this visualization immediately makes
          it clear which members sponsored more of the selected bills than their
          peers allowing the user to identify a bipartisan cohort of members
          likely to support similar legislation.
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

    <div
      v-show="selectedBillNumbers.length > 1"
      id="subject-graph-svg-container"
    />
  </v-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import Graph from '@/d3/graph'
import { sleep } from '@/common/functions'

const bills = [
  {
    number: 195,
    title: 'Pay our Doctors Act of 2019',
    sponsors: ['Mullin, Markwayne [R-OK-2]'],
    cosponsors: [
      'McCollum, Betty [D-MN-4]',
      'Cole, Tom [R-OK-4]',
      'Simpson, Michael K. [R-ID-2]',
      'Ruiz, Raul [D-CA-36]',
      'Clarke, Yvette D. [D-NY-9]',
      'Horn, Kendra S. [D-OK-5]',
      'Bonamici, Suzanne [D-OR-1]',
      'Gianforte, Greg [R-MT-At Large]',
      "O'Halleran, Tom [D-AZ-1]",
      'Young, Don [R-AK-At Large]',
      'Torres, Norma J. [D-CA-35]',
      'Moore, Gwen [D-WI-4]',
      'Curtis, John R. [R-UT-3]',
      'Kilmer, Derek [D-WA-6]',
      'Armstrong, Kelly [R-ND-At Large]',
      'Haaland, Debra A. [D-NM-1]',
    ],
    score: 3,
    numDems: 10,
    numReps: 7,
    numInds: 0,
    numLibs: 0,
    multiParty: true,
    link: 'https://www.congress.gov/bill/116th-congress/house-bill/195',
    policyArea: 'Native Americans',
    subjects: ['Appropriations', 'Executive agency funding and structure'],
  },
  {
    number: 312,
    title: 'Mashpee Wampanoag Tribe Reservation Reaffirmation Act',
    sponsors: ['Keating, William R. [D-MA-9]'],
    cosponsors: [
      'Kennedy, Joseph P., III [D-MA-4]',
      'Young, Don [R-AK-At Large]',
      'Grijalva, Raul M. [D-AZ-3]',
      'LaMalfa, Doug [R-CA-1]',
      'Gallego, Ruben [D-AZ-7]',
      'McClintock, Tom [R-CA-4]',
      'Cole, Tom [R-OK-4]',
      'Fitzpatrick, Brian K. [R-PA-1]',
      'Neal, Richard E. [D-MA-1]',
      'McGovern, James P. [D-MA-2]',
      'Clark, Katherine M. [D-MA-5]',
      'Trahan, Lori [D-MA-3]',
      'Pingree, Chellie [D-ME-1]',
      'Davids, Sharice [D-KS-3]',
      'Haaland, Debra A. [D-NM-1]',
      'Mullin, Markwayne [R-OK-2]',
      "O'Halleran, Tom [D-AZ-1]",
      'Pressley, Ayanna [D-MA-7]',
      'Raskin, Jamie [D-MD-8]',
      'Moulton, Seth [D-MA-6]',
      'McCollum, Betty [D-MN-4]',
      'Van Drew, Jefferson [D-NJ-2]',
      'Lynch, Stephen F. [D-MA-8]',
      'Scanlon, Mary Gay [D-PA-5]',
      'Frankel, Lois [D-FL-21]',
      'Suozzi, Thomas R. [D-NY-3]',
      'Kind, Ron [D-WI-3]',
      'Vargas, Juan [D-CA-51]',
      'Dean, Madeleine [D-PA-4]',
      'Jackson Lee, Sheila [D-TX-18]',
      'Horn, Kendra S. [D-OK-5]',
      'Panetta, Jimmy [D-CA-20]',
      'Pascrell, Bill, Jr. [D-NJ-9]',
      'Doyle, Michael F. [D-PA-18]',
      'DeSaulnier, Mark [D-CA-11]',
    ],
    score: 24,
    numDems: 30,
    numReps: 6,
    numInds: 0,
    numLibs: 0,
    multiParty: true,
    link: 'https://www.congress.gov/bill/116th-congress/house-bill/312',
    policyArea: 'Native Americans',
    subjects: [
      'Federal-Indian relations',
      'Indian lands and resources rights',
      'Massachusetts',
    ],
  },
  {
    number: 317,
    title: 'Santa Ynez Band of Chumash Indians Land Affirmation Act of 2019',
    sponsors: ['LaMalfa, Doug [R-CA-1]'],
    cosponsors: ['Carbajal, Salud O. [D-CA-24]'],
    score: 0,
    numDems: 1,
    numReps: 1,
    numInds: 0,
    numLibs: 0,
    multiParty: true,
    link: 'https://www.congress.gov/bill/116th-congress/house-bill/317',
    policyArea: 'Native Americans',
    subjects: [
      'California',
      'Federal-Indian relations',
      'Gambling',
      'Indian lands and resources rights',
      'Land transfers',
    ],
  },
  {
    number: 375,
    title:
      'To amend the Act of June 18, 1934, to reaffirm the authority of the Secretary of the Interior to take land into trust for Indian Tribes, and for other purposes.',
    sponsors: ['Cole, Tom [R-OK-4]'],
    cosponsors: [
      'McCollum, Betty [D-MN-4]',
      'Moore, Gwen [D-WI-4]',
      'Simpson, Michael K. [R-ID-2]',
      'Kilmer, Derek [D-WA-6]',
      'Calvert, Ken [R-CA-42]',
      'Larsen, Rick [D-WA-2]',
      'Cardenas, Tony [D-CA-29]',
      'Gallego, Ruben [D-AZ-7]',
      'Keating, William R. [D-MA-9]',
      'Schakowsky, Janice D. [D-IL-9]',
      'Haaland, Debra A. [D-NM-1]',
      "O'Halleran, Tom [D-AZ-1]",
      'Gomez, Jimmy [D-CA-34]',
      'Lieu, Ted [D-CA-33]',
      'Pingree, Chellie [D-ME-1]',
      'Heck, Denny [D-WA-10]',
      'Joyce, David P. [R-OH-14]',
      'Cook, Paul [R-CA-8]',
      'Pocan, Mark [D-WI-2]',
      'Ruiz, Raul [D-CA-36]',
      'Torres Small, Xochitl [D-NM-2]',
      'DelBene, Suzan K. [D-WA-1]',
      'Torres, Norma J. [D-CA-35]',
      'Lujan, Ben Ray [D-NM-3]',
      'Aguilar, Pete [D-CA-31]',
      'Kildee, Daniel T. [D-MI-5]',
      'Mullin, Markwayne [R-OK-2]',
      'Davids, Sharice [D-KS-3]',
    ],
    score: 17,
    numDems: 23,
    numReps: 6,
    numInds: 0,
    numLibs: 0,
    multiParty: true,
    link: 'https://www.congress.gov/bill/116th-congress/house-bill/375',
    policyArea: 'Native Americans',
    subjects: [
      'Federal-Indian relations',
      'Indian lands and resources rights',
      'Land transfers',
    ],
  },
]

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
    selectedBillNumbers: [195, 312],
    minimumBillCount: 2,
  }),

  // TODO: handle resize events using RxJS

  created() {
    this.dispatchGetSubjects()
  },

  mounted() {
    this.graph = new Graph('subject-graph-svg-container')
    this.bills = bills
    setTimeout(() => {
      this.updateGraph()
    }, 500)
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

    onChange() {
      // check for changes to the selected bills
      if (
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
</style>

<style>
#subject-graph-svg-container {
  display: inline-block;
  position: relative;
  width: 100%;
  vertical-align: top;
  overflow: hidden;
  height: 800px;
}
</style>
