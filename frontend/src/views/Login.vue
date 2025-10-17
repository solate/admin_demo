<template>
  <div class="login-page">
    <el-card class="login-card">
      <h2 class="title">后台管理系统</h2>
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item label="账号" prop="username">
          <el-input v-model="form.username" placeholder="请输入账号" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" :type="showPwd ? 'text' : 'password'" placeholder="请输入密码">
            <template #suffix>
              <el-icon @click="showPwd = !showPwd" class="clickable">
                <component :is="showPwd ? 'View' : 'Hide'" />
              </el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="onSubmit">登录</el-button>
          <el-button type="default">注册</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
  
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

interface LoginForm { username: string; password: string }

const router = useRouter()
const formRef = ref()
const form = ref<LoginForm>({ username: '', password: '' })
const showPwd = ref(false)
const loading = ref(false)

const rules = {
  username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

async function onSubmit() {
  await formRef.value?.validate()
  loading.value = true
  try {
    // TODO: 替换为实际登录 API
    if (form.value.username && form.value.password) {
      localStorage.setItem('token', 'mock-token')
      ElMessage.success('登录成功')
      router.push('/')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
}
.login-card { width: 360px; }
.title { text-align: center; margin-bottom: 16px; }
.clickable { cursor: pointer; }
</style>


