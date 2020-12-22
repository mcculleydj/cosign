<template>
  <v-navigation-drawer permanent app>
    <v-row align="center">
      <v-col cols="auto" class="ml-3 pr-0">
        <v-img src="logo.png" width="60" />
      </v-col>
      <v-col>
        <div class="ma-2 pl-0 brand">
          Cosign
        </div>
      </v-col>
    </v-row>
    <v-divider />
    <v-list nav>
      <v-list-item
        v-for="route in routes"
        :key="route.text"
        link
        :to="{ name: route.name }"
        color="primary"
        :exact="route.name === 'home'"
        :class="route.name === $route.name ? 'router-link-active' : ''"
      >
        <v-list-item-icon>
          <v-icon>{{ route.icon }}</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>
            {{ route.text }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
import { routes } from '@/router'
export default {
  data: () => ({
    routes,
  }),
}
</script>
