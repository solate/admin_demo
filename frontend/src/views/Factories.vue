<template>
  <div class="factories">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>工厂管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新建工厂
          </el-button>
        </div>
      </template>
      
      <div class="search-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索工厂名称或地址"
          style="width: 300px;"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
      
      <el-table :data="tableData" style="width: 100%; margin-top: 16px;">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="工厂名称" />
        <el-table-column prop="address" label="地址" />
        <el-table-column prop="owner" label="负责人" />
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button size="small" type="primary" plain @click="handleEdit(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button size="small" type="danger" plain @click="handleDelete(row)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 16px; text-align: right;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>
    
    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑工厂' : '新建工厂'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="工厂名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入工厂名称" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="form.address" placeholder="请输入地址" />
        </el-form-item>
        <el-form-item label="负责人" prop="owner">
          <el-input v-model="form.owner" placeholder="请输入负责人" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Delete } from '@element-plus/icons-vue'

interface Factory {
  id: number
  name: string
  address: string
  owner: string
  createdAt: string
}

const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

const tableData = ref<Factory[]>([])

const form = reactive({
  id: 0,
  name: '',
  address: '',
  owner: ''
})

const rules = {
  name: [{ required: true, message: '请输入工厂名称', trigger: 'blur' }],
  address: [{ required: true, message: '请输入地址', trigger: 'blur' }],
  owner: [{ required: true, message: '请输入负责人', trigger: 'blur' }]
}

// 模拟数据
const mockData: Factory[] = [
  { id: 1, name: '北京工厂', address: '北京市朝阳区', owner: '张三', createdAt: '2024-01-15 10:30:00' },
  { id: 2, name: '上海工厂', address: '上海市浦东新区', owner: '李四', createdAt: '2024-01-16 14:20:00' },
  { id: 3, name: '广州工厂', address: '广州市天河区', owner: '王五', createdAt: '2024-01-17 09:15:00' }
]

onMounted(() => {
  loadData()
})

function loadData() {
  // 模拟API调用
  tableData.value = mockData
  total.value = mockData.length
}

function handleSearch() {
  // 模拟搜索
  loadData()
}

function handleAdd() {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
}

function handleEdit(row: Factory) {
  isEdit.value = true
  dialogVisible.value = true
  Object.assign(form, row)
}

function handleDelete(row: Factory) {
  ElMessageBox.confirm('确定要删除这个工厂吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('删除成功')
    loadData()
  })
}

function handleSubmit() {
  formRef.value?.validate((valid: boolean) => {
    if (valid) {
      ElMessage.success(isEdit.value ? '编辑成功' : '新建成功')
      dialogVisible.value = false
      loadData()
    }
  })
}

function resetForm() {
  Object.assign(form, {
    id: 0,
    name: '',
    address: '',
    owner: ''
  })
}

function handleSizeChange(val: number) {
  pageSize.value = val
  loadData()
}

function handleCurrentChange(val: number) {
  currentPage.value = val
  loadData()
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-bar {
  margin-bottom: 16px;
}

.action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}

.action-buttons .el-button {
  margin: 0;
  padding: 4px 8px;
  font-size: 12px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.action-buttons .el-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.action-buttons .el-button .el-icon {
  margin-right: 2px;
  font-size: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .action-buttons {
    flex-direction: column;
    gap: 2px;
  }
  
  .action-buttons .el-button {
    width: 100%;
    justify-content: center;
  }
}
</style>
