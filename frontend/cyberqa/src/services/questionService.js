// Base URL for the API
const API_BASE_URL = '/api';

/**
 * Load questions from the backend API.
 * @returns {Promise<Array>} A promise that resolves to an array of questions.
 */
export async function loadQuestions() {
  const response = await fetch(`${API_BASE_URL}/questions`);
  
  if (!response.ok) {
    throw new Error('Failed to load questions');
  }
  
  return await response.json();
}

/**
 * Submit User A's answers to the backend.
 * @param {Object} answers - The answers from User A.
 * @param {boolean} shareAnswers - Whether User A wants to share their answers.
 * @returns {Promise<string>} A promise that resolves to the invitation token.
 */
export async function submitUserA(answers, shareAnswers) {
  const response = await fetch(`${API_BASE_URL}/submit-user-a`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      answers,
      shareAnswers,
    }),
  });

  if (!response.ok) {
    throw new Error('Failed to submit User A answers');
  }

  const data = await response.json();
  return data.token;
}

/**
 * Submit User B's answers to the backend.
 * @param {string} token - The invitation token.
 * @param {Object} answers - The answers from User B.
 * @param {boolean} shareAnswers - Whether User B wants to share their answers.
 * @returns {Promise<boolean>} A promise that resolves to true if successful.
 */
export async function submitUserB(token, answers, shareAnswers) {
  const response = await fetch(`${API_BASE_URL}/submit-user-b`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      token,
      answers,
      shareAnswers,
    }),
  });

  if (!response.ok) {
    throw new Error('Failed to submit User B answers');
  }

  const data = await response.json();
  return data.success;
}

/**
 * Get results from the backend.
 * @param {string} token - The invitation token.
 * @returns {Promise<Object>} A promise that resolves to the results.
 */
export async function getResults(token) {
  const response = await fetch(`${API_BASE_URL}/results/${token}`);

  if (!response.ok) {
    throw new Error('Failed to fetch results');
  }

  return await response.json();
}