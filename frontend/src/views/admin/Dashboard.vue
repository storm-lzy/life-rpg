<!-- 管理端仪表盘 -->
<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :xs="12" :sm="6">
        <div class="stat-card users">
          <div class="stat-icon">👥</div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.userCount }}</div>
            <div class="stat-label">用户总数</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card gold">
          <div class="stat-icon">🪙</div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.todayGold }}</div>
            <div class="stat-label">今日产出金币</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card tasks">
          <div class="stat-icon">✅</div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.todayTasks }}</div>
            <div class="stat-label">今日完成任务</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card active">
          <div class="stat-icon">🎯</div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.activeTaskCount }}</div>
            <div class="stat-label">活跃任务数</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 图表区 -->
    <el-row :gutter="20" class="charts">
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <span>每日金币产出趋势</span>
          </template>
          <div class="chart-placeholder">
            <div class="bar-chart">
              <div
                v-for="(item, index) in stats.dailyGoldStats"
                :key="index"
                class="bar-item"
              >
                <div
                  class="bar"
                  :style="{ height: getBarHeight(item.gold, maxGold) + 'px' }"
                ></div>
                <div class="bar-label">{{ formatDate(item.date) }}</div>
                <div class="bar-value">{{ item.gold }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <span>每日任务完成趋势</span>
          </template>
          <div class="chart-placeholder">
            <div class="bar-chart">
              <div
                v-for="(item, index) in stats.dailyTaskStats"
                :key="index"
                class="bar-item"
              >
                <div
                  class="bar task-bar"
                  :style="{ height: getBarHeight(item.count, maxTask) + 'px' }"
                ></div>
                <div class="bar-label">{{ formatDate(item.date) }}</div>
                <div class="bar-value">{{ item.count }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷入口 -->
    <el-card class="quick-actions" shadow="never">
      <template #header>
        <span>快捷操作</span>
      </template>
      <el-row :gutter="20">
        <el-col :span="6">
          <router-link to="/admin/game/task" class="quick-link">
            <el-icon :size="32"><List /></el-icon>
            <span>任务管理</span>
          </router-link>
        </el-col>
        <el-col :span="6">
          <router-link to="/admin/game/reward" class="quick-link">
            <el-icon :size="32"><Present /></el-icon>
            <span>奖励管理</span>
          </router-link>
        </el-col>
        <el-col :span="6">
          <router-link to="/admin/announcement" class="quick-link">
            <el-icon :size="32"><Bell /></el-icon>
            <span>发布公告</span>
          </router-link>
        </el-col>
        <el-col :span="6">
          <router-link to="/admin/system/user" class="quick-link">
            <el-icon :size="32"><User /></el-icon>
            <span>用户管理</span>
          </router-link>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { dashboardApi } from '@/api'
import { List, Present, Bell, User } from '@element-plus/icons-vue'

interface Stats {
  userCount: number
  todayGold: number
  todayTasks: number
  activeTaskCount: number
  activeRewardCount: number
  dailyGoldStats: { date: string; gold: number }[]
  dailyTaskStats: { date: string; count: number }[]
}

const stats = ref<Stats>({
  userCount: 0,
  todayGold: 0,
  todayTasks: 0,
  activeTaskCount: 0,
  activeRewardCount: 0,
  dailyGoldStats: [],
  dailyTaskStats: [],
})

const maxGold = computed(() => {
  const values = stats.value.dailyGoldStats.map((i) => i.gold)
  return Math.max(...values, 1)
})

const maxTask = computed(() => {
  const values = stats.value.dailyTaskStats.map((i) => i.count)
  return Math.max(...values, 1)
})

const getBarHeight = (value: number, max: number) => {
  return Math.max((value / max) * 120, 4)
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

onMounted(async () => {
  try {
    const data: any = await dashboardApi.stats()
    stats.value = data
  } catch {
    // 错误处理
  }
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
}

.stat-cards {
  margin-bottom: 20px;
}

.stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.stat-icon {
  font-size: 40px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #888;
  margin-top: 4px;
}

.chart-card {
  margin-bottom: 20px;
}

.chart-placeholder {
  height: 200px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.bar-chart {
  display: flex;
  align-items: flex-end;
  gap: 16px;
  height: 180px;
  width: 100%;
  justify-content: space-around;
  padding: 0 10px;
}

.bar-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.bar {
  width: 40px;
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px 4px 0 0;
  transition: height 0.3s;
}

.bar.task-bar {
  background: linear-gradient(180deg, #07c160 0%, #00a854 100%);
}

.bar-label {
  font-size: 12px;
  color: #888;
}

.bar-value {
  font-size: 12px;
  font-weight: 600;
  color: #333;
  position: absolute;
  top: -20px;
}

.bar-item {
  position: relative;
}

.quick-actions {
  margin-bottom: 20px;
}

.quick-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 24px;
  background: #f8f9fa;
  border-radius: 12px;
  text-decoration: none;
  color: #333;
  transition: all 0.2s;
}

.quick-link:hover {
  background: #667eea;
  color: #fff;
}

.quick-link span {
  font-size: 14px;
}
</style>
