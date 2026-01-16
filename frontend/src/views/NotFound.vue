<!-- 404页面 -->
<template>
  <div class="not-found">
    <div class="content">
      <div class="icon">🎮</div>
      <h1>404</h1>
      <p>页面走丢了...</p>
      <div class="actions">
        <el-button type="primary" @click="goHome">返回首页</el-button>
        <el-button @click="goBack">返回上一页</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const goHome = () => {
  if (userStore.isAdmin) {
    router.push('/admin')
  } else if (userStore.isLoggedIn) {
    router.push('/app')
  } else {
    router.push('/login')
  }
}

const goBack = () => {
  router.back()
}
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.content {
  text-align: center;
  color: #fff;
}

.icon {
  font-size: 80px;
  margin-bottom: 20px;
}

h1 {
  font-size: 80px;
  margin: 0;
  font-weight: 700;
}

p {
  font-size: 20px;
  margin: 20px 0 40px;
  opacity: 0.8;
}

.actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}
</style>
