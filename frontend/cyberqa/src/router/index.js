import { createRouter, createWebHistory } from 'vue-router';
import UserA from '@/views/UserA.vue';
import UserB from '@/views/UserB.vue';
import Results from '@/views/Results.vue';
import Homepage from '@/views/Homepage.vue';

const routes = [
  {
    path: '/',
    name: 'Homepage',
    component: Homepage
  },
  {
    path: '/user-a',
    name: 'UserA',
    component: UserA
  },
  {
    path: '/user-b/:token',
    name: 'UserB',
    component: UserB,
    props: true
  },
  {
    path: '/results/:token',
    name: 'Results',
    component: Results,
    props: true
  },
  {
    path: '/view-results',
    name: 'ViewResults',
    component: () => import('@/views/ViewResults.vue')
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;