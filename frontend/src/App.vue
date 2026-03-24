<template>
    <div class="app-shell">
        <aside class="sidebar">
            <div class="brand">
                <div class="brand-mark">RH</div>
                <div>
                    <p class="brand-title">Recruit Helper</p>
                    <p class="brand-subtitle">Arknights desktop toolkit</p>
                </div>
            </div>

            <nav class="nav-list" aria-label="Primary">
                <button
                    v-for="item in navigationItems"
                    :key="item.key"
                    class="nav-item"
                    :class="{ active: activePage === item.key }"
                    type="button"
                    @click="activePage = item.key"
                >
                    <span class="nav-badge">{{ item.shortLabel }}</span>
                    <span>{{ item.label }}</span>
                </button>
            </nav>
        </aside>

        <main class="content-panel">
            <header class="page-header">
                <div>
                    <p class="eyebrow">Main View</p>
                    <h1>{{ activePageMeta.label }}</h1>
                    <p class="description">{{ activePageMeta.description }}</p>
                </div>
            </header>

            <section class="page-body">
                <WindowToolsPage v-if="activePage === 'window-tools'" />
                <OperatorDataPage v-else />
            </section>
        </main>
    </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import OperatorDataPage from './components/OperatorDataPage.vue'
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
        description: '抓取并浏览公开招募干员基础数据，为后续推荐逻辑做准备。',
    },
]

const activePage = ref('window-tools')

const activePageMeta = computed(() =>
    navigationItems.find((item) => item.key === activePage.value) ?? navigationItems[0],
)
</script>

<style scoped>
.app-shell {
    min-height: 100vh;
    display: grid;
    grid-template-columns: 280px 1fr;
    background:
        radial-gradient(circle at top left, rgba(184, 87, 65, 0.18), transparent 28%),
        radial-gradient(circle at bottom right, rgba(39, 72, 102, 0.22), transparent 26%),
        linear-gradient(135deg, #f3ede4 0%, #efe4d2 100%);
    color: #1e2127;
}

.sidebar {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    padding: 1.5rem;
    background: rgba(26, 34, 43, 0.92);
    color: #f7f3ea;
    box-shadow: inset -1px 0 0 rgba(255, 255, 255, 0.08);
}

.brand {
    display: flex;
    align-items: center;
    gap: 0.9rem;
}

.brand-mark {
    width: 3rem;
    height: 3rem;
    display: grid;
    place-items: center;
    border-radius: 0.9rem;
    background: linear-gradient(135deg, #c96b4b 0%, #f3c08b 100%);
    color: #1f2023;
    font-weight: 800;
    letter-spacing: 0.08em;
}

.brand-title,
.brand-subtitle,
.eyebrow,
.description {
    margin: 0;
}

.brand-title {
    font-size: 1.1rem;
    font-weight: 700;
}

.brand-subtitle {
    margin-top: 0.25rem;
    font-size: 0.88rem;
    color: rgba(247, 243, 234, 0.72);
}

.nav-list {
    display: flex;
    flex-direction: column;
    gap: 0.65rem;
}

.nav-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    width: 100%;
    padding: 0.95rem 1rem;
    border: 0;
    border-radius: 1rem;
    background: rgba(255, 255, 255, 0.06);
    color: inherit;
    font: inherit;
    cursor: pointer;
    text-align: left;
    transition:
        transform 0.18s ease,
        background-color 0.18s ease,
        box-shadow 0.18s ease;
}

.nav-item:hover {
    transform: translateX(3px);
    background: rgba(255, 255, 255, 0.12);
}

.nav-item.active {
    background: linear-gradient(135deg, rgba(201, 107, 75, 0.94), rgba(243, 192, 139, 0.9));
    color: #151719;
    box-shadow: 0 18px 30px rgba(24, 28, 35, 0.24);
}

.nav-badge {
    min-width: 2rem;
    height: 2rem;
    display: inline-grid;
    place-items: center;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.15);
    font-size: 0.76rem;
    font-weight: 800;
    letter-spacing: 0.06em;
}

.content-panel {
    min-width: 0;
    display: flex;
    flex-direction: column;
    padding: 1.5rem;
}

.page-header {
    padding: 0.25rem 0 1.4rem;
}

.eyebrow {
    text-transform: uppercase;
    letter-spacing: 0.16em;
    font-size: 0.72rem;
    color: #8a6858;
}

h1 {
    margin: 0.35rem 0;
    font-size: clamp(2rem, 4vw, 3rem);
    line-height: 1.02;
}

.description {
    max-width: 44rem;
    color: #5f5148;
}

.page-body {
    min-height: 0;
    flex: 1;
}

@media (max-width: 960px) {
    .app-shell {
        grid-template-columns: 1fr;
    }

    .sidebar {
        position: sticky;
        top: 0;
        z-index: 10;
    }

    .nav-list {
        flex-direction: row;
        flex-wrap: wrap;
    }

    .nav-item {
        flex: 1 1 220px;
    }
}
</style>
