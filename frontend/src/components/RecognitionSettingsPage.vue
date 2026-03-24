<template>
    <div class="recognition-page">
        <section class="hero-card panel">
            <div>
                <p class="section-kicker">Recognition</p>
                <h2>识别设置</h2>
                <p class="summary">
                    选择目标窗口、截取窗口画面、框选多个识别区域并保存模板，用于后续判断指定区域是否匹配参考图片。
                </p>
            </div>

            <div class="hero-actions">
                <button class="primary-button" type="button" @click="toggleSelectionMode">
                    {{ selectionMode ? '停止指向选择' : '开始指向窗口' }}
                </button>
                <button
                    class="secondary-button"
                    type="button"
                    :disabled="!hoverWindow.hwnd"
                    @click="captureHoveredWindow"
                >
                    选中当前窗口
                </button>
                <button
                    class="secondary-button"
                    type="button"
                    :disabled="!capturedWindow.hwnd"
                    @click="saveTemplate"
                >
                    保存模板
                </button>
            </div>
        </section>

        <el-alert v-if="errorMessage" class="status-alert" type="error" :closable="false" show-icon>
            {{ errorMessage }}
        </el-alert>

        <section class="panel info-grid">
            <article class="info-card">
                <p class="section-kicker">Hover</p>
                <h3>鼠标下窗口</h3>
                <p>标题: {{ hoverWindow.title || '(无标题)' }}</p>
                <p>类名: {{ hoverWindow.className || '(未知)' }}</p>
                <p>句柄: {{ hoverWindow.hwnd || 0 }}</p>
                <p>坐标: {{ hoverWindow.x }}, {{ hoverWindow.y }}</p>
            </article>

            <article class="info-card">
                <p class="section-kicker">Selected</p>
                <h3>当前模板窗口</h3>
                <p>标题: {{ capturedWindow.title || '(未选择)' }}</p>
                <p>类名: {{ capturedWindow.className || '(未选择)' }}</p>
                <p>句柄: {{ capturedWindow.hwnd || 0 }}</p>
                <p>区域数: {{ regions.length }}</p>
            </article>
        </section>

        <section class="content-grid">
            <article class="panel preview-panel">
                <div class="panel-heading">
                    <div>
                        <p class="section-kicker">Preview</p>
                        <h3>窗口截图</h3>
                    </div>
                    <button
                        class="secondary-button"
                        type="button"
                        :disabled="!capturedWindow.hwnd"
                        @click="captureSelectedWindow"
                    >
                        重新截图
                    </button>
                </div>

                <div v-if="!capturedWindow.imageBase64" class="empty-state">
                    <p>请先指向并选中一个目标窗口。</p>
                </div>
                <div v-else class="preview-shell">
                    <div
                        ref="imageStageRef"
                        class="image-stage"
                        @mousedown="startRegionSelection"
                        @mousemove="updateRegionSelection"
                        @mouseup="finishRegionSelection"
                        @mouseleave="finishRegionSelection"
                    >
                        <img
                            ref="previewImageRef"
                            class="preview-image"
                            :src="previewImageSource"
                            alt="窗口截图"
                            draggable="false"
                        >
                        <div
                            v-for="region in regions"
                            :key="region.id"
                            class="region-box"
                            :style="regionStyle(region)"
                        >
                            <span class="region-label">{{ region.label || region.id }}</span>
                        </div>
                        <div
                            v-if="draftRegion"
                            class="region-box draft"
                            :style="regionStyle(draftRegion)"
                        />
                    </div>
                </div>
            </article>

            <article class="panel regions-panel">
                <div class="panel-heading">
                    <div>
                        <p class="section-kicker">Regions</p>
                        <h3>识别区域</h3>
                    </div>
                    <button v-if="regions.length" class="secondary-button" type="button" @click="clearRegions">
                        清空区域
                    </button>
                </div>

                <div v-if="!regions.length" class="empty-state">
                    <p>在截图上按住鼠标拖拽，即可新增一个识别区域。</p>
                </div>
                <div v-else class="region-list">
                    <article v-for="region in regions" :key="region.id" class="region-item">
                        <div class="region-item-top">
                            <strong>{{ region.id }}</strong>
                            <button class="danger-button" type="button" @click="removeRegion(region.id)">
                                删除
                            </button>
                        </div>
                        <label class="field-label">
                            标签
                            <input v-model="region.label" class="field-input" type="text" placeholder="例如：开始按钮">
                        </label>
                        <p class="region-meta">
                            相对位置: x={{ region.x.toFixed(4) }}, y={{ region.y.toFixed(4) }},
                            w={{ region.width.toFixed(4) }}, h={{ region.height.toFixed(4) }}
                        </p>
                    </article>
                </div>
            </article>
        </section>

        <section class="content-grid">
            <article class="panel templates-panel">
                <div class="panel-heading">
                    <div>
                        <p class="section-kicker">Templates</p>
                        <h3>已保存模板</h3>
                    </div>
                    <button class="secondary-button" type="button" @click="loadTemplates">刷新列表</button>
                </div>

                <div v-if="!templates.length" class="empty-state">
                    <p>当前还没有已保存的识别模板。</p>
                </div>
                <div v-else class="template-list">
                    <article
                        v-for="template in templates"
                        :key="template.id"
                        class="template-item"
                        :class="{ active: selectedTemplateId === template.id }"
                    >
                        <div>
                            <strong>{{ template.title || '(无标题)' }}</strong>
                            <p>{{ template.className || '(未知类名)' }}</p>
                            <p>{{ template.regionCount }} 个区域</p>
                        </div>
                        <div class="template-actions">
                            <button class="secondary-button" type="button" @click="openTemplate(template.id)">
                                加载
                            </button>
                            <button
                                class="secondary-button"
                                type="button"
                                :disabled="!capturedWindow.hwnd"
                                @click="matchTemplate(template.id)"
                            >
                                匹配当前窗口
                            </button>
                        </div>
                    </article>
                </div>
            </article>

            <article class="panel match-panel">
                <div class="panel-heading">
                    <div>
                        <p class="section-kicker">Matching</p>
                        <h3>区域匹配结果</h3>
                    </div>
                </div>

                <div v-if="!matchResults.length" class="empty-state">
                    <p>加载一个模板后，可以对当前选中的窗口执行区域匹配。</p>
                </div>
                <div v-else class="match-list">
                    <article v-for="result in matchResults" :key="result.regionId" class="match-item">
                        <strong>{{ result.label || result.regionId }}</strong>
                        <span :class="['match-state', result.match ? 'match' : 'mismatch']">
                            {{ result.match ? '匹配' : '不匹配' }}
                        </span>
                    </article>
                </div>
            </article>
        </section>
    </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
    CaptureWindowForRecognition,
    GetMousePosition,
    GetWindowInfo,
    GetWindowUnderMouse,
    GetRecognitionTemplate,
    LoadRecognitionTemplates,
    MatchRecognitionTemplate,
    SaveRecognitionTemplate,
} from '../../wailsjs/go/main/App'

const selectionMode = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const hoverWindow = ref({ hwnd: 0, title: '', className: '', x: 0, y: 0 })
const capturedWindow = ref({ hwnd: 0, title: '', className: '', width: 0, height: 0, imageBase64: '' })
const regions = ref([])
const templates = ref([])
const selectedTemplateId = ref('')
const matchResults = ref([])
const draftRegion = ref(null)
const imageStageRef = ref(null)
const previewImageRef = ref(null)
let selectionTimer = null
let dragStartPoint = null

const previewImageSource = computed(() =>
    capturedWindow.value.imageBase64 ? `data:image/png;base64,${capturedWindow.value.imageBase64}` : '',
)

function toggleSelectionMode() {
    selectionMode.value = !selectionMode.value
}

function clearRegions() {
    regions.value = []
}

function removeRegion(regionId) {
    regions.value = regions.value.filter((region) => region.id !== regionId)
}

function buildRegionId() {
    return `region-${String(regions.value.length + 1).padStart(2, '0')}`
}

function regionStyle(region) {
    return {
        left: `${region.x * 100}%`,
        top: `${region.y * 100}%`,
        width: `${region.width * 100}%`,
        height: `${region.height * 100}%`,
    }
}

function getStagePoint(event) {
    const stage = imageStageRef.value
    if (!stage) {
        return null
    }
    const rect = stage.getBoundingClientRect()
    if (!rect.width || !rect.height) {
        return null
    }
    const x = Math.min(Math.max(event.clientX - rect.left, 0), rect.width)
    const y = Math.min(Math.max(event.clientY - rect.top, 0), rect.height)
    return {
        x: x / rect.width,
        y: y / rect.height,
    }
}

function startRegionSelection(event) {
    if (!capturedWindow.value.imageBase64) {
        return
    }
    const point = getStagePoint(event)
    if (!point) {
        return
    }
    dragStartPoint = point
    draftRegion.value = {
        id: 'draft',
        label: '',
        x: point.x,
        y: point.y,
        width: 0,
        height: 0,
    }
}

function updateRegionSelection(event) {
    if (!dragStartPoint) {
        return
    }
    const point = getStagePoint(event)
    if (!point) {
        return
    }
    draftRegion.value = normalizeDraftRegion(dragStartPoint, point)
}

function finishRegionSelection(event) {
    if (!dragStartPoint || !draftRegion.value) {
        return
    }
    const point = getStagePoint(event)
    const finalRegion = point ? normalizeDraftRegion(dragStartPoint, point) : draftRegion.value
    dragStartPoint = null
    draftRegion.value = null

    if (finalRegion.width < 0.01 || finalRegion.height < 0.01) {
        return
    }

    regions.value = [
        ...regions.value,
        {
            id: buildRegionId(),
            label: `区域 ${regions.value.length + 1}`,
            x: finalRegion.x,
            y: finalRegion.y,
            width: finalRegion.width,
            height: finalRegion.height,
        },
    ]
}

function normalizeDraftRegion(start, end) {
    const x = Math.min(start.x, end.x)
    const y = Math.min(start.y, end.y)
    const width = Math.abs(start.x - end.x)
    const height = Math.abs(start.y - end.y)
    return { id: 'draft', label: '', x, y, width, height }
}

async function refreshHoveredWindow() {
    if (!selectionMode.value) {
        selectionTimer = setTimeout(refreshHoveredWindow, 200)
        return
    }

    try {
        const [hwnd, position] = await Promise.all([GetWindowUnderMouse(), GetMousePosition()])
        if (!hwnd) {
            hoverWindow.value = { hwnd: 0, title: '', className: '', x: position.x, y: position.y }
        } else {
            const info = await GetWindowInfo(hwnd)
            hoverWindow.value = {
                hwnd,
                title: info.title,
                className: info.className,
                x: position.x,
                y: position.y,
            }
        }
    } catch (error) {
        console.error('刷新鼠标下窗口失败:', error)
    } finally {
        selectionTimer = setTimeout(refreshHoveredWindow, 200)
    }
}

async function captureWindow(hwnd) {
    loading.value = true
    errorMessage.value = ''
    matchResults.value = []

    try {
        const result = await CaptureWindowForRecognition(hwnd)
        capturedWindow.value = {
            hwnd: result.hwnd,
            title: result.title,
            className: result.className,
            width: result.width,
            height: result.height,
            imageBase64: result.imageBase64,
        }
        regions.value = []
    } catch (error) {
        console.error('截取目标窗口失败:', error)
        errorMessage.value = typeof error === 'string' ? error : error?.message || '截取目标窗口失败'
    } finally {
        loading.value = false
    }
}

async function captureHoveredWindow() {
    if (!hoverWindow.value.hwnd) {
        ElMessage.info('请先将鼠标移动到目标窗口上')
        return
    }
    await captureWindow(hoverWindow.value.hwnd)
}

async function captureSelectedWindow() {
    if (!capturedWindow.value.hwnd) {
        return
    }
    await captureWindow(capturedWindow.value.hwnd)
}

async function saveTemplate() {
    if (!capturedWindow.value.hwnd || !capturedWindow.value.imageBase64) {
        ElMessage.info('请先选择并截取目标窗口')
        return
    }
    if (!regions.value.length) {
        ElMessage.info('请先在截图上划定至少一个区域')
        return
    }

    loading.value = true
    errorMessage.value = ''

    try {
        const saved = await SaveRecognitionTemplate({
            hwnd: capturedWindow.value.hwnd,
            title: capturedWindow.value.title,
            className: capturedWindow.value.className,
            screenshotPng: capturedWindow.value.imageBase64,
            width: capturedWindow.value.width,
            height: capturedWindow.value.height,
            regions: regions.value.map((region) => ({
                id: region.id,
                label: region.label,
                x: region.x,
                y: region.y,
                width: region.width,
                height: region.height,
            })),
        })
        selectedTemplateId.value = saved.id
        ElMessage.success(`已保存模板 ${saved.title || saved.id}`)
        await loadTemplates()
    } catch (error) {
        console.error('保存识别模板失败:', error)
        errorMessage.value = typeof error === 'string' ? error : error?.message || '保存识别模板失败'
        ElMessage.error(errorMessage.value)
    } finally {
        loading.value = false
    }
}

async function loadTemplates() {
    try {
        templates.value = await LoadRecognitionTemplates()
    } catch (error) {
        console.error('加载模板列表失败:', error)
    }
}

async function openTemplate(templateId) {
    loading.value = true
    errorMessage.value = ''
    try {
        const template = await GetRecognitionTemplate(templateId)
        selectedTemplateId.value = template.id
        capturedWindow.value = {
            hwnd: template.hwnd,
            title: template.title,
            className: template.className,
            width: template.width,
            height: template.height,
            imageBase64: template.screenshotPng,
        }
        regions.value = template.regions.map((region) => ({
            id: region.id,
            label: region.label,
            x: region.x,
            y: region.y,
            width: region.width,
            height: region.height,
        }))
        matchResults.value = []
    } catch (error) {
        console.error('加载识别模板失败:', error)
        errorMessage.value = typeof error === 'string' ? error : error?.message || '加载识别模板失败'
    } finally {
        loading.value = false
    }
}

async function matchTemplate(templateId) {
    if (!capturedWindow.value.hwnd) {
        ElMessage.info('请先选择一个当前窗口用于匹配')
        return
    }

    loading.value = true
    errorMessage.value = ''
    try {
        const result = await MatchRecognitionTemplate({
            templateId,
            hwnd: capturedWindow.value.hwnd,
        })
        selectedTemplateId.value = result.templateId
        matchResults.value = result.results
    } catch (error) {
        console.error('执行区域匹配失败:', error)
        errorMessage.value = typeof error === 'string' ? error : error?.message || '执行区域匹配失败'
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    refreshHoveredWindow()
    loadTemplates()
})

onUnmounted(() => {
    if (selectionTimer) {
        clearTimeout(selectionTimer)
    }
})
</script>

<style scoped>
.recognition-page {
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
.hero-actions,
.content-grid,
.panel-heading,
.region-item-top,
.template-actions,
.match-item,
.info-grid {
    display: flex;
    gap: 1rem;
}

.hero-card,
.panel-heading,
.region-item-top,
.match-item {
    justify-content: space-between;
}

.hero-actions,
.content-grid,
.info-grid {
    flex-wrap: wrap;
}

.hero-actions {
    justify-content: flex-end;
    align-content: flex-start;
}

.content-grid > *,
.info-grid > * {
    flex: 1 1 320px;
}

.section-kicker,
.summary,
.region-meta,
.info-card p,
.match-item p {
    margin: 0;
}

.section-kicker {
    color: #5f8fbf;
    font-size: 0.76rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
}

h2,
.h3 {
    margin: 0;
}

.summary,
.region-meta,
.info-card p {
    color: #55728f;
}

.status-alert {
    border-radius: 1rem;
}

.primary-button,
.secondary-button,
.danger-button,
.field-input {
    font: inherit;
}

.primary-button,
.secondary-button,
.danger-button {
    border: 0;
    border-radius: 999px;
    padding: 0.8rem 1.1rem;
    cursor: pointer;
}

.primary-button {
    background: linear-gradient(135deg, #5ba9ff 0%, #8fd0ff 100%);
    color: #0d2a44;
    font-weight: 700;
}

.secondary-button {
    background: rgba(91, 169, 255, 0.14);
    color: #24527c;
}

.secondary-button:disabled {
    cursor: not-allowed;
    opacity: 0.56;
}

.danger-button {
    background: rgba(255, 104, 104, 0.16);
    color: #983737;
}

.empty-state {
    display: grid;
    place-items: center;
    min-height: 14rem;
    text-align: center;
    color: #55728f;
}

.preview-shell {
    overflow: auto;
}

.image-stage {
    position: relative;
    display: inline-block;
    max-width: 100%;
    cursor: crosshair;
    border-radius: 1rem;
    overflow: hidden;
    background: rgba(91, 169, 255, 0.08);
}

.preview-image {
    display: block;
    max-width: min(100%, 820px);
    height: auto;
}

.region-box {
    position: absolute;
    border: 2px solid #2c7be5;
    background: rgba(91, 169, 255, 0.14);
    border-radius: 0.45rem;
    box-sizing: border-box;
}

.region-box.draft {
    border-style: dashed;
}

.region-label {
    position: absolute;
    top: 0;
    left: 0;
    transform: translateY(-100%);
    padding: 0.2rem 0.45rem;
    border-radius: 999px;
    background: #2c7be5;
    color: #fff;
    font-size: 0.78rem;
    white-space: nowrap;
}

.region-list,
.template-list,
.match-list {
    display: grid;
    gap: 0.9rem;
}

.region-item,
.template-item,
.match-item,
.info-card {
    border-radius: 1rem;
    background: rgba(91, 169, 255, 0.08);
    padding: 1rem;
}

.template-item.active {
    outline: 2px solid rgba(91, 169, 255, 0.38);
}

.field-label {
    display: grid;
    gap: 0.4rem;
    color: #24527c;
}

.field-input {
    border: 1px solid rgba(91, 169, 255, 0.24);
    border-radius: 0.8rem;
    padding: 0.7rem 0.85rem;
    background: rgba(255, 255, 255, 0.9);
}

.match-state {
    border-radius: 999px;
    padding: 0.35rem 0.75rem;
    font-weight: 700;
}

.match-state.match {
    background: rgba(82, 196, 26, 0.14);
    color: #2f8f19;
}

.match-state.mismatch {
    background: rgba(255, 120, 117, 0.14);
    color: #c0392b;
}

@media (max-width: 900px) {
    .hero-card,
    .panel-heading,
    .region-item-top,
    .match-item {
        flex-direction: column;
    }

    .hero-actions {
        justify-content: flex-start;
    }
}
</style>
