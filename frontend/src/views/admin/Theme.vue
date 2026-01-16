<!-- H5主题配置页面 -->
<template>
  <div class="theme-config">
    <el-card shadow="never">
      <template #header>
        <span>H5端主题配置</span>
      </template>

      <el-form :model="form" label-width="120px" style="max-width: 600px">
        <el-form-item label="主色调">
          <el-color-picker v-model="form.primaryColor" />
          <span class="color-preview" :style="{ background: form.primaryColor }"></span>
          <span class="color-tip">页面主色调，用于导航栏、按钮等</span>
        </el-form-item>
        <el-form-item label="辅助色">
          <el-color-picker v-model="form.secondaryColor" />
          <span class="color-preview" :style="{ background: form.secondaryColor }"></span>
          <span class="color-tip">辅助色，用于强调元素</span>
        </el-form-item>
        <el-form-item label="金币颜色">
          <el-color-picker v-model="form.goldColor" />
          <span class="color-preview" :style="{ background: form.goldColor }"></span>
          <span class="color-tip">金币图标和数字颜色</span>
        </el-form-item>
        <el-form-item label="经验颜色">
          <el-color-picker v-model="form.expColor" />
          <span class="color-preview" :style="{ background: form.expColor }"></span>
          <span class="color-tip">经验条和相关元素颜色</span>
        </el-form-item>
        <el-form-item label="背景图URL">
          <el-input v-model="form.backgroundUrl" placeholder="可选，H5首页背景图" />
        </el-form-item>
        <el-form-item label="Logo URL">
          <el-input v-model="form.logoUrl" placeholder="可选，H5端Logo图片" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSave">保存配置</el-button>
          <el-button @click="handleReset">恢复默认</el-button>
        </el-form-item>
      </el-form>

      <!-- 预览区 -->
      <el-divider>预览效果</el-divider>
      <div class="preview-area" :style="previewStyle">
        <div class="preview-header">Life RPG</div>
        <div class="preview-card">
          <div class="preview-level">Lv.10</div>
          <div class="preview-stats">
            <span :style="{ color: form.goldColor }">🪙 1280</span>
            <span :style="{ color: form.expColor }">⭐ 3500</span>
          </div>
        </div>
        <div class="preview-btn" :style="{ background: form.primaryColor }">完成任务</div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { themeApi } from '@/api'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const form = reactive({
  primaryColor: '#1989fa',
  secondaryColor: '#ff976a',
  goldColor: '#ffd700',
  expColor: '#07c160',
  backgroundUrl: '',
  logoUrl: '',
})

const defaultTheme = {
  primaryColor: '#1989fa',
  secondaryColor: '#ff976a',
  goldColor: '#ffd700',
  expColor: '#07c160',
  backgroundUrl: '',
  logoUrl: '',
}

const previewStyle = computed(() => ({
  background: `linear-gradient(180deg, ${form.primaryColor} 0%, #f7f8fa 50%)`,
}))

const fetchData = async () => {
  try {
    const data: any = await themeApi.get()
    if (data) {
      Object.assign(form, data)
    }
  } catch { /* ignore */ }
}

const handleSave = async () => {
  loading.value = true
  try {
    await themeApi.update(form)
    ElMessage.success('保存成功')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  Object.assign(form, defaultTheme)
}

onMounted(() => fetchData())
</script>

<style scoped>
.color-preview {
  display: inline-block;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  margin-left: 12px;
  vertical-align: middle;
  border: 1px solid #ddd;
}

.color-tip {
  margin-left: 12px;
  color: #888;
  font-size: 12px;
}

.preview-area {
  width: 300px;
  height: 400px;
  border-radius: 20px;
  padding: 20px;
  margin: 0 auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.preview-header {
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  color: #fff;
  margin-bottom: 20px;
}

.preview-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  margin-bottom: 20px;
}

.preview-level {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 12px;
}

.preview-stats {
  display: flex;
  justify-content: space-around;
  font-size: 16px;
  font-weight: 600;
}

.preview-btn {
  text-align: center;
  padding: 12px;
  border-radius: 8px;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
}
</style>
