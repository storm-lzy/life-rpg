<!-- H5商店 -->
<template>
  <div class="app-shop">
    <van-nav-bar title="奖励商店">
      <template #right>
        <span class="gold-display">🪙 {{ userStore.gold }}</span>
      </template>
    </van-nav-bar>

    <!-- 商品网格 -->
    <div class="shop-grid">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-empty v-if="!rewards.length" description="暂无商品" />
        
        <div class="reward-grid">
          <div
            v-for="reward in rewards"
            :key="reward.id"
            class="reward-card"
            @click="showDetail(reward)"
          >
            <div class="reward-image">
              {{ getEmoji(reward.category) }}
            </div>
            <div class="reward-info">
              <div class="reward-title">{{ reward.title }}</div>
              <div class="reward-stock" v-if="reward.stock !== -1">
                库存: {{ reward.stock }}
              </div>
              <div class="reward-price">
                <span class="price">🪙 {{ reward.cost }}</span>
              </div>
            </div>
          </div>
        </div>
      </van-pull-refresh>
    </div>

    <!-- 商品详情弹窗 -->
    <van-action-sheet v-model:show="showSheet" :title="currentReward?.title">
      <div class="reward-detail" v-if="currentReward">
        <div class="detail-icon">{{ getEmoji(currentReward.category) }}</div>
        <div class="detail-desc">{{ currentReward.description }}</div>
        <div class="detail-price">
          需要 <span class="price-value">{{ currentReward.cost }}</span> 🪙
        </div>
        <div class="detail-balance">
          当前余额: {{ userStore.gold }} 🪙
          <span v-if="userStore.gold < currentReward.cost" class="insufficient">(不足)</span>
        </div>
        <van-button
          type="primary"
          block
          round
          size="large"
          :disabled="userStore.gold < currentReward.cost || currentReward.stock === 0"
          :loading="purchasing"
          @click="handlePurchase"
        >
          {{ currentReward.stock === 0 ? '已售罄' : '立即兑换' }}
        </van-button>
      </div>
    </van-action-sheet>

    <!-- 购买成功 -->
    <van-overlay :show="showSuccess">
      <div class="success-popup">
        <div class="success-icon">🎁</div>
        <div class="success-title">兑换成功!</div>
        <div class="success-reward">{{ currentReward?.title }}</div>
        <div class="success-cost">消耗 {{ currentReward?.cost }} 🪙</div>
      </div>
    </van-overlay>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { rewardApi } from '@/api'
import { showConfirmDialog } from 'vant'

const userStore = useUserStore()
const refreshing = ref(false)
const rewards = ref<any[]>([])

// 详情弹窗
const showSheet = ref(false)
const currentReward = ref<any>(null)
const purchasing = ref(false)

// 成功动画
const showSuccess = ref(false)

// 根据分类获取emoji
const getEmoji = (category: string) => {
  const map: Record<string, string> = {
    '休闲': '☕',
    '美食': '🍜',
    '娱乐': '🎮',
    '旅行': '✈️',
    '购物': '🛍️',
  }
  return map[category] || '🎁'
}

const fetchRewards = async () => {
  try {
    const data: any = await rewardApi.userList()
    rewards.value = data || []
  } catch { /* ignore */ }
}

const onRefresh = async () => {
  await fetchRewards()
  refreshing.value = false
}

const showDetail = (reward: any) => {
  currentReward.value = reward
  showSheet.value = true
}

const handlePurchase = async () => {
  await showConfirmDialog({
    title: '确认兑换',
    message: `确定花费 ${currentReward.value.cost} 金币兑换 "${currentReward.value.title}"?`,
  })

  purchasing.value = true
  try {
    const result: any = await rewardApi.purchase(currentReward.value.id)
    
    // 更新用户金币
    userStore.updateUserStats(result.newGold, userStore.exp, userStore.level)
    
    // 关闭弹窗，显示成功
    showSheet.value = false
    showSuccess.value = true
    
    setTimeout(() => {
      showSuccess.value = false
    }, 2000)
    
    // 刷新列表
    fetchRewards()
  } catch { /* ignore */ } finally {
    purchasing.value = false
  }
}

onMounted(() => fetchRewards())
</script>

<style scoped>
.app-shop {
  background: #f7f8fa;
  min-height: 100vh;
}

.gold-display {
  font-weight: 600;
  color: #ffd700;
}

.shop-grid {
  padding: 16px;
}

.reward-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.reward-card {
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.reward-image {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.reward-info {
  padding: 12px;
}

.reward-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.reward-stock {
  font-size: 12px;
  color: #888;
  margin-bottom: 4px;
}

.reward-price .price {
  font-size: 16px;
  font-weight: 700;
  color: #ff976a;
}

/* 详情弹窗 */
.reward-detail {
  padding: 24px;
  text-align: center;
}

.detail-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.detail-desc {
  color: #666;
  margin-bottom: 16px;
  line-height: 1.6;
}

.detail-price {
  font-size: 20px;
  margin-bottom: 8px;
}

.price-value {
  font-weight: 700;
  color: #ff976a;
}

.detail-balance {
  font-size: 14px;
  color: #888;
  margin-bottom: 24px;
}

.insufficient {
  color: #ee0a24;
}

/* 成功动画 */
.success-popup {
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

.success-icon {
  font-size: 60px;
  margin-bottom: 16px;
}

.success-title {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 8px;
}

.success-reward {
  font-size: 18px;
  color: #333;
  margin-bottom: 8px;
}

.success-cost {
  font-size: 14px;
  color: #888;
}
</style>
