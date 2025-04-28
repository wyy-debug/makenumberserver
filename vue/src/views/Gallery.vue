<template>
  <div class="container">
    <div class="title mb-4">海娜纹身图案库</div>
    
    <!-- 搜索和筛选区域 -->
    <div class="mb-4">
      <!-- 搜索框 -->
      <div class="search-bar mb-3">
        <i class="search-icon"></i>
        <input 
          type="text" 
          placeholder="搜索图案..." 
          v-model="searchQuery"
          @input="handleSearch"
          @keyup.enter="handleSearch"
        />
      </div>
      
      <!-- 分类下拉框 -->
      <div class="category-dropdown">
        <select v-model="categoryIndex" @change="handleCategoryChange">
          <option v-for="(category, index) in categories" :key="index" :value="index">
            {{ category }}
          </option>
        </select>
      </div>
    </div>
    
    <!-- 图片网格 -->
    <div class="gallery-grid">
      <template v-if="galleryItems.length > 0">
        <div v-for="item in galleryItems" :key="item.id" class="gallery-item">
          <img :src="item.imageUrl" alt="图案" />
          <div class="gallery-item-info">
            <div class="gallery-item-title">{{ item.title }}</div>
            <div class="gallery-item-footer">
              <span class="tag">{{ item.category }}</span>
              <div class="like-btn" @click="toggleLike(item.id)">
                <img :src="item.isLiked ? require('@/assets/images/heart-filled.png') : require('@/assets/images/heart.png')" />
              </div>
            </div>
          </div>
        </div>
      </template>
      <template v-else>
        <div class="empty-tip">
          <span>暂无图案</span>
        </div>
      </template>
    </div>
    
    <button 
      v-if="hasMore" 
      class="btn w-full load-more" 
      @click="loadMore" 
      :disabled="loading"
    >
      {{ loading ? '加载中...' : '加载更多图案' }}
    </button>
    <div v-else class="end-tip">已加载全部图案</div>
  </div>
</template>

<script>
import axios from 'axios'
import { mapState } from 'vuex'

export default {
  name: 'Gallery',
  data() {
    return {
      galleryItems: [],
      categories: ['全部', '传统', '现代', '几何', '动物', '花卉', '自然'],
      categoryIndex: 0,
      searchQuery: '',
      page: 1,
      limit: 10,
      hasMore: true,
      loading: false
    }
  },
  computed: {
    ...mapState(['shopId', 'token'])
  },
  created() {
    this.loadGalleryItems()
  },
  methods: {
    async loadGalleryItems(reset = false) {
      if (reset) {
        this.page = 1
        this.galleryItems = []
      }
      
      const category = this.categoryIndex === 0 ? '' : this.categories[this.categoryIndex]
      
      try {
        this.loading = true
        const res = await axios.get('/api/v1/gallery', {
          params: {
            shop_id: this.shopId,
            page: this.page,
            limit: this.limit,
            category,
            keyword: this.searchQuery
          }
        })
        
        const { items, total } = res.data.data
        
        if (reset) {
          this.galleryItems = items
        } else {
          this.galleryItems = [...this.galleryItems, ...items]
        }
        
        this.hasMore = this.galleryItems.length < total
        this.page++
      } catch (error) {
        console.error('获取图库列表失败:', error)
      } finally {
        this.loading = false
      }
    },
    
    handleSearch() {
      this.loadGalleryItems(true)
    },
    
    handleCategoryChange() {
      this.loadGalleryItems(true)
    },
    
    loadMore() {
      this.loadGalleryItems()
    },
    
    async toggleLike(id) {
      if (!this.token) {
        alert('请先登录')
        return
      }
      
      try {
        const item = this.galleryItems.find(item => item.id === id)
        if (!item) return
        
        const isLiked = item.isLiked
        
        // 先在界面上更新状态
        item.isLiked = !isLiked
        
        // 发送请求到服务器
        if (isLiked) {
          await axios.delete(`/api/v1/gallery/${id}/like`)
        } else {
          await axios.post(`/api/v1/gallery/${id}/like`)
        }
      } catch (error) {
        console.error('操作失败:', error)
        // 操作失败，恢复原状态
        const item = this.galleryItems.find(item => item.id === id)
        if (item) {
          item.isLiked = !item.isLiked
        }
        alert('操作失败，请重试')
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.search-bar {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 8px;
  padding: 10px 15px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  
  input {
    flex: 1;
    border: none;
    outline: none;
    font-size: 14px;
  }
  
  .search-icon {
    width: 16px;
    height: 16px;
    margin-right: 10px;
    background-image: url('../assets/images/search.png');
    background-size: contain;
    background-repeat: no-repeat;
  }
}

.category-dropdown {
  background-color: #fff;
  border-radius: 8px;
  padding: 10px 15px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  
  select {
    width: 100%;
    border: none;
    outline: none;
    background: transparent;
    font-size: 14px;
    appearance: none;
    background-image: url('../assets/images/arrow-down.png');
    background-repeat: no-repeat;
    background-position: right center;
    padding-right: 20px;
  }
}

.gallery-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  margin-bottom: 20px;
}

.gallery-item {
  background-color: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  
  img {
    width: 100%;
    height: 150px;
    object-fit: cover;
  }
}

.gallery-item-info {
  padding: 10px;
}

.gallery-item-title {
  font-weight: bold;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.gallery-item-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.like-btn {
  width: 24px;
  height: 24px;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }
}

.empty-tip {
  grid-column: span 2;
  text-align: center;
  padding: 40px 0;
  color: var(--gray-500);
}

.load-more {
  background-color: var(--gray-200);
  color: var(--gray-700);
}

.end-tip {
  text-align: center;
  color: var(--gray-500);
  padding: 10px 0;
  font-size: 12px;
}
</style>