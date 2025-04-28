import { createStore } from 'vuex'
import axios from 'axios'

export default createStore({
  state: {
    userInfo: null,
    shopInfo: null,
    token: localStorage.getItem('token') || '',
    baseUrl: '/api/v1',
    queueData: {
      hasNumber: false,
      queueNumber: null,
      peopleAhead: 0,
      waitTime: 0
    },
    shopId: 1
  },
  getters: {
    isAuthenticated: state => !!state.token,
    getQueueData: state => state.queueData,
    getShopInfo: state => state.shopInfo
  },
  mutations: {
    setUserInfo(state, userInfo) {
      state.userInfo = userInfo
    },
    setToken(state, token) {
      state.token = token
      localStorage.setItem('token', token)
    },
    clearToken(state) {
      state.token = ''
      localStorage.removeItem('token')
    },
    setQueueData(state, queueData) {
      state.queueData = queueData
    },
    setShopInfo(state, shopInfo) {
      state.shopInfo = shopInfo
    }
  },
  actions: {
    async login({ commit, dispatch }, code) {
      try {
        const res = await axios.post(`${this.state.baseUrl}/auth/login`, { code })
        const { token, user } = res.data.data
        commit('setToken', token)
        commit('setUserInfo', user)
        
        // 获取店铺信息
        dispatch('getShopInfo', this.state.shopId)
        
        // 获取排队状态
        dispatch('getQueueStatus', this.state.shopId)
        
        return user
      } catch (error) {
        console.error('登录失败:', error)
        throw error
      }
    },
    
    async getShopInfo({ commit }, shopId) {
      try {
        const res = await axios.get(`${this.state.baseUrl}/shops/${shopId}`)
        commit('setShopInfo', res.data.data)
        return res.data.data
      } catch (error) {
        console.error('获取店铺信息失败:', error)
        throw error
      }
    },
    
    async getQueueStatus({ commit }, shopId) {
      if (!this.state.token) return
      
      try {
        const res = await axios.get(`${this.state.baseUrl}/queue/status`, {
          params: { shop_id: shopId },
          headers: { Authorization: `Bearer ${this.state.token}` }
        })
        commit('setQueueData', res.data.data)
        return res.data.data
      } catch (error) {
        console.error('获取排队状态失败:', error)
        throw error
      }
    },
    
    // 设置请求拦截器以处理token
    setupAxios({ state, commit }) {
      axios.interceptors.request.use(
        config => {
          if (state.token) {
            config.headers.Authorization = `Bearer ${state.token}`
          }
          return config
        },
        error => Promise.reject(error)
      )
      
      axios.interceptors.response.use(
        response => response,
        error => {
          if (error.response && error.response.status === 401) {
            commit('clearToken')
            commit('setUserInfo', null)
          }
          return Promise.reject(error)
        }
      )
    }
  }
}) 