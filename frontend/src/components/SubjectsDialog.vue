<template>
  <v-card>
    <v-card-title class="d-flex justify-space-between">
      Policy Area: {{ bill.policyArea || 'None' }}
      <v-btn icon @click="$emit('close')">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-text>
      <v-row>
        <v-col
          v-for="subject in bill.subjects"
          :key="`subject-dialog-${bill.number}-${subject}`"
          cols="auto"
          class="pa-1"
        >
          <v-chip>{{ subject }}</v-chip>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  props: ['bill'],
}
</script>
