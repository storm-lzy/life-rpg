<!-- 菜单管理页面 -->
<template>
  <div class="menu-manage">
    <el-card shadow="never">
      <div class="toolbar">
        <el-button type="primary" @click="handleAdd()">新增菜单</el-button>
      </div>

      <el-table :data="tableData" v-loading="loading" row-key="id" default-expand-all>
        <el-table-column prop="name" label="菜单名称" width="180" />
        <el-table-column prop="icon" label="图标" width="80">
          <template #default="{ row }">
            {{ row.icon || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由路径" width="200" />
        <el-table-column prop="component" label="组件路径" width="180" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="['info', 'success', 'warning'][row.type - 1]">
              {{ ['目录', '菜单', '按钮'][row.type - 1] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="180">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleAdd(row.id)">新增子菜单</el-button>
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="上级菜单" prop="parentId">
          <el-tree-select
            v-model="form.parentId"
            :data="menuTreeOptions"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            check-strictly
            placeholder="无上级"
            clearable
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="菜单类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio :value="1">目录</el-radio>
            <el-radio :value="2">菜单</el-radio>
            <el-radio :value="3">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item v-if="form.type !== 3" label="路由路径" prop="path">
          <el-input v-model="form.path" />
        </el-form-item>
        <el-form-item v-if="form.type === 2" label="组件路径" prop="component">
          <el-input v-model="form.component" />
        </el-form-item>
        <el-form-item v-if="form.type !== 3" label="图标" prop="icon">
          <el-input v-model="form.icon" placeholder="Element Plus 图标名称" />
        </el-form-item>
        <el-form-item v-if="form.type === 3" label="权限标识" prop="permission">
          <el-input v-model="form.permission" placeholder="如 system:user:add" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" />
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
import { menuApi } from '@/api'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'

const loading = ref(false)
const tableData = ref<any[]>([])

// 表单
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const form = reactive({
  id: undefined as number | undefined,
  parentId: 0,
  name: '',
  path: '',
  component: '',
  icon: '',
  sort: 0,
  type: 2,
  permission: '',
  visible: 1,
  status: 1,
})
const submitLoading = ref(false)

const dialogTitle = computed(() => (form.id ? '编辑菜单' : '新增菜单'))

const menuTreeOptions = computed(() => [{ id: 0, name: '顶级菜单', children: tableData.value }])

const rules: FormRules = {
  name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
}

const fetchData = async () => {
  loading.value = true
  try {
    const data: any = await menuApi.tree()
    tableData.value = data || []
  } finally {
    loading.value = false
  }
}

const handleAdd = (parentId: number = 0) => {
  Object.assign(form, {
    id: undefined,
    parentId,
    name: '',
    path: '',
    component: '',
    icon: '',
    sort: 0,
    type: 2,
    permission: '',
    status: 1,
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
      await menuApi.update(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await menuApi.create(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm('确定要删除该菜单吗？', '提示', { type: 'warning' })
  try {
    await menuApi.delete(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch {
    // 错误处理
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.toolbar {
  margin-bottom: 16px;
}
</style>
