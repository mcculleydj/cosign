<template>
  <v-navigation-drawer permanent app>
    <div class="ma-2 brand">
      Cosign
    </div>
    <v-divider />
    <v-list nav>
      <v-list-item
        v-for="route in routes"
        :key="route.text"
        link
        :to="{ name: route.name }"
        exact
        color="primary"
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
