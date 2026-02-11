<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { toast } from 'vue-sonner'
import { adminUsersApi, type AdminUser } from '@/api/admin/users'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const users = ref<AdminUser[]>([])
const loading = ref(false)
const saving = ref(false)
const showCreateModal = ref(false)
const editingUser = ref<AdminUser | null>(null)

const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0,
  has_next: false,
  has_prev: false
})

const form = ref({
  email: '',
  username: '',
  password: '',
  role: 'user'
})

const isFormValid = computed(() => {
  if (editingUser.value) {
    return form.value.email || form.value.username || form.value.role
  }
  return form.value.email && form.value.username && form.value.password && form.value.role
})

onMounted(() => fetchUsers(1))

async function fetchUsers(page = 1) {
  loading.value = true
  try {
    const response = await adminUsersApi.getUsers(page, 20)
    users.value = response.data.data
    pagination.value = response.data
  } catch (error: any) {
    toast.error(error.response?.data?.error || 'Failed to load users')
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingUser.value = null
  form.value = {
    email: '',
    username: '',
    password: '',
    role: 'user'
  }
  showCreateModal.value = true
}

function editUser(user: AdminUser) {
  editingUser.value = user
  form.value = {
    email: user.email,
    username: user.username,
    password: '',
    role: user.role
  }
  showCreateModal.value = true
}

async function saveUser() {
  if (!isFormValid.value) {
    toast.error('Please fill in all required fields')
    return
  }

  saving.value = true
  try {
    if (editingUser.value) {
      const updateData: any = {}
      if (form.value.email && form.value.email !== editingUser.value.email) {
        updateData.email = form.value.email
      }
      if (form.value.username && form.value.username !== editingUser.value.username) {
        updateData.username = form.value.username
      }
      if (form.value.role && form.value.role !== editingUser.value.role) {
        updateData.role = form.value.role
      }

      if (Object.keys(updateData).length === 0) {
        toast.error('No changes to save')
        return
      }

      await adminUsersApi.updateUser(editingUser.value.id, updateData)
      toast.success('User updated successfully')
    } else {
      await adminUsersApi.createUser(form.value)
      toast.success('User created successfully')
    }
    closeModal()
    await fetchUsers(pagination.value.page)
  } catch (error: any) {
    toast.error(error.response?.data?.error || 'Failed to save user')
  } finally {
    saving.value = false
  }
}

async function deleteUser(userId: number) {
  if (userStore.user?.id === userId) {
    toast.error('You cannot delete your own account')
    return
  }

  if (!confirm('Are you sure you want to delete this user? This action cannot be undone.')) {
    return
  }

  try {
    await adminUsersApi.deleteUser(userId)
    toast.success('User deleted successfully')
    await fetchUsers(pagination.value.page)
  } catch (error: any) {
    toast.error(error.response?.data?.error || 'Failed to delete user')
  }
}

function closeModal() {
  showCreateModal.value = false
  editingUser.value = null
  form.value = { email: '', username: '', password: '', role: 'user' }
}

function changePage(page: number) {
  fetchUsers(page)
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}
</script>

<template>
  <div class="space-y-8 p-8">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">User Management</h2>
        <p class="text-gray-400 mt-1">Manage system users and permissions</p>
      </div>
      <button
        @click="openCreateModal"
        class="px-6 py-2 bg-purple-600 text-white rounded-xl hover:bg-purple-700 transition-colors"
      >
        Create User
      </button>
    </div>

    <!-- Users List -->
    <div class="bg-white rounded-3xl border border-gray-100 shadow-sm overflow-hidden">
      <table class="w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-sm font-medium text-gray-500">User</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-gray-500">Email</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-gray-500">Role</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-gray-500">GitHub</th>
            <th class="px-6 py-3 text-left text-sm font-medium text-gray-500">Created</th>
            <th class="px-6 py-3 text-right text-sm font-medium text-gray-500">Actions</th>
          </tr>
        </thead>
        <tbody v-if="!loading && users.length > 0">
          <tr v-for="user in users" :key="user.id" class="border-t border-gray-100">
            <td class="px-6 py-4 font-medium text-gray-900">{{ user.username }}</td>
            <td class="px-6 py-4 text-gray-600">{{ user.email }}</td>
            <td class="px-6 py-4">
              <span
                :class="
                  user.role === 'admin'
                    ? 'bg-purple-100 text-purple-600'
                    : 'bg-gray-100 text-gray-600'
                "
                class="px-3 py-1 rounded-full text-sm font-medium"
              >
                {{ user.role }}
              </span>
            </td>
            <td class="px-6 py-4 text-gray-600">{{ user.github_username || '-' }}</td>
            <td class="px-6 py-4 text-gray-600">{{ formatDate(user.created_at) }}</td>
            <td class="px-6 py-4 text-right space-x-2">
              <button
                @click="editUser(user)"
                class="text-blue-600 hover:underline font-medium"
              >
                Edit
              </button>
              <button
                @click="deleteUser(user.id)"
                class="text-red-600 hover:underline font-medium"
                :disabled="userStore.user?.id === user.id"
                :class="{ 'opacity-50 cursor-not-allowed': userStore.user?.id === user.id }"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <LoadingSpinner v-if="loading" class="py-12" />

      <div v-if="!loading && users.length === 0" class="py-12 text-center text-gray-500">
        No users found
      </div>

      <!-- Pagination -->
      <div v-if="pagination.total > 0" class="px-6 py-4 border-t border-gray-100 flex justify-between items-center">
        <div class="text-sm text-gray-600">
          Showing {{ (pagination.page - 1) * pagination.limit + 1 }} to
          {{ Math.min(pagination.page * pagination.limit, pagination.total) }} of
          {{ pagination.total }} users
        </div>
        <div class="flex gap-2">
          <button
            @click="changePage(pagination.page - 1)"
            :disabled="!pagination.has_prev"
            class="px-4 py-2 border border-gray-300 rounded-lg text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <button
            @click="changePage(pagination.page + 1)"
            :disabled="!pagination.has_next"
            class="px-4 py-2 border border-gray-300 rounded-lg text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>
    </div>

    <!-- Create/Edit User Modal -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 backdrop-blur-sm"
      @click.self="closeModal"
    >
      <div class="bg-white rounded-3xl p-8 max-w-md w-full mx-4 shadow-2xl">
        <h3 class="text-xl font-bold mb-6">
          {{ editingUser ? 'Edit User' : 'Create User' }}
        </h3>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
            <input
              v-model="form.email"
              type="email"
              class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple-500"
              :placeholder="editingUser ? 'Leave empty to keep current' : 'user@example.com'"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Username</label>
            <input
              v-model="form.username"
              type="text"
              class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple-500"
              :placeholder="editingUser ? 'Leave empty to keep current' : 'username'"
            />
          </div>

          <div v-if="!editingUser">
            <label class="block text-sm font-medium text-gray-700 mb-2">Password</label>
            <input
              v-model="form.password"
              type="password"
              class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple-500"
              placeholder="Minimum 8 characters"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Role</label>
            <select
              v-model="form.role"
              class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple-500"
            >
              <option value="user">User</option>
              <option value="admin">Admin</option>
            </select>
          </div>

          <div class="flex gap-2 pt-4">
            <button
              @click="saveUser"
              :disabled="saving || !isFormValid"
              class="flex-1 px-6 py-2 bg-purple-600 text-white rounded-xl hover:bg-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {{ saving ? 'Saving...' : 'Save' }}
            </button>
            <button
              @click="closeModal"
              class="flex-1 px-6 py-2 border border-gray-300 text-gray-700 rounded-xl hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
