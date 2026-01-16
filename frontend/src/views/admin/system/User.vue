<!-- 用户管理页面 -->
<template>
  <div class="user-manage">
    <el-card shadow="never">
      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchUsername"
          placeholder="搜索用户名"
          style="width: 200px"
          clearable
          @keyup.enter="handleSearch"
        />
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button type="primary" @click="handleAdd">新增用户</el-button>
      </div>

      <!-- 表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="row.role?.key === 'admin' ? 'danger' : 'primary'">
              {{ row.role?.name || '-' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="gold" label="金币" width="100" />
        <el-table-column prop="level" label="等级" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="warning" @click="handleResetPassword(row)">重置密码</el-button>
            <el-button link type="danger" @click="handleDelete(row)" :disabled="row.id === 1">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="!!form.id" />
        </el-form-item>
        <el-form-item v-if="!form.id" label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="默认123456" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item label="角色" prop="roleId">
          <el-select v-model="form.roleId" style="width: 100%">
            <el-option v-for="role in roles" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { userApi, roleApi } from '@/api'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'

interface User {
  id?: number
  username: string
  password?: string
  nickname: string
  roleId: number
  role?: { id: number; name: string; key: string }
  gold: number
  level: number
  status: number
  createdAt: string
}

const loading = ref(false)
const tableData = ref<User[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const searchUsername = ref('')

// 角色列表
const roles = ref<any[]>([])

// 表单
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const form = reactive<User>({
  id: undefined,
  username: '',
  password: '',
  nickname: '',
  roleId: 2,
  gold: 0,
  level: 1,
  status: 1,
  createdAt: '',
})
const submitLoading = ref(false)

const dialogTitle = computed(() => (form.id ? '编辑用户' : '新增用户'))

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  roleId: [{ required: true, message: '请选择角色', trigger: 'change' }],
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const fetchData = async () => {
  loading.value = true
  try {
    const data: any = await userApi.list({
      page: page.value,
      pageSize: pageSize.value,
      username: searchUsername.value,
    })
    tableData.value = data.list || []
    total.value = data.total || 0
  } catch {
    // 错误处理
  } finally {
    loading.value = false
  }
}

const fetchRoles = async () => {
  try {
    const data: any = await roleApi.list()
    roles.value = data || []
  } catch {
    // 错误处理
  }
}

const handleSearch = () => {
  page.value = 1
  fetchData()
}

const handleAdd = () => {
  Object.assign(form, {
    id: undefined,
    username: '',
    password: '',
    nickname: '',
    roleId: 2,
    status: 1,
  })
  dialogVisible.value = true
}

const handleEdit = (row: User) => {
  Object.assign(form, { ...row, password: '' })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    // 只发送需要的字段，不发送 createdAt 等时间字段
    const submitData = {
      username: form.username,
      password: form.password,
      nickname: form.nickname,
      roleId: form.roleId,
      status: form.status,
    }
    
    if (form.id) {
      await userApi.update(form.id, submitData)
      ElMessage.success('更新成功')
    } else {
      await userApi.create(submitData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch {
    // 错误处理
  } finally {
    submitLoading.value = false
  }
}

const handleResetPassword = async (row: User) => {
  await ElMessageBox.confirm('确定要重置该用户密码为123456吗？', '提示', { type: 'warning' })
  try {
    await userApi.resetPassword(row.id!)
    ElMessage.success('密码已重置为123456')
  } catch {
    // 错误处理
  }
}

const handleDelete = async (row: User) => {
  await ElMessageBox.confirm('确定要删除该用户吗？', '提示', { type: 'warning' })
  try {
    await userApi.delete(row.id!)
    ElMessage.success('删除成功')
    fetchData()
  } catch {
    // 错误处理
  }
}

onMounted(() => {
  fetchData()
  fetchRoles()
})
</script>

<style scoped>
.search-bar {
  margin-bottom: 16px;
  display: flex;
  gap: 12px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
