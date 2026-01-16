<!-- 登录页面 -->
<template>
  <div class="login-page">
    <div class="login-container">
      <!-- Logo -->
      <div class="logo-section">
        <div class="logo-icon">🎮</div>
        <h1 class="logo-title">Life RPG</h1>
        <p class="logo-subtitle">将生活游戏化，成为更好的自己</p>
      </div>

      <!-- 登录表单 -->
      <el-form ref="formRef" :model="form" :rules="rules" class="login-form">
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名"
            size="large"
            :prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            size="large"
            :prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-btn"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 快捷入口 -->
      <div class="quick-links">
        <el-button link type="primary" @click="fillAdminAccount">
          管理员演示账号
        </el-button>
        <el-divider direction="vertical" />
        <el-button link @click="showRegister = true">
          注册账号
        </el-button>
      </div>

      <!-- 设备提示 -->
      <div class="device-tips">
        <p>💻 管理端请使用 PC 浏览器访问</p>
        <p>📱 用户端支持手机浏览器访问</p>
      </div>
    </div>

    <!-- 注册对话框 -->
    <el-dialog v-model="showRegister" title="注册账号" width="400px">
      <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules">
        <el-form-item prop="username">
          <el-input v-model="registerForm.username" placeholder="用户名 (至少3个字符)" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="密码 (至少6个字符)"
            show-password
          />
        </el-form-item>
        <el-form-item prop="nickname">
          <el-input v-model="registerForm.nickname" placeholder="昵称 (可选)" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRegister = false">取消</el-button>
        <el-button type="primary" :loading="registerLoading" @click="handleRegister">
          注册
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { authApi } from '@/api'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

// 登录表单
const formRef = ref<FormInstance>()
const form = reactive({
  username: '',
  password: '',
})
const loading = ref(false)

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

// 注册表单
const showRegister = ref(false)
const registerFormRef = ref<FormInstance>()
const registerForm = reactive({
  username: '',
  password: '',
  nickname: '',
})
const registerLoading = ref(false)

const registerRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, message: '用户名至少3个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' },
  ],
}

// 填充管理员账号
const fillAdminAccount = () => {
  form.username = 'admin'
  form.password = '123456'
}

// 登录
const handleLogin = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.login(form.username, form.password)
    await userStore.fetchMenus()
    
    ElMessage.success('登录成功')
    
    // 根据角色跳转
    if (userStore.isAdmin) {
      router.push('/admin')
    } else {
      router.push('/app')
    }
  } catch {
    // 错误已在 API 拦截器处理
  } finally {
    loading.value = false
  }
}

// 注册
const handleRegister = async () => {
  const valid = await registerFormRef.value?.validate().catch(() => false)
  if (!valid) return

  registerLoading.value = true
  try {
    await authApi.register(registerForm)
    ElMessage.success('注册成功，请登录')
    showRegister.value = false
    form.username = registerForm.username
    form.password = ''
  } catch {
    // 错误已在 API 拦截器处理
  } finally {
    registerLoading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 400px;
  background: #fff;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.logo-section {
  text-align: center;
  margin-bottom: 40px;
}

.logo-icon {
  font-size: 60px;
  margin-bottom: 16px;
}

.logo-title {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  margin: 0;
}

.logo-subtitle {
  color: #888;
  margin-top: 8px;
  font-size: 14px;
}

.login-form {
  margin-bottom: 24px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
}

.quick-links {
  text-align: center;
  margin-bottom: 24px;
}

.device-tips {
  text-align: center;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  font-size: 12px;
  color: #666;
}

.device-tips p {
  margin: 4px 0;
}
</style>
