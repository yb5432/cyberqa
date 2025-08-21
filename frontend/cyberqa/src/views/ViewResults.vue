<template>
  <div class="max-w-4xl mx-auto p-4">
    <h1 class="text-3xl font-bold mb-6 text-center bg-clip-text text-transparent bg-gradient-to-r from-pink-400 via-purple-300 to-blue-400">赛博问缘 - 查看结果</h1>
    
    <div class="p-8 rounded-xl bg-black bg-opacity-30 backdrop-blur-lg border border-white border-opacity-20 shadow-lg">
      <h2 class="text-2xl font-semibold mb-6 text-center">输入TOKEN或邀请链接</h2>
      
      <form @submit.prevent="viewResults" class="mb-6">
        <div class="mb-6">
          <label for="tokenInput" class="block text-lg font-medium mb-2">TOKEN或邀请链接</label>
          <input
            id="tokenInput"
            v-model="tokenInput"
            type="text"
            placeholder="请输入TOKEN或完整的邀请链接"
            class="w-full px-4 py-3 bg-black bg-opacity-20 border border-white border-opacity-10 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 text-white"
          />
          <p class="mt-2 text-sm text-gray-400">可以是类似 "abc123xyz" 的TOKEN，也可以是完整的邀请链接</p>
        </div>
        
        <button
          type="submit"
          :disabled="!tokenInput.trim()"
          class="w-full py-4 px-6 bg-gradient-to-r from-blue-500 to-cyan-600 text-white font-bold text-lg rounded-xl shadow-lg hover:from-blue-600 hover:to-cyan-700 transform hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 focus:ring-cyan-400 focus:ring-opacity-50 disabled:opacity-50"
        >
          查看结果
        </button>
      </form>
      
      <div v-if="error" class="mt-6 p-4 rounded-lg bg-red-900 bg-opacity-30 border border-red-500 border-opacity-20">
        <p class="text-center text-red-200">{{ error }}</p>
      </div>
    </div>
    
    <div class="mt-8 text-center">
      <router-link
        to="/"
        class="inline-block py-3 px-6 bg-gradient-to-r from-purple-500 to-pink-600 text-white font-bold text-lg rounded-xl shadow-lg hover:from-purple-600 hover:to-pink-700 transform hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 focus:ring-purple-400 focus:ring-opacity-50"
      >
        返回首页
      </router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ViewResults',
  data() {
    return {
      tokenInput: '',
      error: null
    };
  },
  methods: {
    viewResults() {
      this.error = null;
      
      // Extract token from input (either direct token or URL)
      const token = this.extractToken(this.tokenInput.trim());
      
      if (!token) {
        this.error = '请输入有效的TOKEN或邀请链接';
        return;
      }
      
      // Navigate to results page with the extracted token
      this.$router.push(`/results/${token}`);
    },
    
    extractToken(input) {
      // If it's already a token (alphanumeric string), return it directly
      if (/^[a-zA-Z0-9-]+$/.test(input)) {
        return input;
      }
      
      // Try to extract token from URL
      try {
        const url = new URL(input);
        const pathParts = url.pathname.split('/');
        
        // Check if it's a user-b URL
        if (pathParts.length >= 3 && pathParts[pathParts.length - 2] === 'user-b') {
          return pathParts[pathParts.length - 1];
        }
        
        // Check if it's a results URL
        if (pathParts.length >= 3 && pathParts[pathParts.length - 2] === 'results') {
          return pathParts[pathParts.length - 1];
        }
        
        // If we have a token parameter
        const tokenParam = url.searchParams.get('token');
        if (tokenParam) {
          return tokenParam;
        }
      } catch (e) {
        // Not a valid URL, continue to return null
      }
      
      // If we can't extract a token, return null
      return null;
    }
  }
};
</script>