<template>
    <div class="recruitment-page">
        <section class="hero-card panel">
            <div>
                <p class="section-kicker">Public Recruitment</p>
                <h2>公开招募标签组合</h2>
                <p class="summary">
                    基于本地缓存的公开招募干员数据筛选标签组合。最多可选择 5 个标签，并查看每种有效组合命中的干员。
                </p>
            </div>

            <div class="hero-meta">
                <p class="meta-text">已选标签: {{ selectedTags.length }}/5</p>
                <p class="meta-text">可招募干员: {{ recruitableOperators.length }}</p>
                <p class="meta-text">有效组合: {{ resultGroups.length }}</p>
                <p class="meta-text">数据来源: {{ sourceModeLabel }}</p>
            </div>
        </section>

        <el-alert v-if="errorMessage" class="status-alert" type="error" :closable="false" show-icon>
            {{ errorMessage }}
        </el-alert>

        <section v-if="loading" class="panel empty-state">
            <p class="section-kicker">Loading</p>
            <h3>正在读取本地干员数据</h3>
            <p>公开招募页面会直接复用已缓存的干员数据，不会额外发起新的抓取请求。</p>
        </section>

        <section v-else-if="!hasRecruitmentData" class="panel empty-state">
            <p class="section-kicker">Need Data</p>
            <h3>当前还没有可用的公开招募干员数据</h3>
            <p>请先前往“干员数据”页面刷新或加载本地缓存，然后再回到这里进行标签筛选。</p>
        </section>

        <template v-else>
            <section class="panel selection-panel">
                <div class="selection-header">
                    <div>
                        <p class="section-kicker">Selection</p>
                        <h3>选择公开招募标签</h3>
                    </div>
                    <button
                        v-if="selectedTags.length"
                        class="secondary-button"
                        type="button"
                        @click="clearSelectedTags"
                    >
                        清空选择
                    </button>
                </div>

                <p class="selection-note">
                    组合结果会显示所有非空子组合；没有命中干员的组合会自动隐藏。
                </p>

                <div class="selected-tags" v-if="selectedTags.length">
                    <span v-for="tag in selectedTags" :key="tag" class="selected-tag-chip">{{ tag }}</span>
                </div>

                <div class="tag-group-list">
                    <article v-for="group in recruitmentTagGroups" :key="group.key" class="tag-group-card">
                        <h4>{{ group.label }}</h4>
                        <div class="tag-chip-list">
                            <button
                                v-for="tag in group.tags"
                                :key="tag"
                                class="tag-chip"
                                :class="{ active: isSelected(tag) }"
                                type="button"
                                @click="toggleTag(tag)"
                            >
                                {{ tag }}
                            </button>
                        </div>
                    </article>
                </div>
            </section>

            <section class="panel empty-state" v-if="!selectedTags.length">
                <p class="section-kicker">Combinations</p>
                <h3>先选择一个或多个标签</h3>
                <p>选择标签后会立即显示全部有效组合及其匹配干员。</p>
            </section>

            <section class="panel empty-state" v-else-if="!resultGroups.length">
                <p class="section-kicker">No Matches</p>
                <h3>当前标签组合没有匹配干员</h3>
                <p>可以减少标签数量或更换标签，查看其它有效组合。</p>
            </section>

            <section v-else class="results-layout">
                <article v-for="group in resultGroups" :key="group.key" class="panel result-group">
                    <div class="result-header">
                        <div>
                            <p class="section-kicker">Combination</p>
                            <h3>{{ group.label }}</h3>
                        </div>
                        <span class="result-count">{{ group.operators.length }} 位干员</span>
                    </div>

                    <div class="operator-card-grid">
                        <article
                            v-for="operator in group.operators"
                            :key="`${group.key}-${operator.order}-${operator.name}`"
                            class="operator-card"
                        >
                            <div class="operator-card-top">
                                <strong>{{ operator.name }}</strong>
                                <span>{{ operator.rarity }} 星</span>
                            </div>
                            <p class="operator-meta">
                                {{ operator.metadata.profession || '未知职业' }} / {{ operator.metadata.origin || '未知阵营' }}
                            </p>
                            <div class="operator-tags">
                                <span
                                    v-for="tag in operator.recruitmentTagList"
                                    :key="`${group.key}-${operator.name}-${tag}`"
                                    class="operator-tag"
                                >
                                    {{ tag }}
                                </span>
                            </div>
                        </article>
                    </div>
                </article>
            </section>
        </template>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { LoadCachedOperatorData } from '../../wailsjs/go/main/App'
import { recruitmentTagGroups, recruitmentTagSet } from '../constants/recruitmentTags'

const selectedTags = ref([])
const recruitableOperators = ref([])
const loading = ref(false)
const errorMessage = ref('')
const hasRecruitmentData = ref(false)
const fromCache = ref(false)

const sourceModeLabel = computed(() => (fromCache.value ? '本地缓存' : '未加载'))

const resultGroups = computed(() => {
    const combinations = buildTagCombinations(selectedTags.value)
    return combinations
        .map((tags) => {
            const operators = recruitableOperators.value.filter((operator) =>
                tags.every((tag) => operator.recruitmentTagSet.has(tag)),
            )
            return {
                key: tags.join('|'),
                label: tags.join(' + '),
                tags,
                operators,
            }
        })
        .filter((group) => group.operators.length > 0)
})

function isSelected(tag) {
    return selectedTags.value.includes(tag)
}

function clearSelectedTags() {
    selectedTags.value = []
}

function toggleTag(tag) {
    if (isSelected(tag)) {
        selectedTags.value = selectedTags.value.filter((item) => item !== tag)
        return
    }

    if (selectedTags.value.length >= 5) {
        ElMessage.warning('最多只能同时选择 5 个标签')
        return
    }

    selectedTags.value = [...selectedTags.value, tag]
}

function buildTagCombinations(tags) {
    const combinations = []

    for (let size = tags.length; size >= 1; size -= 1) {
        collectCombinations(tags, size, 0, [], combinations)
    }

    return combinations
}

function collectCombinations(sourceTags, targetSize, startIndex, current, combinations) {
    if (current.length === targetSize) {
        combinations.push([...current])
        return
    }

    for (let index = startIndex; index < sourceTags.length; index += 1) {
        current.push(sourceTags[index])
        collectCombinations(sourceTags, targetSize, index + 1, current, combinations)
        current.pop()
    }
}

function normalizeOperator(operator) {
    if (!operator?.isPublicRecruitable) {
        return null
    }

    const tagValues = new Set()
    const pushTag = (value) => {
        if (value && recruitmentTagSet.has(value)) {
            tagValues.add(value)
        }
    }

    pushTag(operator.metadata?.profession)
    ;(operator.displayTags || []).forEach(pushTag)
    ;(operator.metadata?.recruitmentTags || []).forEach(pushTag)
    ;(operator.metadata?.seniorityTags || []).forEach(pushTag)
    ;(operator.metadata?.raw || []).forEach(pushTag)
    ;(operator.metadata?.extra || []).forEach(pushTag)

    if (!tagValues.size) {
        return null
    }

    return {
        ...operator,
        recruitmentTagSet: tagValues,
        recruitmentTagList: Array.from(tagValues),
    }
}

async function loadRecruitmentData() {
    loading.value = true
    errorMessage.value = ''

    try {
        const result = await LoadCachedOperatorData()
        const operators = Array.isArray(result.operators) ? result.operators : []
        const normalized = operators
            .map(normalizeOperator)
            .filter(Boolean)

        recruitableOperators.value = normalized
        fromCache.value = Boolean(result.fromCache)
        hasRecruitmentData.value = Boolean(result.cacheAvailable && normalized.length)
    } catch (error) {
        console.error('加载公开招募数据失败:', error)
        errorMessage.value = typeof error === 'string' ? error : error?.message || '加载公开招募数据失败'
        recruitableOperators.value = []
        hasRecruitmentData.value = false
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    loadRecruitmentData()
})
</script>

<style scoped>
.recruitment-page {
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

.hero-card,
.selection-header,
.result-header,
.operator-card-top {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
}

.hero-meta {
    display: grid;
    gap: 0.45rem;
    justify-items: end;
    min-width: 14rem;
}

.section-kicker,
.summary,
.meta-text,
.selection-note,
.operator-meta {
    margin: 0;
}

.section-kicker {
    color: #5f8fbf;
    font-size: 0.76rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
}

h2,
h3,
h4 {
    margin: 0.3rem 0 0.45rem;
}

.summary,
.selection-note,
.operator-meta {
    color: #55728f;
}

.meta-text {
    text-align: right;
    color: #6d89a4;
}

.status-alert {
    border-radius: 1rem;
}

.empty-state {
    text-align: center;
    color: #55728f;
}

.selection-panel,
.results-layout,
.tag-group-list,
.operator-card-grid,
.operator-tags,
.tag-chip-list,
.selected-tags {
    display: grid;
    gap: 1rem;
}

.selection-panel {
    gap: 1.2rem;
}

.results-layout {
    gap: 1.25rem;
}

.tag-group-list {
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
}

.tag-group-card {
    border-radius: 1.25rem;
    background: rgba(91, 169, 255, 0.08);
    padding: 1rem;
}

.tag-chip-list,
.selected-tags,
.operator-tags {
    display: flex;
    flex-wrap: wrap;
}

.tag-chip,
.selected-tag-chip,
.operator-tag,
.result-count,
.secondary-button {
    border-radius: 999px;
    font: inherit;
}

.tag-chip,
.secondary-button {
    border: 0;
    cursor: pointer;
}

.tag-chip {
    padding: 0.45rem 0.8rem;
    background: rgba(255, 255, 255, 0.75);
    color: #24527c;
    transition: transform 0.18s ease, background-color 0.18s ease;
}

.tag-chip:hover {
    transform: translateY(-1px);
    background: rgba(91, 169, 255, 0.18);
}

.tag-chip.active {
    background: linear-gradient(135deg, #5ba9ff 0%, #8fd0ff 100%);
    color: #0d2a44;
    font-weight: 700;
}

.selected-tag-chip,
.operator-tag,
.result-count {
    padding: 0.35rem 0.75rem;
    background: rgba(91, 169, 255, 0.14);
    color: #2d628f;
}

.secondary-button {
    padding: 0.7rem 1rem;
    background: rgba(91, 169, 255, 0.14);
    color: #24527c;
}

.result-group {
    display: grid;
    gap: 1rem;
}

.operator-card-grid {
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
}

.operator-card {
    border-radius: 1.2rem;
    padding: 1rem;
    background: rgba(91, 169, 255, 0.08);
    display: grid;
    gap: 0.6rem;
}

@media (max-width: 900px) {
    .hero-card,
    .selection-header,
    .result-header {
        flex-direction: column;
    }

    .hero-meta {
        justify-items: start;
    }

    .meta-text {
        text-align: left;
    }
}
</style>
