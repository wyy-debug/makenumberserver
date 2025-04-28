<template>
  <div class="container">
    <!-- 店铺封面 -->
    <div class="shop-cover mb-6">
      <img :src="shopInfo?.coverImage || '/images/default-cover.png'" />
      <div class="cover-gradient"></div>
      <div class="cover-info">
        <div class="shop-name">{{ shopInfo?.name }}</div>
      </div>
    </div>
    
    <!-- 关于我们 -->
    <div class="card mb-6">
      <div class="subtitle">关于我们</div>
      <div class="text-content mb-4">
        {{ shopInfo?.description || '暂无描述' }}
      </div>
      <div class="tags-container">
        <span class="tag">纯天然材料</span>
        <span class="tag">免费设计</span>
        <span class="tag">可定制图案</span>
        <span class="tag">环境舒适</span>
      </div>
    </div>
    
    <!-- 服务项目 -->
    <div class="card mb-6">
      <div class="subtitle">服务项目</div>
      <div class="service-list">
        <template v-if="services.length > 0">
          <div v-for="item in services" :key="item.id" class="service-item">
            <div>
              <div class="service-name">{{ item.name }}</div>
              <div class="text-sm text-gray">{{ item.description || '' }}</div>
            </div>
            <div class="text-right">
              <div class="text-xs text-gray">约{{ item.duration }}分钟</div>
            </div>
          </div>
        </template>
        <template v-else>
          <div class="empty-tip">暂无服务项目</div>
        </template>
      </div>
    </div>
    
    <!-- 店铺信息 -->
    <div class="card">
      <div class="subtitle">店铺信息</div>
      <div class="shop-info">
        <div class="info-item">
          <div class="info-icon location-icon"></div>
          <div>
            <div>{{ shopInfo?.address || '暂无地址' }}</div>
            <button class="map-btn" @click="openMap">查看地图</button>
          </div>
        </div>
        
        <div class="info-item">
          <div class="info-icon time-icon"></div>
          <div>营业时间: {{ shopInfo?.businessHours || '暂无信息' }}</div>
        </div>
        
        <div class="info-item">
          <div class="info-icon phone-icon"></div>
          <div>联系电话: {{ shopInfo?.phone || '暂无' }}</div>
          <button v-if="shopInfo?.phone" class="call-btn" @click="callPhone">拨打</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import { mapState } from 'vuex'

export default {
  name: 'Intro',
  data() {
    return {
      services: []
    }
  },
  computed: {
    ...mapState(['shopInfo', 'shopId'])
  },
  created() {
    this.getShopInfo()
    this.getServices()
  },
  methods: {
    getShopInfo() {
      this.$store.dispatch('getShopInfo', this.shopId)
    },
    
    async getServices() {
      try {
        const res = await axios.get(`/api/v1/shops/${this.shopId}/services`)
        this.services = res.data.data
      } catch (error) {
        console.error('获取服务列表失败:', error)
      }
    },
    
    openMap() {
      if (this.shopInfo?.latitude && this.shopInfo?.longitude) {
        // 因为是Web应用，我们使用高德地图或百度地图链接
        const url = `https://uri.amap.com/marker?position=${this.shopInfo.longitude},${this.shopInfo.latitude}&name=${this.shopInfo.name}`
        window.open(url, '_blank')
      } else {
        alert('无法获取店铺位置信息')
      }
    },
    
    callPhone() {
      if (this.shopInfo?.phone) {
        window.location.href = `tel:${this.shopInfo.phone}`
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.shop-cover {
  position: relative;
  height: 200px;
  border-radius: 12px;
  overflow: hidden;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .cover-gradient {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 100px;
    background: linear-gradient(to top, rgba(0,0,0,0.7), rgba(0,0,0,0));
  }
  
  .cover-info {
    position: absolute;
    bottom: 20px;
    left: 20px;
    right: 20px;
    color: white;
  }
  
  .shop-name {
    font-size: 24px;
    font-weight: bold;
  }
}

.text-content {
  line-height: 1.6;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  margin-top: 12px;
}

.service-list {
  margin-top: 16px;
}

.service-item {
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid var(--gray-200);
  
  &:last-child {
    border-bottom: none;
  }
}

.service-name {
  font-weight: bold;
  margin-bottom: 4px;
}

.empty-tip {
  text-align: center;
  color: var(--gray-500);
  padding: 20px 0;
}

.shop-info {
  margin-top: 16px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16px;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.info-icon {
  width: 24px;
  height: 24px;
  margin-right: 16px;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
}

.location-icon {
  background-image: url('@/assets/images/location.png');
}

.time-icon {
  background-image: url('@/assets/images/time.png');
}

.phone-icon {
  background-image: url('@/assets/images/phone.png');
}

.map-btn, .call-btn {
  background-color: var(--primary-light);
  color: var(--primary-color);
  border: none;
  border-radius: 4px;
  padding: 4px 8px;
  font-size: 12px;
  margin-top: 4px;
  cursor: pointer;
}

.text-right {
  text-align: right;
}
</style> 