<!-- 角色管理页面 -->
<template>
  <div class="role-manage">
    <el-card shadow="never">
      <div class="toolbar">
        <el-button type="primary" @click="handleAdd">新增角色</el-button>
      </div>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="key" label="角色标识" width="120" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column label="操作" fixed="right" width="220">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="success" @click="handleAssignMenus(row)">分配菜单</el-button>
            <el-button link type="danger" @click="handleDelete(row)" :disabled="row.id <= 2">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="标识" prop="key">
          <el-input v-model="form.key" :disabled="!!form.id" />
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
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 分配菜单对话框 -->
    <el-dialog v-model="menuDialogVisible" title="分配菜单" width="500px">
      <el-tree
        ref="menuTreeRef"
        :data="menuTree"
        :props="{ label: 'name', children: 'children' }"
        show-checkbox
        node-key="id"
        default-expand-all
      />
      <template #footer>
        <el-button @click="menuDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="menuSubmitLoading" @click="handleSubmitMenus">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { roleApi, menuApi } from '@/api'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { ElTree } from 'element-plus'

const loading = ref(false)
const tableData = ref<any[]>([])

// 表单
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const form = reactive({
  id: undefined as number | undefined,
  name: '',
  key: '',
  sort: 0,
  status: 1,
  remark: '',
})
const submitLoading = ref(false)

const dialogTitle = computed(() => (form.id ? '编辑角色' : '新增角色'))

const rules: FormRules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  key: [{ required: true, message: '请输入角色标识', trigger: 'blur' }],
}

// 菜单分配
const menuDialogVisible = ref(false)
const menuTreeRef = ref<InstanceType<typeof ElTree>>()
const menuTree = ref<any[]>([])
const currentRoleId = ref<number>()
const menuSubmitLoading = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    const data: any = await roleApi.list()
    tableData.value = data || []
  } finally {
    loading.value = false
  }
}

const fetchMenuTree = async () => {
  try {
    const data: any = await menuApi.tree()
    menuTree.value = data || []
  } catch {
    // 错误处理
  }
}

const handleAdd = () => {
  Object.assign(form, { id: undefined, name: '', key: '', sort: 0, status: 1, remark: '' })
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
      await roleApi.update(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await roleApi.create(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm('确定要删除该角色吗？', '提示', { type: 'warning' })
  try {
    await roleApi.delete(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch {
    // 错误处理
  }
}

const handleAssignMenus = async (row: any) => {
  currentRoleId.value = row.id
  if (menuTree.value.length === 0) {
    await fetchMenuTree()
  }
  // 获取角色已有菜单
  try {
    const menuIds: any = await roleApi.getMenus(row.id)
    menuTreeRef.value?.setCheckedKeys(menuIds || [])
  } catch {
    // 错误处理
  }
  menuDialogVisible.value = true
}

const handleSubmitMenus = async () => {
  const checkedKeys = menuTreeRef.value?.getCheckedKeys(false) as number[]
  const halfCheckedKeys = menuTreeRef.value?.getHalfCheckedKeys() as number[]
  const menuIds = [...checkedKeys, ...halfCheckedKeys]

  menuSubmitLoading.value = true
  try {
    await roleApi.assignMenus(currentRoleId.value!, menuIds)
    ElMessage.success('分配成功')
    menuDialogVisible.value = false
  } finally {
    menuSubmitLoading.value = false
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
