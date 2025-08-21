<template>
  <div class="max-w-4xl mx-auto p-4">
    <h1 class="text-3xl font-bold mb-6 text-center bg-clip-text text-transparent bg-gradient-to-r from-pink-400 via-purple-300 to-blue-400">赛博问答 - 受邀人</h1>
    <p class="mb-8 text-center text-xl">请回答以下15个问题。</p>
    
    <div v-if="!token" class="p-8 rounded-xl bg-black bg-opacity-30 backdrop-blur-lg border border-white border-opacity-20 shadow-lg text-center">
      <p class="mb-6 text-xl">无效的邀请链接或会话数据不存在。</p>
      <router-link
        to="/user-a"
        class="inline-block py-3 px-6 bg-gradient-to-r from-pink-500 to-purple-600 text-white font-bold text-lg rounded-xl shadow-lg hover:from-pink-600 hover:to-purple-700 transform hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 focus:ring-purple-400 focus:ring-opacity-50"
      >
        返回发起人页面
      </router-link>
    </div>
    
    <form v-else @submit.prevent="submitAnswers" class="mb-8">
      <Question
        v-for="question in questions"
        :key="question.id"
        :question="question"
        :answer="answers[question.question] || ''"
        @update-answer="updateAnswer"
      />
      
      <div class="mb-6 p-6 rounded-xl bg-black bg-opacity-30 backdrop-blur-lg border border-white border-opacity-20 shadow-lg">
        <label class="flex items-center cursor-pointer">
          <input
            type="checkbox"
            v-model="shareAnswers"
            class="h-5 w-5 text-purple-500 border-gray-300 rounded focus:ring-purple-500"
          />
          <span class="ml-3 block text-lg font-medium">
            是否分享您的答案给对方？
          </span>
        </label>
      </div>
      
      <button
        type="submit"
        :disabled="isSubmitting"
        class="w-full py-4 px-6 bg-gradient-to-r from-pink-500 to-purple-600 text-white font-bold text-lg rounded-xl shadow-lg hover:from-pink-600 hover:to-purple-700 transform hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 focus:ring-purple-400 focus:ring-opacity-50 disabled:opacity-50"
      >
        {{ isSubmitting ? '提交中...' : '提交答案' }}
      </button>
    </form>
    
    <div v-if="error" class="p-8 rounded-xl bg-red-900 bg-opacity-30 backdrop-blur-lg border border-red-500 border-opacity-20 shadow-lg">
      <h2 class="text-2xl font-semibold mb-6 text-center">错误</h2>
      <p class="mb-6 text-center">{{ error }}</p>
    </div>
  </div>
</template>

<script>
import { loadQuestions, submitUserB } from '@/services/questionService';
import Question from '@/components/Question.vue';

export default {
  name: 'UserB',
  components: {
    Question
  },
  data() {
    return {
      questions: [],
      answers: {},
      shareAnswers: false,
      token: null,
      isSubmitting: false,
      error: null
    };
  },
  async created() {
    this.token = this.$route.params.token;
    this.questions = await loadQuestions();
  },
  methods: {
    updateAnswer({ questionId, answer }) {
      // Find the question object by ID
      const question = this.questions.find(q => q.id === questionId);
      if (question) {
        // Use the question content as the key
        this.answers[question.question] = answer;
      }
    },
    async submitAnswers() {
      this.isSubmitting = true;
      this.error = null;
      
      try {
        // Submit answers to the backend
        await submitUserB(this.token, this.answers, this.shareAnswers);
        
        // Redirect to the results page
        this.$router.push(`/results/${this.token}`);
      } catch (err) {
        console.error('Failed to submit answers:', err);
        this.error = '提交答案失败，请稍后重试。';
      } finally {
        this.isSubmitting = false;
      }
    }
  }
};
</script>