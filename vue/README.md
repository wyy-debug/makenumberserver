# 海娜纹身排队系统 (Vue版本)

这是一个为海娜纹身店开发的排队系统的Vue实现版本，从微信小程序转换而来。

## 功能特性

- 用户排队取号
- 实时查看排队状态
- 浏览店铺信息和服务项目
- 图案库展示与搜索

## 技术栈

- Vue 3
- Vuex
- Vue Router
- Axios
- Sass

## 项目结构

```
├── public/               # 静态资源
├── src/
│   ├── assets/           # 项目资源文件
│   │   └── styles/       # 样式文件
│   ├── components/       # 公共组件
│   ├── router/           # 路由配置
│   ├── store/            # Vuex状态管理
│   ├── views/            # 页面视图组件
│   ├── App.vue           # 根组件
│   └── main.js           # 入口文件
├── package.json          # 项目依赖
└── vue.config.js         # Vue配置
```

## 安装与运行

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run serve

# 构建生产版本
npm run build
```

## 配置说明

API服务器地址可以在以下文件中修改：

- `vue.config.js`（开发环境代理）
- `src/store/index.js`（baseUrl配置）

## 设计说明

本项目UI设计遵循简洁现代的风格，主色调采用紫色系，符合海娜纹身的品牌形象。响应式设计确保在各种移动设备上的良好显示效果。 