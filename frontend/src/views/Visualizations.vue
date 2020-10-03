<template>
  <v-container>
    <div class="brand">Visualizations</div>
    <v-row>
      <v-col>
        <v-autocomplete
          label="Member"
          v-model="member"
          :items="memberItems"
          solo
        />
      </v-col>
    </v-row>
    <template v-if="member">
      <v-row align="center">
        <v-col cols="auto">
          <v-select
            label="Threshold"
            :items="thresholds"
            v-model="threshold"
            hide-details
            outlined
            dense
          />
        </v-col>
        <v-col cols="auto">
          <v-switch label="Sticky" v-model="sticky" />
        </v-col>
        <v-col>
          <v-slider
            label="Spread"
            hide-details
            v-model="strength"
            min="30"
            max="200"
            inverse-label
            @end="draw()"
          />
        </v-col>
        <v-spacer />
        <v-col cols="auto">
          <v-btn outlined color="primary" @click="resetFn">
            <v-icon>mdi-refresh</v-icon>
            reset
          </v-btn>
        </v-col>
      </v-row>
      <svg />
    </template>
  </v-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import { drawGraph } from '@/d3/graph'

export default {
  computed: {
    ...mapGetters({
      memberItems: 'memberItems',
      memberMap: 'memberMap',
    }),
  },

  data: () => ({
    sticky: false,
    resetFn: () => {},
    removeFn: () => {},
    member: null,
    thresholds: [0, 10, 20, 30, 40, 50],
    threshold: 40,
    strength: 30,
  }),

  created() {
    if (!this.memberItems.length) {
      this.dispatchGetMembers()
    }
  },

  methods: {
    ...mapActions({
      dispatchGetMembers: 'getMembers',
    }),

    async draw() {
      this.removeFn()
      await this.$nextTick()
      const { reset, remove } = drawGraph('svg', this.member, this.memberMap, {
        threshold: this.threshold,
        sticky: this.sticky,
        strength: -this.strength,
      })
      this.resetFn = reset
      this.removeFn = remove
    },
  },

  watch: {
    member() {
      this.draw()
    },

    threshold() {
      this.draw()
    },

    sticky(state) {
      if (!state) {
        this.draw()
      } else {
        this.draw(true)
      }
    },
  },
}
</script>

<style scoped>
svg {
  border: 1px solid black;
  border-radius: 5px;
  height: 500px;
  width: 100%;
}
</style>

<style>
.link {
  stroke: #333;
  stroke-width: 1px;
}

.node {
  cursor: move;
  opacity: 1;
}
.node.rep {
  fill: #f44336;
}
.node.dem {
  fill: #1565c0;
}
.node.oth {
  fill: #4caf50;
}

.node.fixed {
  opacity: 1;
  stroke-width: 1.5px;
  stroke: #000;
}
</style>
