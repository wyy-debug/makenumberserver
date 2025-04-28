<template>
  <div id="app">
    <router-view />
    <tab-bar v-if="$route.meta.showTabBar"/>
  </div>
</template>

<script>
import TabBar from '@/components/TabBar.vue'

export default {
  name: 'App',
  components: {
    TabBar
  }
}
</script>

<style lang="scss">
/* 全局样式会通过main.scss引入 */
</style> 