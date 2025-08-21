<template>
  <div class="mb-6 p-6 rounded-xl bg-black bg-opacity-30 backdrop-blur-lg border border-white border-opacity-20 shadow-lg">
    <h3 class="text-xl font-semibold mb-4">{{ question.question }}</h3>
    <div v-if="question.isMultipleChoice" class="space-y-3">
      <div v-for="option in question.options" :key="option" class="flex items-center">
        <input
          type="radio"
          :name="`question-${question.id}`"
          :value="option"
          v-model="localAnswer"
          @change="onAnswerChange"
          class="h-5 w-5 text-purple-500 border-gray-300 focus:ring-purple-500"
        />
        <label class="ml-3 block text-lg font-medium">
          {{ option }}
        </label>
      </div>
    </div>
    <div v-else class="w-full">
      <textarea
        :value="localAnswer"
        @input="onTextInputChange"
        class="w-full px-4 py-3 bg-black bg-opacity-20 border border-white border-opacity-10 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 text-white placeholder-gray-400"
        rows="4"
        placeholder="请输入您的答案..."
      ></textarea>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Question',
  props: {
    question: {
      type: Object,
      required: true
    },
    answer: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      localAnswer: this.answer
    };
  },
  watch: {
    answer(newVal) {
      this.localAnswer = newVal;
    }
  },
  methods: {
    onAnswerChange(event) {
      this.$emit('update-answer', {
        questionId: this.question.id,
        answer: event.target.value
      });
    },
    onTextInputChange(event) {
      this.localAnswer = event.target.value;
      this.$emit('update-answer', {
        questionId: this.question.id,
        answer: event.target.value
      });
    }
  }
};
</script>