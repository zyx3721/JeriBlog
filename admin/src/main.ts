/*
项目名称：JeriBlog
文件名称：main.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Admin 应用入口文件
*/

import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import './assets/css/main.scss'
import 'remixicon/fonts/remixicon.css'

const app = createApp(App)

app.use(router)
app.use(ElementPlus, {
  locale: zhCn,
})
app.mount('#app')