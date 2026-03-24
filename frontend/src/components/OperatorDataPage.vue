<template>
    <div class="operator-page">
        <section class="hero-card">
            <div>
                <p class="section-kicker">Recruitment Dataset</p>
                <h2>公开招募干员数据</h2>
                <p class="summary">
                    从 biligame wiki 的公开招募工具页面抓取 `.contentDetail` 条目，解析干员名、星级、标签和原始元数据。
                </p>
            </div>

            <div class="hero-actions">
                <button class="fetch-button" type="button" :disabled="loading" @click="handleFetch">
                    {{ loading ? '获取中…' : '获取干员数据' }}
                </button>
                <p class="meta-text">当前记录数: {{ operators.length }}</p>
                <p class="meta-text">来源: {{ sourceUrl }}</p>
            </div>
        </section>

        <el-alert v-if="errorMessage" class="status-alert" type="error" :closable="false" show-icon>
            {{ errorMessage }}
        </el-alert>

        <section v-if="!hasLoaded && !loading" class="empty-state panel">
            <p class="section-kicker">Empty</p>
            <h3>当前还没有干员数据</h3>
            <p>点击“获取干员数据”后，页面会从公开招募工具抓取并解析全部干员条目。</p>
        </section>

        <section v-else-if="loading && !operators.length" class="empty-state panel">
            <p class="section-kicker">Loading</p>
            <h3>正在抓取页面并解析干员数据</h3>
            <p>这一步会请求 wiki 页面，再把 HTML 中的 `.contentDetail` 条目转换为结构化记录。</p>
        </section>

        <section v-else class="data-layout">
            <article class="panel stats-panel">
                <p class="section-kicker">Overview</p>
                <h3>数据概览</h3>
                <div class="stats-grid">
                    <div class="stat-card">
                        <span class="stat-label">干员数量</span>
                        <strong>{{ operators.length }}</strong>
                    </div>
                    <div class="stat-card">
                        <span class="stat-label">公开招募可用</span>
                        <strong>{{ publicRecruitableCount }}</strong>
                    </div>
                    <div class="stat-card">
                        <span class="stat-label">最高星级</span>
                        <strong>{{ highestRarityLabel }}</strong>
                    </div>
                </div>
            </article>

            <article class="panel table-panel">
                <div class="table-header">
                    <div>
                        <p class="section-kicker">Results</p>
                        <h3>干员列表</h3>
                    </div>
                </div>

                <div class="table-scroll">
                    <table>
                        <thead>
                            <tr>
                                <th>名称</th>
                                <th>星级</th>
                                <th>标签</th>
                                <th>职业 / 阵营</th>
                                <th>公开招募</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="operator in operators" :key="operator.name">
                                <td>
                                    <div class="name-cell">
                                        <span class="name">{{ operator.name }}</span>
                                        <span class="name-meta">{{ operator.metadata.gender || '未知性别' }}</span>
                                    </div>
                                </td>
                                <td>{{ operator.rarity }} 星</td>
                                <td>
                                    <div class="tag-list">
                                        <span
                                            v-for="tag in operator.displayTags"
                                            :key="`${operator.name}-${tag}`"
                                            class="tag"
                                        >
                                            {{ tag }}
                                        </span>
                                    </div>
                                </td>
                                <td>{{ operator.metadata.profession || '未知职业' }} / {{ operator.metadata.origin || '未知阵营' }}</td>
                                <td>{{ operator.isPublicRecruitable ? '是' : '否' }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </article>
        </section>
    </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { FetchOperatorData } from '../../wailsjs/go/main/App'

const operators = ref([])
const loading = ref(false)
const hasLoaded = ref(false)
const errorMessage = ref('')
const sourceUrl = ref('https://wiki.biligame.com/arknights/公开招募工具')

const publicRecruitableCount = computed(() => operators.value.filter((item) => item.isPublicRecruitable).length)
const highestRarityLabel = computed(() => {
    if (!operators.value.length) {
        return '-'
    }
    return `${Math.max(...operators.value.map((item) => item.rarity))} 星`
})

async function handleFetch() {
    loading.value = true
    errorMessage.value = ''

    try {
        const result = await FetchOperatorData()
        operators.value = result.operators ?? []
        sourceUrl.value = result.sourceUrl || sourceUrl.value
        hasLoaded.value = true
        ElMessage.success(`已获取 ${operators.value.length} 条干员数据`)
    } catch (error) {
        console.error('获取干员数据失败:', error)
        errorMessage.value = typeof error === 'string' ? error : error?.message || '获取干员数据失败'
        ElMessage.error(errorMessage.value)
    } finally {
        loading.value = false
    }
}
</script>

<style scoped>
.operator-page {
    display: grid;
    gap: 1.25rem;
}

.hero-card,
.panel {
    border-radius: 1.5rem;
    background: rgba(255, 251, 245, 0.84);
    box-shadow:
        0 18px 60px rgba(95, 68, 51, 0.12),
        inset 0 1px 0 rgba(255, 255, 255, 0.8);
    padding: 1.25rem;
    backdrop-filter: blur(18px);
}

.hero-card {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
}

.section-kicker,
.summary,
.meta-text {
    margin: 0;
}

.section-kicker {
    color: #a47460;
    font-size: 0.76rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
}

h2,
h3 {
    margin: 0.3rem 0 0.45rem;
}

.summary {
    max-width: 44rem;
    color: #5f5148;
}

.hero-actions {
    display: grid;
    gap: 0.6rem;
    justify-items: end;
    min-width: 17rem;
}

.fetch-button {
    border: 0;
    border-radius: 999px;
    padding: 0.85rem 1.2rem;
    background: linear-gradient(135deg, #c96b4b 0%, #efb37d 100%);
    color: #17191d;
    font: inherit;
    font-weight: 700;
    cursor: pointer;
}

.fetch-button:disabled {
    cursor: wait;
    opacity: 0.72;
}

.meta-text {
    text-align: right;
    color: #7d685d;
    word-break: break-all;
}

.status-alert {
    border-radius: 1rem;
}

.empty-state {
    text-align: center;
    color: #5f5148;
}

.data-layout {
    display: grid;
    gap: 1.25rem;
}

.stats-grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 0.9rem;
}

.stat-card {
    padding: 1rem;
    border-radius: 1rem;
    background: rgba(99, 83, 73, 0.08);
    display: grid;
    gap: 0.35rem;
}

.stat-label {
    color: #7d685d;
    font-size: 0.9rem;
}

.table-scroll {
    overflow: auto;
}

table {
    width: 100%;
    border-collapse: collapse;
}

th,
td {
    padding: 0.9rem 0.75rem;
    border-bottom: 1px solid rgba(30, 33, 39, 0.08);
    text-align: left;
    vertical-align: top;
}

th {
    color: #7d685d;
    font-weight: 700;
}

.name-cell {
    display: grid;
    gap: 0.18rem;
}

.name {
    font-weight: 700;
}

.name-meta {
    color: #8f7c70;
    font-size: 0.88rem;
}

.tag-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.45rem;
}

.tag {
    display: inline-flex;
    align-items: center;
    border-radius: 999px;
    padding: 0.3rem 0.65rem;
    background: rgba(201, 107, 75, 0.14);
    color: #744833;
    font-size: 0.86rem;
}

@media (max-width: 900px) {
    .hero-card {
        flex-direction: column;
    }

    .hero-actions {
        justify-items: start;
        min-width: 0;
    }

    .meta-text {
        text-align: left;
    }

    .stats-grid {
        grid-template-columns: 1fr;
    }
}
</style>
