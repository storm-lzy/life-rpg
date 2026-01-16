<!-- H5移动端布局 -->
<template>
  <div class="mobile-layout" :style="themeStyle">
    <!-- 主内容区 -->
    <div class="mobile-content">
      <router-view v-slot="{ Component }">
        <transition name="slide" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </div>

    <!-- 底部导航栏 -->
    <van-tabbar v-model="active" route :placeholder="true" class="tabbar">
      <van-tabbar-item to="/app/home" icon="home-o" name="home">首页</van-tabbar-item>
      <van-tabbar-item to="/app/quest" icon="todo-list-o" name="quest">任务</van-tabbar-item>
      <van-tabbar-item to="/app/shop" icon="shop-o" name="shop">商店</van-tabbar-item>
      <van-tabbar-item to="/app/profile" icon="user-o" name="profile">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { themeApi } from '@/api'

const active = ref('home')

// 主题配置
const theme = ref({
  primaryColor: '#1989fa',
  secondaryColor: '#ff976a',
  goldColor: '#ffd700',
  expColor: '#07c160',
})

const themeStyle = computed(() => ({
  '--primary-color': theme.value.primaryColor,
  '--secondary-color': theme.value.secondaryColor,
  '--gold-color': theme.value.goldColor,
  '--exp-color': theme.value.expColor,
}))

onMounted(async () => {
  try {
    const data: any = await themeApi.get()
    if (data) {
      theme.value = data
    }
  } catch {
    // 使用默认主题
  }
})
</script>

<style scoped>
.mobile-layout {
  min-height: 100vh;
  background: linear-gradient(180deg, var(--primary-color) 0%, #f7f8fa 30%);
}

.mobile-content {
  padding-bottom: 50px;
  min-height: calc(100vh - 50px);
}

.tabbar {
  border-top: 1px solid #eee;
}

/* 页面切换动画 */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.slide-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style>
