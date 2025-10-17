<template>
  <div class="products">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>商品管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新建商品
          </el-button>
        </div>
      </template>
      
      <div class="search-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索商品名称或SKU"
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
        <el-table-column prop="name" label="商品名称" />
        <el-table-column prop="sku" label="SKU" width="120" />
        <el-table-column prop="category" label="类别" width="100" />
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">
            ¥{{ row.price }}
          </template>
        </el-table-column>
        <el-table-column prop="unit" label="单位" width="80" />
        <el-table-column prop="stock" label="库存" width="100">
          <template #default="{ row }">
            <el-tag :type="row.stock > 10 ? 'success' : row.stock > 0 ? 'warning' : 'danger'">
              {{ row.stock }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button size="small" type="primary" plain @click="handleEdit(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button size="small" type="success" plain @click="handleStockIn(row)">
                <el-icon><Plus /></el-icon>
                入库
              </el-button>
              <el-button size="small" type="warning" plain @click="handleStockOut(row)">
                <el-icon><Minus /></el-icon>
                出库
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
    
    <!-- 商品编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑商品' : '新建商品'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="SKU" prop="sku">
          <el-input v-model="form.sku" placeholder="请输入SKU" />
        </el-form-item>
        <el-form-item label="类别" prop="category">
          <el-input v-model="form.category" placeholder="请输入类别" />
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="form.price" :min="0" :precision="2" style="width: 100%;" />
        </el-form-item>
        <el-form-item label="单位" prop="unit">
          <el-input v-model="form.unit" placeholder="请输入单位" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
    
    <!-- 库存操作对话框 -->
    <el-dialog
      v-model="stockDialogVisible"
      :title="stockType === 'in' ? '商品入库' : '商品出库'"
      width="400px"
    >
      <el-form :model="stockForm" :rules="stockRules" ref="stockFormRef" label-width="80px">
        <el-form-item label="商品名称">
          <el-input v-model="currentProduct.name" disabled />
        </el-form-item>
        <el-form-item label="当前库存">
          <el-input v-model="currentProduct.stock" disabled />
        </el-form-item>
        <el-form-item :label="stockType === 'in' ? '入库数量' : '出库数量'" prop="quantity">
          <el-input-number v-model="stockForm.quantity" :min="1" style="width: 100%;" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="stockForm.remark" type="textarea" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="stockDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleStockSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Plus, Minus, Delete } from '@element-plus/icons-vue'

interface Product {
  id: number
  name: string
  sku: string
  category: string
  price: number
  unit: string
  stock: number
  remark: string
  createdAt: string
}

const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const dialogVisible = ref(false)
const stockDialogVisible = ref(false)
const isEdit = ref(false)
const stockType = ref<'in' | 'out'>('in')
const formRef = ref()
const stockFormRef = ref()

const tableData = ref<Product[]>([])
const currentProduct = ref<Product>({} as Product)

const form = reactive({
  id: 0,
  name: '',
  sku: '',
  category: '',
  price: 0,
  unit: '',
  remark: ''
})

const stockForm = reactive({
  quantity: 1,
  remark: ''
})

const rules = {
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  sku: [{ required: true, message: '请输入SKU', trigger: 'blur' }],
  category: [{ required: true, message: '请输入类别', trigger: 'blur' }],
  price: [{ required: true, message: '请输入价格', trigger: 'blur' }],
  unit: [{ required: true, message: '请输入单位', trigger: 'blur' }]
}

const stockRules = {
  quantity: [{ required: true, message: '请输入数量', trigger: 'blur' }]
}

// 模拟数据
const mockData: Product[] = [
  { id: 1, name: '苹果手机', sku: 'IPHONE-15', category: '手机', price: 5999, unit: '台', stock: 50, remark: '最新款', createdAt: '2024-01-15 10:30:00' },
  { id: 2, name: '华为笔记本', sku: 'HUAWEI-MATE', category: '电脑', price: 6999, unit: '台', stock: 25, remark: '办公专用', createdAt: '2024-01-16 14:20:00' },
  { id: 3, name: '小米耳机', sku: 'MI-AIRPODS', category: '配件', price: 299, unit: '个', stock: 5, remark: '无线蓝牙', createdAt: '2024-01-17 09:15:00' }
]

onMounted(() => {
  loadData()
})

function loadData() {
  tableData.value = mockData
  total.value = mockData.length
}

function handleSearch() {
  loadData()
}

function handleAdd() {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
}

function handleEdit(row: Product) {
  isEdit.value = true
  dialogVisible.value = true
  Object.assign(form, row)
}

function handleStockIn(row: Product) {
  stockType.value = 'in'
  currentProduct.value = row
  stockDialogVisible.value = true
  resetStockForm()
}

function handleStockOut(row: Product) {
  stockType.value = 'out'
  currentProduct.value = row
  stockDialogVisible.value = true
  resetStockForm()
}

function handleDelete(row: Product) {
  ElMessageBox.confirm('确定要删除这个商品吗？', '提示', {
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

function handleStockSubmit() {
  stockFormRef.value?.validate((valid: boolean) => {
    if (valid) {
      const action = stockType.value === 'in' ? '入库' : '出库'
      ElMessage.success(`${action}成功`)
      stockDialogVisible.value = false
      loadData()
    }
  })
}

function resetForm() {
  Object.assign(form, {
    id: 0,
    name: '',
    sku: '',
    category: '',
    price: 0,
    unit: '',
    remark: ''
  })
}

function resetStockForm() {
  Object.assign(stockForm, {
    quantity: 1,
    remark: ''
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
