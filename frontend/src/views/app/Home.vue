<!-- H5首页 -->
<template>
  <div class="app-home">
    <!-- 头部用户信息 -->
    <div class="header">
      <div class="user-info">
        <van-image
          round
          width="60"
          height="60"
          :src="userStore.userInfo?.avatar || 'https://fastly.jsdelivr.net/npm/@vant/assets/cat.jpeg'"
        />
        <div class="user-detail">
          <div class="nickname">{{ userStore.userInfo?.nickname || userStore.username }}</div>
          <div class="level-badge" :style="{ background: themeVars.expColor }">
            Lv.{{ userStore.level }}
          </div>
        </div>
      </div>
      <div class="stats">
        <div class="stat-item">
          <span class="stat-value" :style="{ color: themeVars.goldColor }">🪙 {{ userStore.gold }}</span>
          <span class="stat-label">金币</span>
        </div>
      </div>
    </div>

    <!-- 经验条 -->
    <div class="exp-section">
      <div class="exp-header">
        <span>经验值</span>
        <span>{{ profile.expProgress }} / {{ profile.nextLevelExp }}</span>
      </div>
      <van-progress
        :percentage="profile.expPercentage || 0"
        :color="themeVars.expColor"
        :track-color="'rgba(255,255,255,0.3)'"
        stroke-width="8"
        :show-pivot="false"
      />
    </div>

    <!-- 公告轮播 -->
    <div class="announcement-section" v-if="announcements.length">
      <van-swipe :autoplay="5000" indicator-color="white" class="announcement-swipe">
        <van-swipe-item v-for="item in announcements" :key="item.id">
          <div class="announcement-item" @click="showAnnouncement(item)">
            <van-tag :type="{ notice: 'primary', activity: 'warning', update: 'success' }[item.type]" class="tag">
              {{ { notice: '通知', activity: '活动', update: '更新' }[item.type] }}
            </van-tag>
            <span class="title">{{ item.title }}</span>
          </div>
        </van-swipe-item>
      </van-swipe>
    </div>

    <!-- 快捷入口 -->
    <div class="quick-entry">
      <router-link to="/app/quest" class="entry-item">
        <div class="entry-icon">📋</div>
        <div class="entry-text">任务大厅</div>
      </router-link>
      <router-link to="/app/shop" class="entry-item">
        <div class="entry-icon">🛒</div>
        <div class="entry-text">奖励商店</div>
      </router-link>
      <router-link to="/app/profile" class="entry-item">
        <div class="entry-icon">📊</div>
        <div class="entry-text">我的记录</div>
      </router-link>
    </div>

    <!-- 今日任务 -->
    <div class="section">
      <div class="section-header">
        <span class="section-title">今日任务</span>
        <router-link to="/app/quest" class="section-more">查看全部</router-link>
      </div>
      <div class="task-list">
        <div
          v-for="task in todayTasks"
          :key="task.id"
          class="task-card"
          :class="{ completed: task.completed }"
        >
          <div class="task-icon">{{ task.icon }}</div>
          <div class="task-info">
            <div class="task-title">{{ task.title }}</div>
            <div class="task-reward">
              <span :style="{ color: themeVars.goldColor }">+{{ task.goldReward }}🪙</span>
              <span :style="{ color: themeVars.expColor }">+{{ task.expReward }}⭐</span>
            </div>
          </div>
          <van-button
            v-if="!task.completed"
            type="primary"
            size="small"
            :loading="task.loading"
            @click="completeTask(task)"
          >完成</van-button>
          <van-tag v-else type="success">已完成</van-tag>
        </div>
      </div>
    </div>

    <!-- 公告详情弹窗 -->
    <van-dialog v-model:show="showAnnouncementDialog" :title="currentAnnouncement?.title">
      <div class="announcement-content">
        {{ currentAnnouncement?.content }}
      </div>
    </van-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { dashboardApi, taskApi, announcementApi, themeApi } from '@/api'
import { showSuccessToast, showConfirmDialog } from 'vant'

const userStore = useUserStore()

// 主题配置
const themeVars = reactive({
  primaryColor: '#1989fa',
  goldColor: '#ffd700',
  expColor: '#07c160',
})

// 用户资料
const profile = reactive({
  expProgress: 0,
  nextLevelExp: 100,
  expPercentage: 0,
})

// 公告
const announcements = ref<any[]>([])
const showAnnouncementDialog = ref(false)
const currentAnnouncement = ref<any>(null)

// 任务
const todayTasks = ref<any[]>([])

const fetchTheme = async () => {
  try {
    const data: any = await themeApi.get()
    if (data) {
      Object.assign(themeVars, data)
    }
  } catch { /* ignore */ }
}

const fetchProfile = async () => {
  try {
    const data: any = await dashboardApi.profile()
    Object.assign(profile, data)
    if (data.user) {
      userStore.updateUserStats(data.user.gold, data.user.exp, data.user.level)
    }
  } catch { /* ignore */ }
}

const fetchAnnouncements = async () => {
  try {
    const data: any = await announcementApi.userList()
    announcements.value = data || []
  } catch { /* ignore */ }
}

const fetchTasks = async () => {
  try {
    const data: any = await taskApi.userList()
    todayTasks.value = (data || []).slice(0, 5).map((t: any) => ({ ...t, loading: false }))
  } catch { /* ignore */ }
}

const showAnnouncement = (item: any) => {
  currentAnnouncement.value = item
  showAnnouncementDialog.value = true
}

const completeTask = async (task: any) => {
  task.loading = true
  try {
    const result: any = await taskApi.complete(task.id)
    task.completed = true
    
    // 更新用户数据
    userStore.updateUserStats(result.newGold, result.newExp, result.newLevel)
    
    // 显示奖励动画
    if (result.levelUp) {
      await showConfirmDialog({
        title: '🎉 升级了！',
        message: `恭喜升到 Lv.${result.newLevel}！\n获得 ${result.goldReward}🪙 ${result.expReward}⭐`,
        confirmButtonText: '太棒了',
        showCancelButton: false,
      })
    } else {
      showSuccessToast(`+${result.goldReward}🪙 +${result.expReward}⭐`)
    }
    
    // 刷新资料
    fetchProfile()
  } catch { /* ignore */ } finally {
    task.loading = false
  }
}

onMounted(() => {
  fetchTheme()
  fetchProfile()
  fetchAnnouncements()
  fetchTasks()
})
</script>

<style scoped>
.app-home {
  padding: 20px 16px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  margin-bottom: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-detail {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.nickname {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.level-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 10px;
  font-size: 12px;
  color: #fff;
  font-weight: 600;
}

.stats {
  text-align: right;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
}

.stat-label {
  display: block;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
}

.exp-section {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 12px;
  padding: 12px 16px;
  margin-bottom: 16px;
}

.exp-header {
  display: flex;
  justify-content: space-between;
  color: #fff;
  font-size: 12px;
  margin-bottom: 8px;
}

.announcement-section {
  margin-bottom: 16px;
}

.announcement-swipe {
  border-radius: 12px;
  overflow: hidden;
}

.announcement-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #fff;
}

.announcement-item .title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.quick-entry {
  display: flex;
  justify-content: space-around;
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 16px;
}

.entry-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  color: #333;
}

.entry-icon {
  font-size: 32px;
}

.entry-text {
  font-size: 12px;
}

.section {
  background: #fff;
  border-radius: 16px;
  padding: 16px;
  margin-bottom: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
}

.section-more {
  font-size: 12px;
  color: #1989fa;
  text-decoration: none;
}

.task-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 12px;
  margin-bottom: 8px;
}

.task-card.completed {
  opacity: 0.6;
}

.task-icon {
  font-size: 28px;
}

.task-info {
  flex: 1;
}

.task-title {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.task-reward {
  display: flex;
  gap: 12px;
  font-size: 12px;
}

.announcement-content {
  padding: 16px;
  line-height: 1.6;
}
</style>
