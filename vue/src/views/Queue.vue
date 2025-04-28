<template>
  <div class="container">
    <div class="title">海娜纹身排队</div>
    
    <!-- 当前排队状态 -->
    <div class="status-card card">
      <div class="text-center">
        <div class="subtitle mb-1">当前状态</div>
        <div :class="['status-text', queueData.hasNumber ? 'active' : '']">
          {{ queueData.hasNumber ? '排队中' : '未取号' }}
        </div>
        <template v-if="queueData.hasNumber">
          <div class="mb-1">您的号码: {{ queueData.queueNumber }}</div>
          <div class="mb-1">前方等待: {{ queueData.peopleAhead }} 人</div>
          <div>预计等待: {{ queueData.waitTime }} 分钟</div>
        </template>
      </div>
    </div>
    
    <!-- 取号区域 -->
    <div class="card">
      <template v-if="!queueData.hasNumber">
        <div class="subtitle">选择服务</div>
        <div class="mb-4">
          <div v-for="item in services" :key="item.id" class="radio-item">
            <label>
              <input 
                type="radio" 
                :value="item.id" 
                v-model="selectedService" 
              />
              <span>{{ item.name }} (约{{ item.duration }}分钟)</span>
            </label>
          </div>
        </div>
        <button 
          class="btn btn-primary w-full" 
          @click="getNumber" 
          :disabled="loading"
        >
          {{ loading ? '处理中...' : '立即取号' }}
        </button>
      </template>
      
      <template v-else>
        <div class="text-center">
          <div class="queue-number">{{ queueData.queueNumber }}</div>
          <div class="mb-4">请留意叫号提醒</div>
          <button 
            class="btn btn-danger" 
            @click="cancelQueue" 
            :disabled="loading"
          >
            {{ loading ? '处理中...' : '取消排队' }}
          </button>
        </div>
      </template>
    </div>
    
    <!-- 当前叫号情况 -->
    <div class="card">
      <div class="subtitle">当前叫号</div>
      <div class="call-container">
        <div class="call-status">
          <div class="call-number text-success">{{ currentQueue.currentServing || '-' }}</div>
          <div class="text-sm text-gray">正在服务</div>
        </div>
        <div class="call-status">
          <div class="call-number text-info">{{ currentQueue.currentWaiting || '-' }}</div>
          <div class="text-sm text-gray">请就位</div>
        </div>
      </div>
      <div class="text-xs text-gray">
        今日已服务: {{ currentQueue.totalServed || 0 }} 位顾客
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import { mapState } from 'vuex'

export default {
  name: 'Queue',
  data() {
    return {
      selectedService: '',
      services: [],
      currentQueue: {
        currentServing: '',
        currentWaiting: '',
        totalServed: 0
      },
      loading: false
    }
  },
  computed: {
    ...mapState(['queueData', 'shopId', 'token'])
  },
  created() {
    this.getServices()
    this.getCurrentQueue()
  },
  mounted() {
    this.$store.dispatch('getQueueStatus', this.shopId)
  },
  methods: {
    // 获取服务列表
    async getServices() {
      try {
        this.loading = true
        const res = await axios.get(`/api/v1/shops/${this.shopId}/services`)
        this.services = res.data.data
        this.selectedService = this.services.length > 0 ? this.services[0].id : ''
      } catch (error) {
        console.error('获取服务列表失败:', error)
        alert('获取服务列表失败')
      } finally {
        this.loading = false
      }
    },
    
    // 获取当前叫号情况
    async getCurrentQueue() {
      try {
        const res = await axios.get('/api/v1/queue/current', {
          params: { shop_id: this.shopId }
        })
        this.currentQueue = res.data.data
      } catch (error) {
        console.error('获取当前叫号情况失败:', error)
      }
    },

    // 取号
    async getNumber() {
      if (!this.token) {
        alert('请先登录')
        return
      }
      
      if (!this.selectedService) {
        alert('请选择服务')
        return
      }
      
      try {
        this.loading = true
        await axios.post('/api/v1/queue/number', {
          shop_id: this.shopId,
          service_id: this.selectedService
        })
        
        // 更新排队状态
        await this.$store.dispatch('getQueueStatus', this.shopId)
        await this.getCurrentQueue()
        
        alert('取号成功')
      } catch (error) {
        console.error('取号失败:', error)
        alert(error.response?.data?.message || '取号失败')
      } finally {
        this.loading = false
      }
    },

    // 取消排队
    async cancelQueue() {
      if (confirm('确定要取消当前排队吗？')) {
        try {
          this.loading = true
          await axios.delete('/api/v1/queue/number', {
            data: { shop_id: this.shopId }
          })
          
          // 更新排队状态
          await this.$store.dispatch('getQueueStatus', this.shopId)
          await this.getCurrentQueue()
          
          alert('已取消排队')
        } catch (error) {
          console.error('取消排队失败:', error)
          alert(error.response?.data?.message || '取消排队失败')
        } finally {
          this.loading = false
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.status-card {
  background-color: #fff;
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.status-text {
  font-size: 24px;
  font-weight: bold;
  margin: 10px 0;
  color: var(--gray-500);
  
  &.active {
    color: var(--primary-color);
  }
}

.radio-item {
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  
  input {
    margin-right: 10px;
  }
}

.queue-number {
  font-size: 48px;
  font-weight: bold;
  color: var(--primary-color);
  margin: 15px 0;
}

.call-container {
  display: flex;
  justify-content: space-around;
  margin: 20px 0;
}

.call-status {
  text-align: center;
}

.call-number {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 5px;
}
</style> 