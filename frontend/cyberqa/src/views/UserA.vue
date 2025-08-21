<template>
  <div class="max-w-4xl mx-auto p-4">
    <h1 class="text-3xl font-bold mb-6 text-center bg-clip-text text-transparent bg-gradient-to-r from-pink-400 via-purple-300 to-blue-400">赛博问缘 - 发起人</h1>
    <p class="mb-8 text-center text-xl">请回答以下15个问题，完成后将生成一个邀请链接发送给你的同伴。</p>
    
    <form @submit.prevent="submitAnswers" class="mb-8">
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
        {{ isSubmitting ? '提交中...' : '提交并生成邀请链接' }}
      </button>
    </form>
    
    <div v-if="invitationLink" class="p-8 rounded-xl bg-black bg-opacity-30 backdrop-blur-lg border border-white border-opacity-20 shadow-lg">
      <h2 class="text-2xl font-semibold mb-6 text-center">邀请链接已生成</h2>
      <p class="mb-6 text-center">请将以下链接发送给你的同伴：</p>
      <a
        :href="invitationLink"
        target="_blank"
        class="block mb-6 p-4 bg-black bg-opacity-20 border border-white border-opacity-10 rounded-lg break-words text-blue-300 hover:text-blue-100 hover:underline"
      >
        {{ invitationLink }}
      </a>
      <button
        @click="copyLink"
        class="w-full py-3 px-6 bg-gradient-to-r from-green-500 to-teal-600 text-white font-bold text-lg rounded-xl shadow-lg hover:from-green-600 hover:to-teal-700 transform hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 focus:ring-green-400 focus:ring-opacity-50"
      >
        复制链接
      </button>
    </div>
    
    <div v-if="error" class="p-8 rounded-xl bg-red-900 bg-opacity-30 backdrop-blur-lg border border-red-500 border-opacity-20 shadow-lg">
      <h2 class="text-2xl font-semibold mb-6 text-center">错误</h2>
      <p class="mb-6 text-center">{{ error }}</p>
    </div>
  </div>
</template>

<script>
import { loadQuestions, submitUserA } from '@/services/questionService';
import Question from '@/components/Question.vue';

export default {
  name: 'UserA',
  components: {
    Question
  },
  data() {
    return {
      questions: [],
      answers: {},
      shareAnswers: false,
      invitationLink: null,
      isSubmitting: false,
      error: null
    };
  },
  async created() {
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
        const token = await submitUserA(this.answers, this.shareAnswers);
        
        // Generate invitation link
        this.invitationLink = `${window.location.origin}/user-b/${token}`;
      } catch (err) {
        console.error('Failed to submit answers:', err);
        this.error = '提交答案失败，请稍后重试。';
      } finally {
        this.isSubmitting = false;
      }
    },
    copyLink() {
      navigator.clipboard.writeText(this.invitationLink).then(() => {
        alert('链接已复制到剪贴板');
      }).catch(err => {
        console.error('Failed to copy link: ', err);
        // Fallback for older browsers
        const textArea = document.createElement('textarea');
        textArea.value = this.invitationLink;
        document.body.appendChild(textArea);
        textArea.select();
        try {
          document.execCommand('copy');
          alert('链接已复制到剪贴板');
        } catch (err) {
          console.error('Fallback: Oops, unable to copy', err);
          alert('无法复制链接，请手动复制');
        }
        document.body.removeChild(textArea);
      });
    }
  }
};
</script>