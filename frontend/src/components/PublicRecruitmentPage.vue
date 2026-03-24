<template>
    <div class="recruitment-page">
        <section class="hero-card panel">
            <div>
                <p class="section-kicker">Public Recruitment</p>
                <h2>公开招募标签组合</h2>
                <p class="summary">
                    基于本地缓存的公开招募干员数据筛选标签组合。你也可以开启模板驱动的自动识别，让页面按照识别设置中的模板持续同步当前招募标签。
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
                    <div class="selection-actions">
                        <label class="toggle-chip">
                            <input v-model="autoRecognitionEnabled" type="checkbox" @change="handleAutoRecognitionToggle">
                            <span>自动识别</span>
                        </label>
                        <button
                            v-if="selectedTags.length"
                            class="secondary-button"
                            type="button"
                            @click="clearSelectedTags"
                        >
                            清空选择
                        </button>
                    </div>
                </div>

                <p class="selection-note">
                    组合结果会显示所有非空子组合，没有命中干员的组合会自动隐藏。自动识别成功时会直接替换当前选中标签；识别失败时不会改动当前选择。
                </p>

                <div class="automation-panel">
                    <div>
                        <p class="section-kicker">Recognition</p>
                        <h4>{{ autoRecognitionEnabled ? '正在轮询识别模板' : '自动识别已关闭' }}</h4>
                        <p class="automation-copy">模板: {{ autoRecognitionTemplateLabel }}</p>
                        <p class="automation-copy">状态: {{ autoRecognitionStatusMessage }}</p>
                    </div>
                    <div class="automation-meta">
                        <span class="result-count">{{ autoRecognitionBusy ? '识别中' : '空闲' }}</span>
                        <span v-if="autoRecognitionLastSuccessAt" class="result-count">上次同步 {{ autoRecognitionLastSuccessAt }}</span>
                    </div>
                </div>

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
                <p>可以减少标签数量或更换标签，查看其他有效组合。</p>
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
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
    LoadCachedOperatorData,
    LoadRecognitionTemplates,
    RunPublicRecruitmentRecognition,
} from '../../wailsjs/go/main/App'
import { recruitmentTagGroups, recruitmentTagSet } from '../constants/recruitmentTags'

const selectedTags = ref([])
const recruitableOperators = ref([])
const loading = ref(false)
const errorMessage = ref('')
const hasRecruitmentData = ref(false)
const fromCache = ref(false)
const autoRecognitionEnabled = ref(false)
const autoRecognitionBusy = ref(false)
const autoRecognitionTemplate = ref(null)
const autoRecognitionStatusMessage = ref('等待开启自动识别')
const autoRecognitionLastSuccessAt = ref('')

let autoRecognitionTimer = null
let autoRecognitionSession = 0

const sourceModeLabel = computed(() => (fromCache.value ? '本地缓存' : '未加载'))
const autoRecognitionTemplateLabel = computed(() => {
    if (!autoRecognitionTemplate.value) {
        return '未找到识别模板'
    }

    const title = autoRecognitionTemplate.value.title || '(无标题)'
    const className = autoRecognitionTemplate.value.className || '(未知类名)'
    return `${title} / ${className}`
})

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

function setSelectedTags(tags) {
    selectedTags.value = [...tags]
}

function clearSelectedTags() {
    setSelectedTags([])
}

function toggleTag(tag) {
    if (isSelected(tag)) {
        setSelectedTags(selectedTags.value.filter((item) => item !== tag))
        return
    }

    if (selectedTags.value.length >= 5) {
        ElMessage.warning('最多只能同时选择 5 个标签')
        return
    }

    setSelectedTags([...selectedTags.value, tag])
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

async function loadRecognitionTemplateSummary() {
    try {
        const templates = await LoadRecognitionTemplates()
        autoRecognitionTemplate.value = Array.isArray(templates) && templates.length ? templates[0] : null
        if (!autoRecognitionTemplate.value && autoRecognitionEnabled.value) {
            autoRecognitionStatusMessage.value = '未找到可用识别模板'
        }
    } catch (error) {
        console.error('加载识别模板失败:', error)
        autoRecognitionTemplate.value = null
        autoRecognitionStatusMessage.value = '识别模板加载失败'
    }
}

function clearAutoRecognitionTimer() {
    if (autoRecognitionTimer) {
        window.clearTimeout(autoRecognitionTimer)
        autoRecognitionTimer = null
    }
}

function stopAutoRecognitionLoop() {
    autoRecognitionSession += 1
    clearAutoRecognitionTimer()
    autoRecognitionBusy.value = false
}

function scheduleNextAutoRecognitionRun(sessionId) {
    clearAutoRecognitionTimer()
    autoRecognitionTimer = window.setTimeout(() => {
        void runAutoRecognition(sessionId)
    }, 500)
}

function describeRecognitionFailure(result) {
    switch (result?.failureReason) {
    case 'no_template':
        return '未配置识别模板'
    case 'no_window':
        return '目标窗口不可用'
    case 'ambiguous_window':
        return '存在多个同标题同类名窗口，无法唯一定位'
    case 'capture_failed':
        return '窗口截图失败'
    case 'incomplete_match':
        return result?.failureMessage || '存在未命中或多命中的区域，本轮不更新标签'
    default:
        return result?.failureMessage || '识别未返回可用标签'
    }
}

async function runAutoRecognition(sessionId) {
    if (!autoRecognitionEnabled.value || sessionId !== autoRecognitionSession) {
        return
    }
    if (!autoRecognitionTemplate.value?.id) {
        autoRecognitionStatusMessage.value = '未找到可用识别模板'
        return
    }

    autoRecognitionBusy.value = true
    try {
        const result = await RunPublicRecruitmentRecognition({ templateId: autoRecognitionTemplate.value.id })
        if (!autoRecognitionEnabled.value || sessionId !== autoRecognitionSession) {
            return
        }

        if (result?.success) {
            setSelectedTags(Array.isArray(result.recognizedTags) ? result.recognizedTags : [])
            autoRecognitionLastSuccessAt.value = new Date().toLocaleTimeString()
            autoRecognitionStatusMessage.value = `已识别: ${(result.recognizedTags || []).join('、')}`
        } else {
            autoRecognitionStatusMessage.value = describeRecognitionFailure(result)
        }
    } catch (error) {
        console.error('公开招募自动识别失败:', error)
        if (!autoRecognitionEnabled.value || sessionId !== autoRecognitionSession) {
            return
        }
        autoRecognitionStatusMessage.value = typeof error === 'string' ? error : error?.message || '识别请求失败'
    } finally {
        if (!autoRecognitionEnabled.value || sessionId !== autoRecognitionSession) {
            autoRecognitionBusy.value = false
            return
        }
        autoRecognitionBusy.value = false
        scheduleNextAutoRecognitionRun(sessionId)
    }
}

async function handleAutoRecognitionToggle() {
    if (!autoRecognitionEnabled.value) {
        stopAutoRecognitionLoop()
        autoRecognitionStatusMessage.value = '自动识别已关闭'
        return
    }

    stopAutoRecognitionLoop()
    autoRecognitionSession += 1
    await loadRecognitionTemplateSummary()
    if (!autoRecognitionTemplate.value?.id) {
        autoRecognitionStatusMessage.value = '未找到可用识别模板'
        return
    }

    autoRecognitionStatusMessage.value = '正在启动自动识别'
    void runAutoRecognition(autoRecognitionSession)
}

onMounted(() => {
    void loadRecruitmentData()
    void loadRecognitionTemplateSummary()
})

onBeforeUnmount(() => {
    stopAutoRecognitionLoop()
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
.operator-meta,
.automation-copy {
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
.operator-meta,
.automation-copy {
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

.selection-actions,
.automation-meta {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
}

.automation-panel {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 1rem;
    border-radius: 1.25rem;
    background: rgba(91, 169, 255, 0.08);
    padding: 1rem;
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
.secondary-button,
.toggle-chip {
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

.secondary-button,
.toggle-chip {
    padding: 0.7rem 1rem;
    background: rgba(91, 169, 255, 0.14);
    color: #24527c;
}

.toggle-chip {
    display: inline-flex;
    align-items: center;
    gap: 0.55rem;
    font-weight: 600;
}

.toggle-chip input {
    accent-color: #2d628f;
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
