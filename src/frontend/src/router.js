import Vue from 'vue'
import Router from 'vue-router'

import Welcome from '@/views/Welcome.vue'
import Control from '@/views/Control.vue'
import NotImplemented from '@/views/NotImplemented.vue'



Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Welcome',
      component: Welcome,
    },
    {
      path: '/control',
      name: 'Control',
      component: Control,
    },
    {
      path: '/planner',
      name: 'Planner',
      component: NotImplemented,
    },
    {
      path: '/settings',
      name: 'Settings',
      component: NotImplemented,
    }
  ]
})