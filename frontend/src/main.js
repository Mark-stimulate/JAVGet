import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

const key = sessionStorage.getItem('javget_api_key') || prompt('请输入 JAVGet 访问密钥') || ''
if (key) {
    sessionStorage.setItem('javget_api_key', key)
    document.cookie = `javget_key=${encodeURIComponent(key)}; Path=/file/; SameSite=Strict`
}

const app = createApp(App)
app.use(router)
app.mount('#app')

// 全局错误处理
app.config.errorHandler = (err) => {
    console.error('Global error:', err)
}
