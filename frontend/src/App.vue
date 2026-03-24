<template>
    <div class="app-shell">
        <el-drawer v-model="drawerOpen" direction="ltr" size="280px" class="app-drawer">
            <template #header>
                <div class="drawer-header">
                    <div class="brand-mark">RH</div>
                    <div>
                        <p class="brand-title">公开招募助手</p>
                        <p class="brand-subtitle">桌面工作台</p>
                    </div>
                </div>
            </template>

            <nav class="drawer-nav" aria-label="Primary navigation">
                <button
                    v-for="item in navigationItems"
                    :key="item.key"
                    class="drawer-nav-item"
                    :class="{ active: activePage === item.key }"
                    type="button"
                    @click="selectPage(item.key)"
                >
                    <span class="nav-badge">{{ item.shortLabel }}</span>
                    <span>{{ item.label }}</span>
                </button>
            </nav>
        </el-drawer>

        <main class="shell-main">
            <header class="topbar">
                <button class="menu-trigger" type="button" @click="drawerOpen = true">
                    <span class="menu-trigger-label">MENU</span>
                    <span>{{ activePageMeta.label }}</span>
                </button>

                <div class="topbar-copy">
                    <p class="eyebrow">Main View</p>
                    <h1>{{ activePageMeta.label }}</h1>
                    <p class="description">{{ activePageMeta.description }}</p>
                </div>
            </header>

            <section class="page-body">
                <WindowToolsPage v-if="activePage === 'window-tools'" />
                <OperatorDataPage v-else-if="activePage === 'operator-data'" />
                <PublicRecruitmentPage v-else />
            </section>
        </main>
    </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import OperatorDataPage from './components/OperatorDataPage.vue'
import PublicRecruitmentPage from './components/PublicRecruitmentPage.vue'
import WindowToolsPage from './components/WindowToolsPage.vue'

const navigationItems = [
    {
        key: 'window-tools',
        label: '窗口工具',
        shortLabel: 'WT',
        description: '保留现有 Win32 窗口枚举、选择和高亮能力。',
    },
    {
        key: 'operator-data',
        label: '干员数据',
        shortLabel: 'OD',
        description: '优先加载本地缓存的公开招募干员数据，并支持手动刷新。',
    },
    {
        key: 'public-recruitment',
        label: '公开招募',
        shortLabel: 'PR',
        description: '选择公开招募标签并查看全部有效组合对应的干员结果。',
    },
]

const drawerOpen = ref(false)
const activePage = ref('window-tools')

const activePageMeta = computed(() =>
    navigationItems.find((item) => item.key === activePage.value) ?? navigationItems[0],
)

function selectPage(pageKey) {
    activePage.value = pageKey
    drawerOpen.value = false
}
</script>

<style scoped>
.app-shell {
    min-height: 100vh;
    background:
        radial-gradient(circle at top left, rgba(143, 205, 255, 0.55), transparent 28%),
        radial-gradient(circle at bottom right, rgba(210, 235, 255, 0.8), transparent 34%),
        linear-gradient(180deg, #f8fcff 0%, #edf6ff 100%);
    color: #17324d;
}

:deep(.app-drawer .el-drawer) {
    background: linear-gradient(180deg, #f5fbff 0%, #ebf5ff 100%);
}

:deep(.app-drawer .el-drawer__header) {
    margin-bottom: 1rem;
    padding-bottom: 0;
}

.shell-main {
    padding: 1.5rem;
}

.topbar {
    display: flex;
    align-items: flex-start;
    gap: 1.25rem;
    margin-bottom: 1.5rem;
}

.menu-trigger {
    display: inline-flex;
    align-items: center;
    gap: 0.9rem;
    min-width: 14rem;
    border: 0;
    border-radius: 999px;
    padding: 0.95rem 1.2rem;
    background: linear-gradient(135deg, #5ba9ff 0%, #8fd0ff 100%);
    color: #0d2a44;
    font: inherit;
    font-weight: 700;
    cursor: pointer;
    box-shadow: 0 14px 32px rgba(91, 169, 255, 0.24);
}

.menu-trigger-label {
    display: inline-grid;
    place-items: center;
    min-width: 3rem;
    height: 2rem;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.58);
    font-size: 0.76rem;
    letter-spacing: 0.08em;
}

.drawer-header {
    display: flex;
    align-items: center;
    gap: 0.85rem;
}

.brand-mark {
    width: 3rem;
    height: 3rem;
    display: grid;
    place-items: center;
    border-radius: 1rem;
    background: linear-gradient(135deg, #5ba9ff 0%, #a5deff 100%);
    color: #0d2a44;
    font-weight: 800;
    letter-spacing: 0.06em;
}

.brand-title,
.brand-subtitle,
.eyebrow,
.description {
    margin: 0;
}

.brand-title {
    font-size: 1.08rem;
    font-weight: 700;
}

.brand-subtitle {
    margin-top: 0.15rem;
    color: #6a86a0;
    font-size: 0.9rem;
}

.drawer-nav {
    display: grid;
    gap: 0.75rem;
}

.drawer-nav-item {
    display: flex;
    align-items: center;
    gap: 0.8rem;
    width: 100%;
    border: 0;
    border-radius: 1rem;
    padding: 0.95rem 1rem;
    background: rgba(91, 169, 255, 0.08);
    color: #17324d;
    font: inherit;
    cursor: pointer;
    text-align: left;
    transition: background-color 0.18s ease, transform 0.18s ease;
}

.drawer-nav-item:hover {
    transform: translateX(2px);
    background: rgba(91, 169, 255, 0.16);
}

.drawer-nav-item.active {
    background: linear-gradient(135deg, rgba(91, 169, 255, 0.96), rgba(171, 225, 255, 0.9));
    color: #0d2a44;
}

.nav-badge {
    min-width: 2rem;
    height: 2rem;
    display: inline-grid;
    place-items: center;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.62);
    font-size: 0.76rem;
    font-weight: 800;
    letter-spacing: 0.06em;
}

.topbar-copy {
    padding-top: 0.2rem;
}

.eyebrow {
    text-transform: uppercase;
    letter-spacing: 0.16em;
    font-size: 0.72rem;
    color: #6f93b2;
}

h1 {
    margin: 0.35rem 0;
    font-size: clamp(2rem, 4vw, 3rem);
    line-height: 1.02;
}

.description {
    max-width: 42rem;
    color: #58728c;
}

.page-body {
    min-height: 0;
}

@media (max-width: 900px) {
    .topbar {
        flex-direction: column;
    }

    .menu-trigger {
        min-width: 0;
    }
}
</style>
