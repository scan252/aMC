import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../pages/Dashboard.vue'
import GachaAnalysis from '../pages/GachaAnalysis.vue'
import Account from '../pages/Account.vue'
import Codex from '../pages/Codex.vue'
import EchoTool from '../pages/EchoTool.vue'
import Planner from '../pages/Planner.vue'
import News from '../pages/News.vue'
import Codes from '../pages/Codes.vue'
import Settings from '../pages/Settings.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'dashboard', component: Dashboard, meta: { title: '总览' } },
    { path: '/gacha', name: 'gacha', component: GachaAnalysis, meta: { title: '唤取分析' } },
    { path: '/account', name: 'account', component: Account, meta: { title: '账号中心' } },
    { path: '/codex', name: 'codex', component: Codex, meta: { title: '角色图鉴' } },
    { path: '/echo', name: 'echo', component: EchoTool, meta: { title: '声骸工具' } },
    { path: '/planner', name: 'planner', component: Planner, meta: { title: '养成计算' } },
    { path: '/news', name: 'news', component: News, meta: { title: '资讯日历' } },
    { path: '/codes', name: 'codes', component: Codes, meta: { title: '兑换码' } },
    { path: '/settings', name: 'settings', component: Settings, meta: { title: '设置' } },
  ],
})

router.afterEach((to) => {
  document.title = to.meta.title ? `aMC Suite · ${to.meta.title}` : 'aMC Suite'
})

export default router
