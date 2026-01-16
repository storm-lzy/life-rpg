<!-- H5个人中心 -->
<template>
  <div class="app-profile">
    <van-nav-bar title="我的" />

    <!-- 用户卡片 -->
    <div class="user-card">
      <van-image
        round
        width="70"
        height="70"
        :src="userStore.userInfo?.avatar || 'https://fastly.jsdelivr.net/npm/@vant/assets/cat.jpeg'"
      />
      <div class="user-info">
        <div class="nickname">{{ userStore.userInfo?.nickname || userStore.username }}</div>
        <div class="level">Lv.{{ userStore.level }}</div>
      </div>
      <div class="stats-row">
        <div class="stat">
          <div class="stat-value gold">{{ userStore.gold }}</div>
          <div class="stat-label">金币</div>
        </div>
        <div class="stat">
          <div class="stat-value exp">{{ userStore.exp }}</div>
          <div class="stat-label">经验</div>
        </div>
      </div>
    </div>

    <!-- 功能列表 -->
    <van-cell-group inset class="menu-group">
      <van-cell title="金币记录" is-link @click="activeTab = 'gold'; showLogs = true">
        <template #icon><span class="cell-icon">🪙</span></template>
      </van-cell>
      <van-cell title="经验记录" is-link @click="activeTab = 'exp'; showLogs = true">
        <template #icon><span class="cell-icon">⭐</span></template>
      </van-cell>
    </van-cell-group>

    <van-cell-group inset class="menu-group">
      <van-cell title="退出登录" is-link @click="handleLogout">
        <template #icon><span class="cell-icon">🚪</span></template>
      </van-cell>
    </van-cell-group>

    <!-- 流水记录弹窗 -->
    <van-popup v-model:show="showLogs" position="bottom" :style="{ height: '70%' }" round>
      <div class="logs-popup">
        <div class="logs-header">
          <span>{{ activeTab === 'gold' ? '金币记录' : '经验记录' }}</span>
          <van-icon name="cross" @click="showLogs = false" />
        </div>
        
        <van-tabs v-model:active="logType" v-if="activeTab === 'gold'">
          <van-tab title="全部" name="" />
          <van-tab title="获得" name="gold_in" />
          <van-tab title="消耗" name="gold_out" />
        </van-tabs>

        <div class="logs-list">
          <van-pull-refresh v-model="refreshingLogs" @refresh="fetchLogs">
            <van-list
              v-model:loading="loadingLogs"
              :finished="logsFinished"
              @load="loadMoreLogs"
            >
              <div v-for="log in logs" :key="log.id" class="log-item">
                <div class="log-icon">
                  {{ log.type === 'gold_in' ? '🪙' : log.type === 'gold_out' ? '💸' : '⭐' }}
                </div>
                <div class="log-info">
                  <div class="log-desc">{{ log.description }}</div>
                  <div class="log-time">{{ formatTime(log.createdAt) }}</div>
                </div>
                <div
                  class="log-amount"
                  :class="{ positive: log.type.endsWith('_in'), negative: log.type.endsWith('_out') }"
                >
                  {{ log.type.endsWith('_in') ? '+' : '-' }}{{ log.amount }}
                </div>
              </div>
              <van-empty v-if="!logs.length && !loadingLogs" description="暂无记录" />
            </van-list>
          </van-pull-refresh>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { dashboardApi } from '@/api'
import { showConfirmDialog } from 'vant'

const router = useRouter()
const userStore = useUserStore()

// 流水弹窗
const showLogs = ref(false)
const activeTab = ref('gold')
const logType = ref('')
const logs = ref<any[]>([])
const logsPage = ref(1)
const loadingLogs = ref(false)
const refreshingLogs = ref(false)
const logsFinished = ref(false)

const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

const fetchLogs = async () => {
  logsPage.value = 1
  logsFinished.value = false
  
  try {
    const typeParam = activeTab.value === 'gold' ? logType.value : 'exp_in'
    const data: any = await dashboardApi.logs({ page: 1, pageSize: 20, type: typeParam })
    logs.value = data.list || []
    logsFinished.value = logs.value.length >= (data.total || 0)
  } catch { /* ignore */ } finally {
    refreshingLogs.value = false
  }
}

const loadMoreLogs = async () => {
  logsPage.value++
  try {
    const typeParam = activeTab.value === 'gold' ? logType.value : 'exp_in'
    const data: any = await dashboardApi.logs({ page: logsPage.value, pageSize: 20, type: typeParam })
    logs.value.push(...(data.list || []))
    logsFinished.value = logs.value.length >= (data.total || 0)
  } catch { /* ignore */ } finally {
    loadingLogs.value = false
  }
}

watch([showLogs, logType], ([show]) => {
  if (show) {
    fetchLogs()
  }
})

const handleLogout = async () => {
  await showConfirmDialog({
    title: '确认退出',
    message: '确定要退出登录吗?',
  })
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-profile {
  background: #f7f8fa;
  min-height: 100vh;
}

.user-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  margin: 16px;
  border-radius: 20px;
  padding: 24px;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 16px;
}

.user-info {
  flex: 1;
}

.nickname {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 4px;
}

.level {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  background: rgba(255, 255, 255, 0.2);
  padding: 2px 10px;
  border-radius: 10px;
  display: inline-block;
}

.stats-row {
  width: 100%;
  display: flex;
  gap: 24px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
}

.stat {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
}

.stat-value.gold {
  color: #ffd700;
}

.stat-value.exp {
  color: #7efff5;
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
}

.menu-group {
  margin: 16px;
}

.cell-icon {
  font-size: 20px;
  margin-right: 8px;
}

/* 流水弹窗 */
.logs-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  font-size: 16px;
  font-weight: 600;
}

.logs-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 16px;
}

.log-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.log-icon {
  font-size: 24px;
}

.log-info {
  flex: 1;
}

.log-desc {
  font-size: 14px;
}

.log-time {
  font-size: 12px;
  color: #888;
}

.log-amount {
  font-size: 16px;
  font-weight: 600;
}

.log-amount.positive {
  color: #07c160;
}

.log-amount.negative {
  color: #ee0a24;
}
</style>
