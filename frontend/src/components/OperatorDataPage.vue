<template>
    <div class="operator-page">
        <section class="hero-card panel">
            <div>
                <p class="section-kicker">Operator Cache</p>
                <h2>公开招募干员数据</h2>
                <p class="summary">
                    打开页面时优先读取本地缓存；手动刷新时会重新抓取 wiki 页面，并将结构化数据和图片写回本地。
                </p>
            </div>

            <div class="hero-actions">
                <button class="primary-button" type="button" :disabled="loading" @click="handleFetch">
                    {{ loading ? '刷新中…' : '刷新远程数据' }}
                </button>
                <p class="meta-text">当前记录数: {{ operators.length }}</p>
                <p class="meta-text">当前来源: {{ sourceModeLabel }}</p>
                <p class="meta-text">抓取时间: {{ fetchedAtLabel }}</p>
                <p class="meta-text">来源地址: {{ sourceUrl }}</p>
            </div>
        </section>

        <el-alert v-if="errorMessage" class="status-alert" type="error" :closable="false" show-icon>
            {{ errorMessage }}
        </el-alert>

        <section v-if="loading && !hasLoaded" class="panel empty-state">
            <p class="section-kicker">Loading</p>
            <h3>正在加载本地缓存或远程数据</h3>
            <p>页面会先尝试读取本地缓存，只有手动刷新时才请求远程 wiki 页面。</p>
        </section>

        <section v-else-if="!hasLoaded" class="panel empty-state">
            <p class="section-kicker">Empty</p>
            <h3>当前还没有本地缓存数据</h3>
            <p>点击“刷新远程数据”后，会从 wiki 抓取干员数据并把数据和图片保存到本地。</p>
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
                                <th>顺序</th>
                                <th>图片</th>
                                <th>名称</th>
                                <th>星级</th>
                                <th>标签</th>
                                <th>职业 / 阵营</th>
                                <th>公开招募</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="operator in operators" :key="`${operator.order}-${operator.name}`">
                                <td>{{ operator.order + 1 }}</td>
                                <td>
                                    <img
                                        v-if="operatorImage(operator)"
                                        class="operator-image"
                                        :src="operatorImage(operator)"
                                        :alt="operator.name"
                                    >
                                    <div v-else class="image-fallback">无图</div>
                                </td>
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
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { FetchOperatorData, GetCachedOperatorImage, LoadCachedOperatorData } from '../../wailsjs/go/main/App'

const operators = ref([])
const operatorImageSources = ref({})
const loading = ref(false)
const hasLoaded = ref(false)
const errorMessage = ref('')
const sourceUrl = ref('https://wiki.biligame.com/arknights/公开招募工具')
const fetchedAt = ref('')
const fromCache = ref(false)
let imageRequestToken = 0

const publicRecruitableCount = computed(() => operators.value.filter((item) => item.isPublicRecruitable).length)
const highestRarityLabel = computed(() => {
    if (!operators.value.length) {
        return '-'
    }
    return `${Math.max(...operators.value.map((item) => item.rarity))} 星`
})
const sourceModeLabel = computed(() => (fromCache.value ? '本地缓存' : '远程刷新结果'))
const fetchedAtLabel = computed(() => (fetchedAt.value ? new Date(fetchedAt.value).toLocaleString() : '-'))

async function applyResult(result) {
    operators.value = Array.isArray(result.operators) ? result.operators : []
    sourceUrl.value = result.sourceUrl || sourceUrl.value
    fetchedAt.value = result.fetchedAt || ''
    fromCache.value = Boolean(result.fromCache)
    hasLoaded.value = Boolean(result.cacheAvailable || operators.value.length)
    await loadOperatorImages(operators.value)
}

function operatorImage(operator) {
    if (!operator?.localImagePath) {
        return ''
    }
    return operatorImageSources.value[operator.localImagePath] || ''
}

async function loadOperatorImages(records) {
    const currentToken = ++imageRequestToken
    operatorImageSources.value = {}

    const uniquePaths = [...new Set(records.map((item) => item.localImagePath).filter(Boolean))]
    if (!uniquePaths.length) {
        return
    }

    const resolvedEntries = await Promise.all(uniquePaths.map(async (relativePath) => {
        try {
            const result = await GetCachedOperatorImage(relativePath)
            if (!result?.found || !result.dataBase64) {
                return [relativePath, '']
            }
            const mimeType = result.mimeType || 'image/jpeg'
            return [relativePath, `data:${mimeType};base64,${result.dataBase64}`]
        } catch (error) {
            console.error(`加载干员图片失败: ${relativePath}`, error)
            return [relativePath, '']
        }
    }))

    if (currentToken !== imageRequestToken) {
        return
    }

    operatorImageSources.value = Object.fromEntries(
        resolvedEntries.filter(([, imageSource]) => Boolean(imageSource)),
    )
}

async function loadCache() {
    loading.value = true
    errorMessage.value = ''

    try {
        const result = await LoadCachedOperatorData()
        if (result.cacheAvailable) {
            await applyResult(result)
        } else {
            hasLoaded.value = false
            operatorImageSources.value = {}
        }
    } catch (error) {
        console.error('加载本地缓存失败:', error)
        errorMessage.value = typeof error === 'string' ? error : error?.message || '加载本地缓存失败'
    } finally {
        loading.value = false
    }
}

async function handleFetch() {
    loading.value = true
    errorMessage.value = ''

    try {
        const result = await FetchOperatorData()
        await applyResult(result)
        ElMessage.success(`已刷新 ${operators.value.length} 条干员数据，本地缓存已更新`)
    } catch (error) {
        console.error('刷新远程数据失败:', error)
        errorMessage.value = typeof error === 'string' ? error : error?.message || '刷新远程数据失败'
        ElMessage.error(errorMessage.value)
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    loadCache()
})
</script>

<style scoped>
.operator-page {
    display: grid;
    gap: 1.25rem;
}

.panel {
    border-radius: 1.5rem;
    background: rgba(255, 255, 255, 0.82);
    box-shadow:
        0 18px 52px rgba(101, 157, 212, 0.14),
        inset 0 1px 0 rgba(255, 255, 255, 0.92);
    padding: 1.25rem;
    backdrop-filter: blur(14px);
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
    color: #5f8fbf;
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
    color: #55728f;
}

.hero-actions {
    display: grid;
    gap: 0.6rem;
    justify-items: end;
    min-width: 17rem;
}

.primary-button {
    border: 0;
    border-radius: 999px;
    padding: 0.85rem 1.2rem;
    background: linear-gradient(135deg, #5ba9ff 0%, #8fd0ff 100%);
    color: #0d2a44;
    font: inherit;
    font-weight: 700;
    cursor: pointer;
}

.primary-button:disabled {
    cursor: wait;
    opacity: 0.72;
}

.meta-text {
    text-align: right;
    color: #6d89a4;
    word-break: break-all;
}

.status-alert {
    border-radius: 1rem;
}

.empty-state {
    text-align: center;
    color: #55728f;
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
    background: rgba(91, 169, 255, 0.08);
    display: grid;
    gap: 0.35rem;
}

.stat-label {
    color: #6d89a4;
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
    border-bottom: 1px solid rgba(91, 169, 255, 0.12);
    text-align: left;
    vertical-align: top;
}

th {
    color: #6d89a4;
    font-weight: 700;
}

.operator-image,
.image-fallback {
    width: 52px;
    height: 52px;
    border-radius: 0.9rem;
}

.operator-image {
    display: block;
    object-fit: cover;
    background: rgba(91, 169, 255, 0.08);
}

.image-fallback {
    display: grid;
    place-items: center;
    background: rgba(91, 169, 255, 0.1);
    color: #6d89a4;
    font-size: 0.82rem;
}

.name-cell {
    display: grid;
    gap: 0.18rem;
}

.name {
    font-weight: 700;
}

.name-meta {
    color: #7a97b4;
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
    background: rgba(91, 169, 255, 0.14);
    color: #2d628f;
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
