import client from '../client'

export interface AdminUser {
  id: number
  email: string
  username: string
  role: string
  github_username?: string
  created_at: string
  updated_at: string
}

export const adminUsersApi = {
  getUsers(page = 1, limit = 20) {
    return client.get('/admin/users', { params: { page, limit } })
  },

  getUser(id: number) {
    return client.get(`/admin/users/${id}`)
  },

  createUser(data: { email: string; username: string; password: string; role: string }) {
    return client.post('/admin/users', data)
  },

  updateUser(id: number, data: { email?: string; username?: string; role?: string }) {
    return client.put(`/admin/users/${id}`, data)
  },

  deleteUser(id: number) {
    return client.delete(`/admin/users/${id}`)
  }
}
