<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900 text-white p-4">
    <div class="max-w-4xl w-full text-center">
      <h1 class="text-5xl md:text-7xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-pink-400 via-purple-300 to-blue-400">
        赛博问缘
      </h1>
      <p class="text-xl md:text-2xl mb-10 max-w-2xl mx-auto leading-relaxed">
        探索数字世界中的心灵契合，通过15个深度问题，发现你与TA的奇妙连接。
      </p>
      
      <div class="flex flex-col sm:flex-row justify-center gap-6 mb-16">
        <div class="bg-black bg-opacity-30 backdrop-blur-lg rounded-xl p-6 border border-white border-opacity-20 shadow-2xl">
          <h2 class="text-2xl font-semibold mb-3">如何参与</h2>
          <ul class="text-left space-y-2 min-w-[280px]" >
            <li class="flex items-start">
              <span class="text-green-400 mr-2">✓</span>
              <span>发起人回答多轮问题并生成邀请链接</span>
            </li>
            <li class="flex items-start">
              <span class="text-green-400 mr-2">✓</span>
              <span>将链接分享给同伴</span>
            </li>
            <li class="flex items-start">
              <span class="text-green-400 mr-2">✓</span>
              <span>同伴回答相同问题</span>
            </li>
            <li class="flex items-start">
              <span class="text-green-400 mr-2">✓</span>
              <span>AI分析匹配度并生成结果</span>
            </li>
          </ul>
        </div>
        
        <div class="bg-black bg-opacity-30 backdrop-blur-lg rounded-xl p-6 border border-white border-opacity-20 shadow-2xl">
          <h2 class="text-2xl font-semibold mb-3">特色功能</h2>
          <ul class="text-left space-y-2 min-w-[280px]">
            <li class="flex items-start">
              <span class="text-purple-400 mr-2">✦</span>
              <span>AI驱动的匹配算法</span>
            </li>
            <li class="flex items-start">
              <span class="text-purple-400 mr-2">✦</span>
              <span>深度心理学问题</span>
            </li>
            <li class="flex items-start">
              <span class="text-purple-400 mr-2">✦</span>
              <span>隐私保护设计</span>
            </li>
            <li class="flex items-start">
              <span class="text-purple-400 mr-2">✦</span>
              <span>赛博朋克主题界面</span>
            </li>
          </ul>
        </div>
      </div>
      
      <button 
        @click="startQuiz"
        class="px-10 py-4 bg-gradient-to-r from-pink-500 to-purple-600 text-white font-bold text-xl rounded-full shadow-lg hover:from-pink-600 hover:to-purple-700 transform hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 focus:ring-purple-400 focus:ring-opacity-50"
      >
        开始探索
      </button>
      <button 
        @click="viewResults"
        class="mt-6 px-10 py-4 bg-gradient-to-r from-blue-500 to-cyan-600 text-white font-bold text-xl rounded-full shadow-lg hover:from-blue-600 hover:to-cyan-700 transform hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none focus:ring-4 focus:ring-cyan-400 focus:ring-opacity-50"
      >
        查看结果
      </button>
      
      <div class="mt-16">
        <h3 class="text-2xl font-semibold mb-6">常见问题</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-3xl mx-auto text-left">
          <div class="bg-black bg-opacity-20 p-4 rounded-lg">
            <h4 class="font-bold text-lg mb-2">数据会公开吗？</h4>
            <p class="text-gray-300">问答数据可以自由选择是否公开，选择不与对方分享则仅自己可见。</p>
          </div>
          <div class="bg-black bg-opacity-20 p-4 rounded-lg">
            <h4 class="font-bold text-lg mb-2">匹配度准确吗？</h4>
            <p class="text-gray-300">我们的AI算法基于心理学研究，提供有意义的匹配度分析。</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Homepage',
  methods: {
    startQuiz() {
      this.$router.push('/user-a');
    },
    viewResults() {
      // For now, we'll navigate to a new route where users can input the token
      // In a later step, we'll implement the token input functionality
      this.$router.push('/view-results');
    }
  }
};
</script>

<style scoped>
/* Add any additional styles if needed */
</style>