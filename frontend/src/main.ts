// 应用入口
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// Element Plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// Vant
import {
    Button, NavBar, Tabbar, TabbarItem, Image as VanImage,
    Tag, Swipe, SwipeItem, Progress, Empty, PullRefresh,
    Tab, Tabs, List, Overlay, ActionSheet, Popup, Cell, CellGroup,
    Icon, Dialog, showToast, showSuccessToast, showConfirmDialog
} from 'vant'
import 'vant/lib/index.css'

// 全局样式
import './styles/index.css'

const app = createApp(App)

// 使用 Pinia
app.use(createPinia())

// 使用路由
app.use(router)

// 使用 Element Plus
app.use(ElementPlus, { locale: zhCn })

// 注册 Element Plus 图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

// 注册 Vant 组件
app.use(Button)
app.use(NavBar)
app.use(Tabbar)
app.use(TabbarItem)
app.use(VanImage)
app.use(Tag)
app.use(Swipe)
app.use(SwipeItem)
app.use(Progress)
app.use(Empty)
app.use(PullRefresh)
app.use(Tab)
app.use(Tabs)
app.use(List)
app.use(Overlay)
app.use(ActionSheet)
app.use(Popup)
app.use(Cell)
app.use(CellGroup)
app.use(Icon)
app.use(Dialog)

app.mount('#app')
