<template>
  <v-container>
    <div class="brand">Project Description</div>
    <p>
      When I was first learning Vue, I wanted to build an application from
      scratch that relied on an API and persistance layer to search and serve
      data in response to frontend user actions. The first step was to pick a
      reasonably large publicly available dataset that was interesting enough to
      build an application around. Living about a mile from the Capitol in D.C.
      inspired me to use Congress.
    </p>
    <p>
      An XML summary of every bill introduced in the House is available
      <a href="https://www.govinfo.gov/bulkdata" target="_blank">here</a>. I
      thought it would be interesting to focus on which members were
      co-sponsoring legislation together, especially across party lines.
      Combining sponsorship data with the subject or policy areas covered by the
      legislation gives the user an idea of who is willing to work together
      despite party affiliation on particular issues. That said, this
      application was originally built to learn Vue in the context of a
      full-stack application and not a side-hustle to disrupt the lobbying
      industry. I later extended this application to serve as a portfolio
      project which exhibits full-stack skills like data processing, cloud
      computing, and API development.
    </p>
    <h3>Stack</h3>
    <v-list>
      <v-list-item v-for="(s, i) in stack" :key="`stack-${i}`">
        <v-list-item-content>
          <v-list-item-title>{{ s }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-container>
</template>

<script>
export default {
  data: () => ({
    stack: [
      'Go for parallelized data processing',
      'MongoDB (Atlas) for persistance',
      'Go for the API',
      'Vue for the client-side rendered SPA',
      'Vuex for frontend state management',
      'D3 for data visualization',
      'Vuetify for UI',
      'Axios for XHR',
      'Nginx for secure web hosting',
      'AWS for application hosting',
    ],
  }),
}
</script>

<style scoped></style>
