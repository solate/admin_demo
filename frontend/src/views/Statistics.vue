<template>
  <div class="statistics">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-value">¥89,234</div>
            <div class="stat-label">商品总额</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-value">156</div>
            <div class="stat-label">商品数量</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-value">1,234</div>
            <div class="stat-label">入库数量</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-value">2,456</div>
            <div class="stat-label">库存数量</div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>入库趋势</span>
              <el-select v-model="timeRange" style="width: 120px;">
                <el-option label="日" value="daily" />
                <el-option label="周" value="weekly" />
                <el-option label="月" value="monthly" />
              </el-select>
            </div>
          </template>
          <div class="chart-placeholder">
            <el-empty description="入库趋势图表待开发" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>出库趋势</span>
              <el-select v-model="timeRange" style="width: 120px;">
                <el-option label="日" value="daily" />
                <el-option label="周" value="weekly" />
                <el-option label="月" value="monthly" />
              </el-select>
            </div>
          </template>
          <div class="chart-placeholder">
            <el-empty description="出库趋势图表待开发" />
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card style="margin-top: 20px;">
      <template #header>
        <span>商品分类统计</span>
      </template>
      <el-table :data="categoryData" style="width: 100%;">
        <el-table-column prop="category" label="分类" />
        <el-table-column prop="count" label="商品数量" />
        <el-table-column prop="totalValue" label="总价值">
          <template #default="{ row }">
            ¥{{ row.totalValue.toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="avgPrice" label="平均价格">
          <template #default="{ row }">
            ¥{{ row.avgPrice.toFixed(2) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const timeRange = ref('daily')

const categoryData = ref([
  { category: '手机', count: 45, totalValue: 269955, avgPrice: 5999 },
  { category: '电脑', count: 32, totalValue: 223968, avgPrice: 6999 },
  { category: '配件', count: 79, totalValue: 23621, avgPrice: 299 }
])

onMounted(() => {
  // 加载统计数据
})
</script>

<style scoped>
.stat-card {
  text-align: center;
}

.stat-item {
  padding: 20px 0;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-placeholder {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
