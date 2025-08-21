<template>
  <div class="max-w-4xl mx-auto p-4">
    <h1 class="text-3xl font-bold mb-6 text-center bg-clip-text text-transparent bg-gradient-to-r from-pink-400 via-purple-300 to-blue-400">赛博问缘 - 结果</h1>
    
    <div v-if="loading" class="p-8 rounded-xl bg-black bg-opacity-30 backdrop-blur-lg border border-white border-opacity-20 shadow-lg text-center">
      <p class="mb-6 text-xl">加载结果中...</p>
    </div>
    
    <div v-else-if="error" class="p-8 rounded-xl bg-red-900 bg-opacity-30 backdrop-blur-lg border border-red-500 border-opacity-20 shadow-lg text-center">
      <p class="mb-6 text-xl">{{ error }}</p>
      <router-link
        to="/user-a"
        class="inline-block py-3 px-6 bg-gradient-to-r from-pink-500 to-purple-600 text-white font-bold text-lg rounded-xl shadow-lg hover:from-pink-600 hover:to-purple-700 transform hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 focus:ring-purple-400 focus:ring-opacity-50"
      >
        返回发起人页面
      </router-link>
    </div>
    
    <div v-else-if="!sessionData" class="p-8 rounded-xl bg-black bg-opacity-30 backdrop-blur-lg border border-white border-opacity-20 shadow-lg text-center">
      <p class="mb-6 text-xl">无效的链接或会话数据不存在。</p>
      <router-link
        to="/user-a"
        class="inline-block py-3 px-6 bg-gradient-to-r from-pink-500 to-purple-600 text-white font-bold text-lg rounded-xl shadow-lg hover:from-pink-600 hover:to-purple-700 transform hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 focus:ring-purple-400 focus:ring-opacity-50"
      >
        返回受邀人页面
      </router-link>
    </div>
    
    <div v-else class="mt-6">
      <div class="p-8 rounded-xl bg-black bg-opacity-30 backdrop-blur-lg border border-white border-opacity-20 shadow-lg text-center">
        <h2 class="text-2xl font-semibold mb-4">匹配度: {{ sessionData.compatibility }}%</h2>
        <p class="text-lg">{{ sessionData.summary }}</p>
      </div>
      
      <div v-if="sessionData.userAShared || sessionData.userBShared" class="mt-8">
        <h3 class="text-xl font-semibold mb-6 pb-2 text-center">答案详情</h3>
        
        <div v-if="sessionData.userAShared" class="mb-8">
          <h4 class="text-lg font-semibold mb-4 text-pink-300">发起人的答案</h4>
          <div v-for="(answer, question) in sessionData.userAAnswers" :key="question" class="mb-4 p-4 rounded-lg bg-black bg-opacity-20 border border-white border-opacity-10">
            <p class="font-medium mb-2 text-lg">{{ question }}</p>
            <p class="text-gray-200">{{ answer }}</p>
          </div>
        </div>
        
        <div v-if="sessionData.userBShared" class="mb-8">
          <h4 class="text-lg font-semibold mb-4 text-purple-300">受邀人的答案</h4>
          <div v-for="(answer, question) in sessionData.userBAnswers" :key="question" class="mb-4 p-4 rounded-lg bg-black bg-opacity-20 border border-white border-opacity-10">
            <p class="font-medium mb-2 text-lg">{{ question }}</p>
            <p class="text-gray-200">{{ answer }}</p>
          </div>
        </div>
      </div>
      
      <div v-else class="mt-8 p-6 rounded-xl bg-black bg-opacity-30 backdrop-blur-lg border border-white border-opacity-20 shadow-lg text-center italic text-gray-300">
        <p>双方都选择不公开答案。</p>
      </div>
      
      <router-link
        to="/user-a"
        class="inline-block mt-6 py-3 px-6 bg-gradient-to-r from-pink-500 to-purple-600 text-white font-bold text-lg rounded-xl shadow-lg hover:from-pink-600 hover:to-purple-700 transform hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 focus:ring-purple-400 focus:ring-opacity-50"
      >
        重新开始
      </router-link>
    </div>
  </div>
</template>

<script>
import { loadQuestions, getResults } from '@/services/questionService';

export default {
  name: 'Results',
  data() {
    return {
      questions: [],
      token: null,
      sessionData: null,
      loading: true,
      error: null
    };
  },
  async created() {
    this.token = this.$route.params.token;
    this.questions = await loadQuestions();
    
    // Fetch results from the backend
    await this.fetchResults();
  },
  methods: {
    async fetchResults() {
      this.loading = true;
      this.error = null;
      
      try {
        const results = await getResults(this.token);
        this.sessionData = results;
      } catch (err) {
        console.error('Failed to fetch results:', err);
        this.error = '获取结果失败，请稍后重试。';
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>