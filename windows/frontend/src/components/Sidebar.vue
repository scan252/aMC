<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()

interface NavItem {
  to: string
  label: string
  icon: string
}

interface NavGroup {
  caption: string
  items: NavItem[]
}

// 图标使用 1.6px 圆头描边的 24×24 网格，与 design/01 一致
const groups: NavGroup[] = [
  {
    caption: '工作台',
    items: [
      {
        to: '/',
        label: '总览',
        icon: 'M4 20V10M10 20V4M16 20v-7M21 20H3',
      },
      {
        to: '/gacha',
        label: '唤取分析',
        icon: 'M3 3v18h18M7 15l4-6 4 3 5-8',
      },
      {
        to: '/account',
        label: '账号中心',
        icon: 'M12 3a4 4 0 1 1 0 8 4 4 0 0 1 0-8ZM4 21c1.5-4 4.5-6 8-6s6.5 2 8 6',
      },
    ],
  },
  {
    caption: '图鉴与工具',
    items: [
      {
        to: '/codex',
        label: '角色图鉴',
        icon: 'M4 19V5a2 2 0 0 1 2-2h13v16H6a2 2 0 0 0-2 2zm0 0a2 2 0 0 0 2 2h13',
      },
      {
        to: '/echo',
        label: '声骸工具',
        icon: 'M12 3l8 4.5v9L12 21l-8-4.5v-9L12 3zm0 9l8-4.5M12 12v9M12 12L4 7.5',
      },
      {
        to: '/planner',
        label: '养成计算',
        icon: 'M5 3h14a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1zM8 7h8M8 11h3m2 0h3M8 15h3m2 0h3',
      },
    ],
  },
  {
    caption: '资讯',
    items: [
      {
        to: '/news',
        label: '资讯日历',
        icon: 'M4 7h16a1 1 0 0 1 1 1v11a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V8a1 1 0 0 1 1-1zm4-4v4m8-4v4M3 11h18',
      },
      {
        to: '/codes',
        label: '兑换码',
        icon: 'M5 9h14a1 1 0 0 1 1 1v10a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V10a1 1 0 0 1 1-1zm7 1v11M4 9h16l-1.6-4.2a1 1 0 0 0-.94-.65H6.54a1 1 0 0 0-.94.65L4 9zm4.5-5.5C9.5 3 10 4 9.6 5.2L9.4 6M14.5 3.5C14.5 3 14 4 14.4 5.2l.2.8',
      },
    ],
  },
]
</script>

<template>
  <aside class="sidebar glass">
    <div class="logo">
      <b class="gold-text">aMC</b>
      <span>SUITE</span>
    </div>
    <nav>
      <template v-for="group in groups" :key="group.caption">
        <div class="nav-cap">{{ group.caption }}</div>
        <router-link
          v-for="item in group.items"
          :key="item.to"
          :to="item.to"
          class="nav-item"
          :class="{ active: route.path === item.to }"
        >
          <svg viewBox="0 0 24 24"><path :d="item.icon" /></svg>
          <span>{{ item.label }}</span>
        </router-link>
      </template>
    </nav>
    <div class="side-foot">
      <router-link to="/settings" class="nav-item" :class="{ active: route.path === '/settings' }">
        <svg viewBox="0 0 24 24">
          <path d="M4 8h9m4 0h3M4 16h3m4 0h9" />
          <circle cx="15.5" cy="8" r="2.4" />
          <circle cx="9" cy="16" r="2.4" />
        </svg>
        <span>设置</span>
      </router-link>
      <div class="ver">v0.1.0 · Windows</div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  padding: 20px 14px 16px;
  min-height: 0;
}
.logo {
  display: flex;
  align-items: baseline;
  gap: 8px;
  padding: 2px 10px 20px;
}
.logo b {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.5px;
}
.logo span {
  font-size: 10px;
  letter-spacing: 0.34em;
  color: var(--txt-3);
  font-weight: 600;
}
nav {
  display: flex;
  flex-direction: column;
  gap: 3px;
  flex: 1;
  overflow-y: auto;
}
.nav-cap {
  font-size: 10px;
  letter-spacing: 0.3em;
  color: var(--txt-3);
  padding: 10px 12px 8px;
  font-weight: 600;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 11px;
  height: 40px;
  padding: 0 12px;
  border-radius: var(--radius-sm);
  color: var(--txt-2);
  font-size: 13.5px;
  text-decoration: none;
  transition: all 0.22s ease;
  position: relative;
  flex-shrink: 0;
}
.nav-item svg {
  width: 17px;
  height: 17px;
  stroke: currentColor;
  fill: none;
  stroke-width: 1.7;
  stroke-linecap: round;
  stroke-linejoin: round;
  opacity: 0.8;
}
.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: var(--txt-1);
}
.nav-item.active {
  background: linear-gradient(120deg, rgba(244, 206, 106, 0.16), rgba(222, 158, 59, 0.07));
  color: var(--warn);
  box-shadow: inset 0 0 0 1px rgba(244, 206, 106, 0.18);
}
.nav-item.active::before {
  content: '';
  position: absolute;
  left: -14px;
  top: 11px;
  bottom: 11px;
  width: 3px;
  border-radius: 3px;
  background: linear-gradient(var(--gold-1), var(--gold-2));
  box-shadow: 0 0 12px rgba(244, 206, 106, 0.8);
}
.side-foot {
  padding-top: 12px;
  border-top: 1px solid var(--line-2);
}
.ver {
  font-size: 11px;
  color: var(--txt-3);
  padding: 10px 12px 0;
  letter-spacing: 0.04em;
}
</style>
