import Vue from 'vue'
import Router from 'vue-router'
import LogView from '@/components/LogView'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'LogView',
      component: LogView
    }
  ]
})
