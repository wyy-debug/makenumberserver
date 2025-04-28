<template>
  <div class="tab-bar">
    <div 
      v-for="tab in tabs" 
      :key="tab.path" 
      class="tab-item" 
      :class="{ active: isActive(tab.path) }"
      @click="navigateTo(tab.path)"
    >
      <div class="tab-icon">
        <img :src="isActive(tab.path) ? tab.activeIcon : tab.icon" />
      </div>
      <div class="tab-text">{{ tab.text }}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TabBar',
  data() {
    return {
      tabs: [
        {
          path: '/queue',
          text: '排队',
          icon: '/images/icons/queue.png',
          activeIcon: '/images/icons/queue-active.png'
        },
        {
          path: '/intro',
          text: '介绍',
          icon: '/images/icons/intro.png',
          activeIcon: '/images/icons/intro-active.png'
        },
        {
          path: '/gallery',
          text: '图库',
          icon: '/images/icons/gallery.png',
          activeIcon: '/images/icons/gallery-active.png'
        }
      ]
    }
  },
  methods: {
    isActive(path) {
      return this.$route.path === path
    },
    navigateTo(path) {
      if (this.$route.path !== path) {
        this.$router.push(path)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.tab-bar {
  display: flex;
  height: 50px;
  background-color: #ffffff;
  border-top: 1px solid var(--gray-200);
}

.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #8a8a8a;
  
  &.active {
    color: var(--primary-color);
  }
}

.tab-icon {
  width: 24px;
  height: 24px;
  margin-bottom: 4px;
  
  img {
    width: 100%;
    height: 100%;
  }
}

.tab-text {
  font-size: 12px;
}
</style> 