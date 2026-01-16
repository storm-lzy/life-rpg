<!-- H5任务大厅 -->
<template>
  <div class="app-quest">
    <van-nav-bar title="任务大厅" />

    <!-- 分类标签 -->
    <van-tabs v-model:active="activeTab" sticky>
      <van-tab title="全部" name="all" />
      <van-tab title="每日任务" name="daily" />
      <van-tab title="一次性" name="once" />
    </van-tabs>

    <!-- 任务列表 -->
    <div class="task-list">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-empty v-if="!filteredTasks.length" description="暂无任务" />
        
        <div
          v-for="task in filteredTasks"
          :key="task.id"
          class="task-card"
          :class="{ completed: task.completed }"
        >
          <div class="task-left">
            <div class="task-icon">{{ task.icon }}</div>
            <div class="task-info">
              <div class="task-title">{{ task.title }}</div>
              <div class="task-desc">{{ task.description }}</div>
              <div class="task-meta">
                <van-tag v-if="task.category" plain size="small">{{ task.category }}</van-tag>
                <van-tag :type="task.type === 'daily' ? 'primary' : 'warning'" size="small">
                  {{ task.type === 'daily' ? '每日' : '一次性' }}
                </van-tag>
              </div>
            </div>
          </div>
          <div class="task-right">
            <div class="task-reward">
              <div class="reward-item gold">+{{ task.goldReward }}🪙</div>
              <div class="reward-item exp">+{{ task.expReward }}⭐</div>
            </div>
            <van-button
              v-if="!task.completed"
              type="primary"
              size="small"
              round
              :loading="task.loading"
              @click="completeTask(task)"
            >
              完成任务
            </van-button>
            <van-tag v-else type="success" size="large">✓ 已完成</van-tag>
          </div>
        </div>
      </van-pull-refresh>
    </div>

    <!-- 完成动画 -->
    <van-overlay :show="showReward" @click="showReward = false">
      <div class="reward-popup">
        <div class="reward-icon">🎉</div>
        <div class="reward-title">任务完成！</div>
        <div class="reward-detail">
          <span class="gold">+{{ rewardInfo.goldReward }}🪙</span>
          <span class="exp">+{{ rewardInfo.expReward }}⭐</span>
        </div>
        <div v-if="rewardInfo.levelUp" class="level-up">
          🎊 恭喜升级到 Lv.{{ rewardInfo.newLevel }}！
        </div>
      </div>
    </van-overlay>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { taskApi } from '@/api'

const userStore = useUserStore()
const activeTab = ref('all')
const refreshing = ref(false)
const tasks = ref<any[]>([])

// 奖励弹窗
const showReward = ref(false)
const rewardInfo = ref({
  goldReward: 0,
  expReward: 0,
  levelUp: false,
  newLevel: 1,
})

const filteredTasks = computed(() => {
  if (activeTab.value === 'all') return tasks.value
  return tasks.value.filter((t) => t.type === activeTab.value)
})

const fetchTasks = async () => {
  try {
    const data: any = await taskApi.userList()
    tasks.value = (data || []).map((t: any) => ({ ...t, loading: false }))
  } catch { /* ignore */ }
}

const onRefresh = async () => {
  await fetchTasks()
  refreshing.value = false
}

const completeTask = async (task: any) => {
  task.loading = true
  try {
    const result: any = await taskApi.complete(task.id)
    task.completed = true
    
    // 更新用户数据
    userStore.updateUserStats(result.newGold, result.newExp, result.newLevel)
    
    // 显示奖励
    rewardInfo.value = result
    showReward.value = true
    
    setTimeout(() => {
      showReward.value = false
    }, 2500)
  } catch { /* ignore */ } finally {
    task.loading = false
  }
}

onMounted(() => fetchTasks())
</script>

<style scoped>
.app-quest {
  background: #f7f8fa;
  min-height: 100vh;
}

.task-list {
  padding: 16px;
}

.task-card {
  display: flex;
  justify-content: space-between;
  background: #fff;
  border-radius: 16px;
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.task-card.completed {
  opacity: 0.7;
}

.task-left {
  display: flex;
  gap: 12px;
  flex: 1;
}

.task-icon {
  font-size: 36px;
}

.task-info {
  flex: 1;
}

.task-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 4px;
}

.task-desc {
  font-size: 12px;
  color: #888;
  margin-bottom: 8px;
}

.task-meta {
  display: flex;
  gap: 6px;
}

.task-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: space-between;
}

.task-reward {
  text-align: right;
}

.reward-item {
  font-size: 14px;
  font-weight: 600;
}

.reward-item.gold {
  color: #ffd700;
}

.reward-item.exp {
  color: #07c160;
}

/* 奖励弹窗 */
.reward-popup {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #fff;
  border-radius: 20px;
  padding: 40px;
  text-align: center;
  animation: popIn 0.3s ease;
}

@keyframes popIn {
  0% { transform: translate(-50%, -50%) scale(0.5); opacity: 0; }
  100% { transform: translate(-50%, -50%) scale(1); opacity: 1; }
}

.reward-icon {
  font-size: 60px;
  margin-bottom: 16px;
}

.reward-title {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 16px;
}

.reward-detail {
  display: flex;
  justify-content: center;
  gap: 24px;
  font-size: 20px;
  font-weight: 600;
}

.reward-detail .gold {
  color: #ffd700;
}

.reward-detail .exp {
  color: #07c160;
}

.level-up {
  margin-top: 16px;
  font-size: 16px;
  color: #ff976a;
  font-weight: 600;
}
</style>
