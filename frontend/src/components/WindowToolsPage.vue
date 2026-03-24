<template>
    <div class="page-grid">
        <section class="panel panel-tree">
            <div class="panel-heading">
                <div>
                    <p class="section-kicker">Inspect</p>
                    <h2>窗口树</h2>
                </div>
                <el-button size="small" @click="refreshTree">刷新</el-button>
            </div>

            <el-tree
                ref="treeRef"
                :data="treeData"
                :props="treeProps"
                :load="loadNode"
                lazy
                node-key="hwnd"
                highlight-current
                :current-node-key="selectedHwnd"
                :expand-on-click-node="false"
                @current-change="handleCurrentChange"
            >
                <template #default="{ data }">
                    <span class="tree-node" :title="`句柄: ${data.hwnd}`">
                        <span class="tree-node-icon">•</span>
                        <span>{{ data.label }}</span>
                    </span>
                </template>
            </el-tree>
        </section>

        <section class="page-stack">
            <article class="panel">
                <div class="panel-heading">
                    <div>
                        <p class="section-kicker">Selected</p>
                        <h2>窗口详情</h2>
                    </div>
                </div>

                <div v-if="selectedHwnd !== null && selectedWindowInfo" class="details">
                    <div class="info-row">
                        <span class="label">标题</span>
                        <span class="value">{{ selectedWindowInfo.title || '(无标题)' }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">类名</span>
                        <span class="value">{{ selectedWindowInfo.className || '(未知)' }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">句柄</span>
                        <span class="value">{{ selectedHwnd }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">父窗口</span>
                        <span class="value">{{ parentHwnd || '无' }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">子窗口数</span>
                        <span class="value">{{ childWindows.length }}</span>
                    </div>

                    <div class="actions">
                        <button class="action action-emphasis" type="button" @click="highlightSelectedWindow">
                            高亮窗口
                        </button>
                        <button class="action" type="button" @click="loadParentWindow">查看父窗口</button>
                    </div>

                    <div v-if="childWindows.length > 0" class="child-list">
                        <h3>子窗口</h3>
                        <ul>
                            <li v-for="child in childWindows" :key="child">{{ child }}</li>
                        </ul>
                    </div>
                </div>
                <div v-else-if="selectedHwnd !== null" class="placeholder">加载窗口详情中…</div>
                <div v-else class="placeholder">从当前视图中选择一个窗口以查看详情。</div>
            </article>

            <article class="panel dual-panel">
                <div>
                    <div class="panel-heading compact">
                        <div>
                            <p class="section-kicker">Cursor</p>
                            <h2>鼠标下窗口</h2>
                        </div>
                    </div>

                    <div class="mouse-card">
                        <template v-if="mouseWindowInfo">
                            <p>坐标: X={{ mouseWindowInfo.x }}, Y={{ mouseWindowInfo.y }}</p>
                            <p>标题: {{ mouseWindowInfo.title || '(无标题)' }}</p>
                            <p>类名: {{ mouseWindowInfo.className || '(未知)' }}</p>
                            <p>句柄: {{ mouseWindowInfo.hwnd || 0 }}</p>
                            <button
                                class="action action-emphasis"
                                type="button"
                                :disabled="mouseWindowInfo.hwnd === 0"
                                @click="highlightMouseWindow"
                            >
                                高亮此窗口
                            </button>
                        </template>
                        <p v-else>移动鼠标到任意窗口以查看当前目标。</p>
                    </div>
                </div>

                <div>
                    <div class="panel-heading compact">
                        <div>
                            <p class="section-kicker">Lookup</p>
                            <h2>手动查询</h2>
                        </div>
                    </div>

                    <div class="manual-card">
                        <el-input
                            v-model.number="manualHwnd"
                            type="number"
                            placeholder="输入窗口句柄"
                            clearable
                        />
                        <button class="action" type="button" @click="loadManualWindow">查询</button>
                    </div>
                </div>
            </article>
        </section>
    </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
    GetChildWindows,
    GetMousePosition,
    GetParentWindow,
    GetTopWindows,
    GetWindowInfo,
    GetWindowUnderMouse,
    HighlightWindow,
} from '../../wailsjs/go/main/App'

const treeData = ref([])
const selectedHwnd = ref(null)
const selectedWindowInfo = ref(null)
const parentHwnd = ref(null)
const childWindows = ref([])
const mouseWindowInfo = ref(null)
const manualHwnd = ref(null)
const treeRef = ref(null)
let mouseUpdateTimeout = null

const treeProps = {
    label: 'label',
    children: 'children',
    isLeaf: 'leaf',
}

function formatWindowLabel(title, className) {
    return `${title && title.trim() ? title : '(无标题)'} (${className || '未知类名'})`
}

async function loadWindows() {
    try {
        const windows = await GetTopWindows()
        treeData.value = windows.map((win) => ({
            hwnd: win.hwnd,
            title: win.title,
            className: win.className,
            label: formatWindowLabel(win.title, win.className),
            leaf: false,
        }))
    } catch (error) {
        console.error('获取窗口列表失败:', error)
        ElMessage.error(`获取窗口列表失败: ${error}`)
    }
}

async function loadNode(node, resolve) {
    if (!node.data?.hwnd) {
        resolve([])
        return
    }

    try {
        const childHandles = await GetChildWindows(node.data.hwnd)
        if (!childHandles.length) {
            resolve([])
            return
        }

        const childNodes = await Promise.all(
            childHandles.map(async (childHwnd) => {
                try {
                    const info = await GetWindowInfo(childHwnd)
                    return {
                        hwnd: childHwnd,
                        title: info.title,
                        className: info.className,
                        label: formatWindowLabel(info.title, info.className),
                        leaf: false,
                    }
                } catch (error) {
                    console.error(`获取子窗口信息失败: ${childHwnd}`, error)
                    return {
                        hwnd: childHwnd,
                        title: '',
                        className: '未知',
                        label: `(未知窗口) [${childHwnd}]`,
                        leaf: true,
                    }
                }
            }),
        )

        resolve(childNodes)
    } catch (error) {
        console.error(`加载子窗口失败: ${node.data.hwnd}`, error)
        ElMessage.warning(`加载子窗口失败: ${error}`)
        resolve([])
    }
}

async function refreshTree() {
    selectedHwnd.value = null
    selectedWindowInfo.value = null
    parentHwnd.value = null
    childWindows.value = []
    await loadWindows()
}

async function selectWindow(hwnd) {
    selectedHwnd.value = hwnd
    selectedWindowInfo.value = null
    parentHwnd.value = null
    childWindows.value = []

    try {
        const info = await GetWindowInfo(hwnd)
        selectedWindowInfo.value = info
        const parent = await GetParentWindow(hwnd)
        parentHwnd.value = parent !== 0 ? parent : null
        childWindows.value = await GetChildWindows(hwnd)
    } catch (error) {
        console.error('获取窗口详情失败:', error)
        ElMessage.error(`获取窗口详情失败: ${error}`)
    }
}

async function handleCurrentChange(data) {
    if (!data?.hwnd) {
        return
    }
    await selectWindow(data.hwnd)
}

async function highlightSelectedWindow() {
    if (!selectedHwnd.value) {
        return
    }
    try {
        await HighlightWindow(selectedHwnd.value)
    } catch (error) {
        console.error('高亮窗口失败:', error)
        ElMessage.error(`高亮窗口失败: ${error}`)
    }
}

async function loadParentWindow() {
    if (!parentHwnd.value) {
        ElMessage.info('当前窗口没有父窗口')
        return
    }

    treeRef.value?.setCurrentKey(parentHwnd.value)
    await selectWindow(parentHwnd.value)
}

async function updateMouseWindow() {
    try {
        const hwnd = await GetWindowUnderMouse()
        const position = await GetMousePosition()

        if (hwnd === 0) {
            mouseWindowInfo.value = {
                hwnd: 0,
                title: '',
                className: '',
                x: position.x,
                y: position.y,
            }
        } else {
            const info = await GetWindowInfo(hwnd)
            mouseWindowInfo.value = {
                hwnd,
                title: info.title,
                className: info.className,
                x: position.x,
                y: position.y,
            }
        }
    } catch (error) {
        console.error('获取鼠标下窗口失败:', error)
    } finally {
        mouseUpdateTimeout = setTimeout(updateMouseWindow, 200)
    }
}

async function highlightMouseWindow() {
    if (!mouseWindowInfo.value?.hwnd) {
        return
    }
    try {
        await HighlightWindow(mouseWindowInfo.value.hwnd)
    } catch (error) {
        console.error('高亮鼠标下窗口失败:', error)
        ElMessage.error(`高亮失败: ${error}`)
    }
}

async function loadManualWindow() {
    if (!manualHwnd.value) {
        ElMessage.info('请输入有效的窗口句柄')
        return
    }

    try {
        await GetWindowInfo(manualHwnd.value)
        treeRef.value?.setCurrentKey(manualHwnd.value)
        await selectWindow(manualHwnd.value)
    } catch (error) {
        console.error('手动查询窗口失败:', error)
        ElMessage.error('无效的窗口句柄或无法获取窗口信息')
    }
}

onMounted(() => {
    loadWindows()
    updateMouseWindow()
})

onUnmounted(() => {
    if (mouseUpdateTimeout) {
        clearTimeout(mouseUpdateTimeout)
    }
})
</script>

<style scoped>
.page-grid {
    height: 100%;
    display: grid;
    grid-template-columns: minmax(280px, 360px) minmax(0, 1fr);
    gap: 1.25rem;
}

.page-stack {
    display: grid;
    gap: 1.25rem;
    min-width: 0;
}

.panel {
    min-width: 0;
    border-radius: 1.5rem;
    background: rgba(255, 255, 255, 0.82);
    box-shadow:
        0 18px 52px rgba(101, 157, 212, 0.14),
        inset 0 1px 0 rgba(255, 255, 255, 0.92);
    padding: 1.25rem;
    backdrop-filter: blur(14px);
}

.panel-tree {
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.panel-heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1rem;
}

.panel-heading.compact {
    margin-bottom: 0.75rem;
}

.section-kicker {
    margin: 0 0 0.2rem;
    color: #5f8fbf;
    font-size: 0.76rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
}

h2,
h3 {
    margin: 0;
}

.tree-node {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    min-width: 0;
}

.tree-node-icon {
    color: #5ba9ff;
    flex-shrink: 0;
    font-size: 1.2rem;
    line-height: 1;
}

.details {
    display: grid;
    gap: 0.7rem;
}

.info-row {
    display: grid;
    grid-template-columns: 6rem 1fr;
    gap: 0.75rem;
    align-items: start;
}

.label {
    color: #6c88a3;
    font-weight: 700;
}

.value {
    word-break: break-all;
}

.actions {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    margin-top: 0.35rem;
}

.action {
    border: 0;
    border-radius: 999px;
    padding: 0.7rem 1rem;
    font: inherit;
    cursor: pointer;
    color: #17324d;
    background: rgba(91, 169, 255, 0.12);
    transition:
        transform 0.18s ease,
        background-color 0.18s ease;
}

.action:hover {
    transform: translateY(-1px);
    background: rgba(91, 169, 255, 0.2);
}

.action:disabled {
    cursor: not-allowed;
    opacity: 0.55;
}

.action-emphasis {
    background: linear-gradient(135deg, #5ba9ff 0%, #8fd0ff 100%);
    color: #0d2a44;
}

.child-list {
    margin-top: 0.75rem;
}

.child-list ul {
    margin: 0.65rem 0 0;
    padding: 0;
    list-style: none;
    max-height: 220px;
    overflow: auto;
    border-radius: 1rem;
    background: rgba(91, 169, 255, 0.08);
}

.child-list li {
    padding: 0.65rem 0.8rem;
    border-bottom: 1px solid rgba(91, 169, 255, 0.12);
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
}

.placeholder,
.mouse-card p {
    margin: 0;
}

.placeholder {
    color: #64819d;
}

.dual-panel {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 1rem;
}

.mouse-card,
.manual-card {
    display: grid;
    gap: 0.75rem;
    min-height: 100%;
    padding: 1rem;
    border-radius: 1rem;
    background: rgba(91, 169, 255, 0.08);
}

@media (max-width: 1100px) {
    .page-grid {
        grid-template-columns: 1fr;
    }
}

@media (max-width: 720px) {
    .dual-panel {
        grid-template-columns: 1fr;
    }
}
</style>
