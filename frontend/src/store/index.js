import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

import { memberToString } from '@/common/functions'

const api = axios.create({
  baseURL: `${window.location.origin}/api/`,
})

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    members: [],
    memberMap: {},
    cell: null,
    policyAreas: [],
    subjects: [],
  },

  getters: {
    memberMap: state => state.memberMap,
    memberItems: state =>
      state.members.map(m => ({
        text: memberToString(m),
        value: m,
      })),
    cell: state => state.cell,
    policyAreas: state => state.policyAreas,
    subjects: state => state.subjects,
  },

  mutations: {
    SET_MEMBERS(state, { members, memberMap }) {
      state.members = members
      state.memberMap = memberMap
    },

    SET_CELL(state, cell) {
      state.cell = cell
    },

    SET_POLICY_AREAS(state, policyAreas) {
      state.policyAreas = policyAreas
    },

    SET_SUBJECTS(state, subjects) {
      state.subjects = subjects
    },
  },

  actions: {
    async getMembers({ commit }) {
      try {
        const response = await api.get('members')
        commit('SET_MEMBERS', response.data)
      } catch (err) {
        console.error(err)
      }
    },

    async getCell({ commit }, position) {
      try {
        const response = await api.get(`cell/${position}`)
        const cell = response.data
        if (cell.bills.length) {
          cell.bills.sort((b1, b2) => b1.number - b2.number)
        }
        commit('SET_CELL', cell)
      } catch (err) {
        console.error(err)
      }
    },

    setCell({ commit }, cell) {
      commit('SET_CELL', cell)
    },

    async getBillsByNumbers(_, billNumbers) {
      try {
        const response = await api.get('bills/number', {
          params: { billNumbers },
        })
        return response.data
      } catch (err) {
        console.error(err)
      }
    },

    async getBillsByTitle(_, { params, previousSource }) {
      if (previousSource) {
        previousSource.cancel()
      }

      const source = axios.CancelToken.source()

      try {
        const response = await api.get('bills/title', {
          params,
          cancelToken: source.token,
        })
        return { bills: response.data, source }
      } catch (err) {
        console.error(err)
      }
    },

    async getBillsBySubject(_, { params, previousSource }) {
      if (previousSource) {
        previousSource.cancel()
      }

      const source = axios.CancelToken.source()

      try {
        const response = await api.get('bills/subject', {
          params,
          cancelToken: source.token,
        })
        return { bills: response.data, source }
      } catch (err) {
        console.error(err)
      }
    },

    async getSubjects({ commit }) {
      try {
        const response = await api.get('subjects')
        commit('SET_POLICY_AREAS', response.data.policyAreas)
        commit('SET_SUBJECTS', response.data.subjects)
      } catch (err) {
        console.error(err)
      }
    },

    async getCells(_, { params, previousSource }) {
      if (previousSource) {
        previousSource.cancel()
      }
      const source = axios.CancelToken.source()
      try {
        const response = await api.get('cells', {
          params,
          cancelToken: source.token,
        })
        console.log('r', response)
        return { cells: response.data, source }
      } catch (err) {
        console.error(err)
      }
    },
  },
})
