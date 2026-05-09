import { defineStore } from 'pinia'
import { api } from '../api/api'

export const usePostStore = defineStore('posts', {
  state: () => ({
    posts: [],
    loading: false,
    error: null,
    hasMore: true
  }),
  actions: {
    async fetchPosts(groupId = 0, limit = 10, offset = 0) {
      this.loading = true
      try {
        const query = `limit=${limit}&offset=${offset}${groupId ? `&group_id=${groupId}` : ''}`
        const data = await api.get(`/posts?${query}`)
        if (offset === 0) {
          this.posts = data.posts || []
        } else {
          this.posts = [...this.posts, ...(data.posts || [])]
        }
        this.hasMore = data.has_more
      } catch (err) {
        this.error = err.message
        console.error('Failed to fetch posts:', err)
      } finally {
        this.loading = false
      }
    },
    async createPost(content, privacy = 'public', groupId = 0, allowedUsers = [], imageFile = null) {
      try {
        let data;
        if (imageFile) {
          const formData = new FormData()
          formData.append('content', content)
          formData.append('privacy', privacy)
          formData.append('group_id', groupId)
          if (allowedUsers.length > 0) {
            formData.append('allowed_users', JSON.stringify(allowedUsers))
          }
          formData.append('image', imageFile)
          data = await api.post('/posts', formData)
        } else {
          data = await api.post('/posts', {
            content,
            privacy,
            group_id: groupId,
            allowed_users: allowedUsers
          })
        }
        this.posts.unshift(data.post)
        return data.post
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async addComment(postId, content, imageFile = null) {
      try {
        let comment;
        if (imageFile) {
          const formData = new FormData()
          formData.append('post_id', postId)
          formData.append('content', content)
          formData.append('image', imageFile)
          comment = await api.post('/comments', formData)
        } else {
          comment = await api.post('/comments', { post_id: postId, content })
        }
        
        const post = this.posts.find(p => p.id === postId)
        if (post) {
          if (!post.comments) post.comments = []
          post.comments.push(comment)
        }
        return comment
      } catch (err) {
        this.error = err.message
        throw err
      }
    }
  }
})
