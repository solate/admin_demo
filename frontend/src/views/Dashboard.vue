<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-value">12</div>
            <div class="stat-label">工厂数量</div>
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
            <div class="stat-value">¥89,234</div>
            <div class="stat-label">商品总额</div>
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
    
    <el-card class="chart-card" style="margin-top: 20px;">
      <template #header>
        <span>数据趋势</span>
      </template>
      <div class="chart-placeholder">
        <el-empty description="图表功能待开发" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
// 首页仪表板
</script>

<style scoped>
.dashboard {
  padding: 0;
}

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

.chart-card {
  min-height: 300px;
}

.chart-placeholder {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
