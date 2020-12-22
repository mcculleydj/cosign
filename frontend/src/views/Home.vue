<template>
  <v-container>
    <div class="brand">Project Description</div>
    <div class="large-text mt-3 mb-6">
      <p>
        I built this to learn Vue and Go in the context of a complete web
        application serving data through an API. In my experience, the best way
        to learn the practical aspects of a new language or framework is to
        build something from scratch that includes, but extends beyond the
        borders of the thing you are learning so that you confront problems at
        the interface. For example, how to build for production using the Vue
        CLI or how to target a different operating system during compilation
        using Go.
      </p>
      <p>
        The first step was to pick a reasonably large publicly available dataset
        that was interesting enough to build an application around. Living a
        mile from the Capitol Building in DC inspired me to use congressional
        data. An XML describing every bill introduced in the House is available
        <a href="https://www.govinfo.gov/bulkdata" target="_blank">here</a>. The
        focus of this application is on which members have co-sponsored
        legislation together, especially across party lines. Combining
        sponsorship data with the subject and policy areas the legislation
        covers gives the user an idea of who is willing to work together despite
        party affiliation on particular issues.
      </p>
    </div>
    <v-card>
      <v-img src="stack.png"></v-img>
    </v-card>
  </v-container>
</template>

<script>
export default {}
</script>

<style scoped>
.large-text {
  font-size: 1.1rem;
}
</style>
