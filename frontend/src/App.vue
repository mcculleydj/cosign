<template>
  <v-app>
    <NavDrawer />
    <v-main class="main">
      <router-view />
    </v-main>
  </v-app>
</template>

<script>
import NavDrawer from '@/components/NavDrawer'

export default {
  components: {
    NavDrawer,
  },
}
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=Arizonia');

.brand {
  font-family: 'Arizonia';
  font-size: 44px;
  margin-top: 8px;
}

/* boo this hack! boo! */
.v-autocomplete__content {
  max-width: 0px !important;
}
</style>

<style scoped>
.main {
  background-color: rgba(0, 0, 0, 0.05);
}
</style>
