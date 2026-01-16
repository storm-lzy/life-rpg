<!-- 任务管理页面 -->
<template>
  <div class="task-manage">
    <el-card shadow="never">
      <div class="search-bar">
        <el-select v-model="searchType" placeholder="任务类型" clearable style="width: 120px">
          <el-option label="每日任务" value="daily" />
          <el-option label="一次性" value="once" />
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button type="primary" @click="handleAdd">新增任务</el-button>
      </div>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="icon" label="图标" width="60" />
        <el-table-column prop="title" label="任务名称" width="150" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="goldReward" label="金币奖励" width="100">
          <template #default="{ row }">
            <span style="color: #ffd700;">🪙 {{ row.goldReward }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="expReward" label="经验奖励" width="100">
          <template #default="{ row }">
            <span style="color: #07c160;">⭐ {{ row.expReward }}</span>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'daily' ? 'primary' : 'warning'">
              {{ row.type === 'daily' ? '每日' : '一次性' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.isActive ? 'success' : 'info'">
              {{ row.isActive ? '上架' : '下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="150">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="任务名称" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-input v-model="form.icon" placeholder="支持Emoji，如🏃" />
        </el-form-item>
        <el-form-item label="任务类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio value="daily">每日任务</el-radio>
            <el-radio value="once">一次性任务</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-input v-model="form.category" placeholder="如：健康、学习、工作" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="金币奖励" prop="goldReward">
              <el-input-number v-model="form.goldReward" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="经验奖励" prop="expReward">
              <el-input-number v-model="form.expReward" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态" prop="isActive">
          <el-switch v-model="form.isActive" active-text="上架" inactive-text="下架" />
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
import { taskApi } from '@/api'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const searchType = ref('')

const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const form = reactive({
  id: undefined as number | undefined,
  title: '',
  description: '',
  icon: '',
  type: 'daily',
  category: '',
  goldReward: 10,
  expReward: 5,
  sort: 0,
  isActive: true,
})
const submitLoading = ref(false)

const dialogTitle = computed(() => (form.id ? '编辑任务' : '新增任务'))

const rules: FormRules = {
  title: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
}

const fetchData = async () => {
  loading.value = true
  try {
    const data: any = await taskApi.list({
      page: page.value,
      pageSize: pageSize.value,
      type: searchType.value,
    })
    tableData.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  page.value = 1
  fetchData()
}

const handleAdd = () => {
  Object.assign(form, {
    id: undefined, title: '', description: '', icon: '📝',
    type: 'daily', category: '', goldReward: 10, expReward: 5, sort: 0, isActive: true,
  })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    if (form.id) {
      await taskApi.update(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await taskApi.create(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm('确定要删除该任务吗？', '提示', { type: 'warning' })
  try {
    await taskApi.delete(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch { /* ignore */ }
}

onMounted(() => fetchData())
</script>

<style scoped>
.search-bar { margin-bottom: 16px; display: flex; gap: 12px; }
.pagination { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
