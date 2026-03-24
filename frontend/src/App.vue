<template>
    <div class="app">
        <header>
            <h1>Window Tools</h1>
        </header>

        <div class="main-container">
            <!-- 左侧：窗口树形框 -->
            <aside class="window-tree">
                <div class="tree-header">
                    <h2>窗口树</h2>
                    <el-button size="small" @click="refreshTree" :icon="Refresh">刷新</el-button>
                </div>
                <el-tree ref="treeRef" :data="treeData" :props="treeProps" :load="loadNode" lazy node-key="hwnd"
                    highlight-current :current-node-key="selectedHwnd" @current-change="handleCurrentChange"
                    :expand-on-click-node="false">
                    <template #default="{ node, data }">
                        <span class="tree-node" :title="`句柄: ${data.hwnd}`">
                            <el-icon>
                                <Monitor />
                            </el-icon>
                            <span>{{ data.label }}</span>
                        </span>
                    </template>
                </el-tree>
            </aside>

            <!-- 右侧：窗口详情与操作 -->
            <main class="window-detail">
                <h2>窗口详情</h2>
                <div v-if="selectedHwnd !== null && selectedWindowInfo" class="detail-card">
                    <div class="info-row">
                        <span class="label">标题：</span>
                        <span class="value">{{ selectedWindowInfo.title || '(无标题)' }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">类名：</span>
                        <span class="value">{{ selectedWindowInfo.className }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">句柄：</span>
                        <span class="value">{{ selectedHwnd }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">父窗口：</span>
                        <span class="value">{{ parentHwnd || '无' }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">子窗口数：</span>
                        <span class="value">{{ childWindows.length }}</span>
                    </div>
                    <div class="actions">
                        <button @click="highlightSelectedWindow" class="btn highlight">高亮窗口</button>
                        <button @click="loadParentWindow" class="btn info">查看父窗口</button>
                    </div>
                    <div v-if="childWindows.length > 0" class="child-list">
                        <h3>子窗口列表</h3>
                        <ul>
                            <li v-for="child in childWindows" :key="child">
                                {{ child }}
                            </li>
                        </ul>
                    </div>
                </div>
                <div v-else-if="selectedHwnd !== null && !selectedWindowInfo" class="loading">
                    加载中...
                </div>
                <div v-else class="no-selection">请从左侧选择窗口</div>

                <!-- 鼠标下窗口区域 -->
                <div class="mouse-window">
                    <h2>鼠标下窗口</h2>
                    <div class="mouse-info">
                        <div v-if="mouseWindowInfo">
                            <div>坐标：X={{ mouseWindowInfo.x }}, Y={{ mouseWindowInfo.y }}</div>
                            <div>标题：{{ mouseWindowInfo.title || '(无标题)' }}</div>
                            <div>类名：{{ mouseWindowInfo.className }}</div>
                            <div>句柄：{{ mouseWindowInfo.hwnd }}</div>
                            <button @click="highlightMouseWindow" class="btn highlight">高亮此窗口</button>
                        </div>
                        <div v-else>移动鼠标到任意窗口...</div>
                    </div>
                </div>

                <!-- 手动输入句柄 -->
                <div class="manual-input">
                    <h2>手动查询窗口</h2>
                    <input v-model.number="manualHwnd" type="number" placeholder="输入窗口句柄" />
                    <button @click="loadManualWindow" class="btn info">查询</button>
                </div>
            </main>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Monitor, Refresh } from '@element-plus/icons-vue'
import {
    GetTopWindows,
    GetWindowUnderMouse,
    GetWindowInfo,
    HighlightWindow,
    GetParentWindow,
    GetChildWindows,
    GetMousePosition
} from '../wailsjs/go/main/App'

// 响应式数据
const treeData = ref([])
const selectedHwnd = ref(null)
const selectedWindowInfo = ref(null)
const parentHwnd = ref(null)
const childWindows = ref([])
const mouseWindowInfo = ref(null)
const manualHwnd = ref(null)
const treeRef = ref(null)

// 定时器，用于实时更新鼠标下窗口
let mouseUpdateTimeout = null

// 树形框配置
const treeProps = {
    label: 'label',
    children: 'children',
    isLeaf: 'leaf'
}

// 加载所有顶级窗口
async function loadWindows() {
    try {
        const windows = await GetTopWindows()
        // 构建树节点数据
        treeData.value = windows.map(win => ({
            hwnd: win.hwnd,
            title: win.title,
            className: win.className,
            label: formatWindowLabel(win.title, win.className, win.hwnd),
            leaf: false // 非叶子节点，允许懒加载
        }))
    } catch (err) {
        console.error('获取窗口列表失败:', err)
        ElMessage.error('获取窗口列表失败: ' + err)
    }
}

// 格式化窗口显示文本
function formatWindowLabel(title, className, hwnd) {
    const displayTitle = title && title.trim() ? title : '(无标题)'
    return `${displayTitle} (${className})`
}

// 懒加载子节点
async function loadNode(node, resolve) {
    // 根节点不处理（根节点数据已由 treeData 提供）
    if (!node.data || !node.data.hwnd) {
        return resolve([])
    }

    const hwnd = node.data.hwnd

    try {
        // 获取子窗口句柄列表
        const childHandles = await GetChildWindows(hwnd)

        if (!childHandles || childHandles.length === 0) {
            return resolve([])
        }

        // 并发获取所有子窗口的详细信息
        const childNodes = await Promise.all(
            childHandles.map(async (childHwnd) => {
                try {
                    const { title, className } = await GetWindowInfo(childHwnd)
                    return {
                        hwnd: childHwnd,
                        title: title,
                        className: className,
                        label: formatWindowLabel(title, className, childHwnd),
                        leaf: false // 允许继续展开，实际是否有子节点在展开时再判断
                    }
                } catch (err) {
                    console.error(`获取子窗口信息失败 (句柄: ${childHwnd}):`, err)
                    // 返回一个占位节点，至少显示句柄
                    return {
                        hwnd: childHwnd,
                        title: '',
                        className: '未知',
                        label: `(未知窗口) [${childHwnd}]`,
                        leaf: true // 标记为叶子，避免再次展开
                    }
                }
            })
        )

        resolve(childNodes)
    } catch (err) {
        console.error(`加载子窗口失败 (句柄: ${hwnd}):`, err)
        ElMessage.warning(`加载子窗口失败: ${err}`)
        resolve([])
    }
}

// 刷新树形框
async function refreshTree() {
    await loadWindows()
    // 清空选中状态
    selectedHwnd.value = null
    selectedWindowInfo.value = null
    parentHwnd.value = null
    childWindows.value = []
}

// 选中节点
async function handleCurrentChange(data) {
    if (!data || !data.hwnd) return
    await selectWindow(data.hwnd)
}

// 选择窗口并加载详细信息
async function selectWindow(hwnd) {
    selectedHwnd.value = hwnd
    selectedWindowInfo.value = null
    parentHwnd.value = null
    childWindows.value = []

    try {
        // 获取窗口信息
        const { title, className } = await GetWindowInfo(hwnd)
        selectedWindowInfo.value = { title, className }
        // 获取父窗口
        const parent = await GetParentWindow(hwnd)
        parentHwnd.value = parent !== 0 ? parent : null
        // 获取子窗口列表（仅显示句柄，用于展示）
        const children = await GetChildWindows(hwnd)
        childWindows.value = children
    } catch (err) {
        console.error('获取窗口详情失败:', err)
        ElMessage.error('获取窗口详情失败: ' + err)
    }
}

// 高亮当前选中的窗口
async function highlightSelectedWindow() {
    if (!selectedHwnd.value) return
    try {
        await HighlightWindow(selectedHwnd.value)
    } catch (err) {
        console.error('高亮窗口失败:', err)
        ElMessage.error('高亮窗口失败: ' + err)
    }
}

// 加载父窗口信息（将选中父窗口）
async function loadParentWindow() {
    if (parentHwnd.value) {
        // 在树中选中父窗口
        if (treeRef.value) {
            treeRef.value.setCurrentKey(parentHwnd.value)
        }
        await selectWindow(parentHwnd.value)
    } else {
        ElMessage.info('当前窗口没有父窗口')
    }
}

// 获取鼠标下窗口信息并更新
async function updateMouseWindow() {
    try {
        // 同时获取鼠标下窗口句柄和鼠标坐标
        const hwnd = await GetWindowUnderMouse()
        const pos = await GetMousePosition()

        let windowInfo = null
        if (hwnd !== 0) {
            const { title, className } = await GetWindowInfo(hwnd)
            windowInfo = { hwnd, title, className, x: pos.x, y: pos.y }
        } else {
            windowInfo = { hwnd: 0, title: '', className: '', x: pos.x, y: pos.y }
        }
        mouseWindowInfo.value = windowInfo
    } catch (err) {
        console.error('获取鼠标下窗口失败:', err)
    } finally {
        // 无论成功或失败，都延迟 100ms 后再次执行
        mouseUpdateTimeout = setTimeout(updateMouseWindow, 100)
    }
}

// 高亮鼠标下的窗口
async function highlightMouseWindow() {
    if (!mouseWindowInfo.value || mouseWindowInfo.value.hwnd === 0) return
    try {
        await HighlightWindow(mouseWindowInfo.value.hwnd)
    } catch (err) {
        console.error('高亮鼠标下窗口失败:', err)
        ElMessage.error('高亮失败: ' + err)
    }
}

// 手动查询窗口
async function loadManualWindow() {
    if (!manualHwnd.value) return
    const hwnd = manualHwnd.value

    try {
        // 尝试加载窗口信息
        const [title, className] = await GetWindowInfo(hwnd)

        // 在树中查找该节点（如果存在）
        if (treeRef.value) {
            treeRef.value.setCurrentKey(hwnd)
        }

        // 直接选中该窗口
        await selectWindow(hwnd)

        // 可选：如果窗口不在当前树中，提示用户
        const existsInTree = findNodeInTree(treeData.value, hwnd)
        if (!existsInTree) {
            ElMessage.info('该窗口不在树形列表中，可能是子窗口或已关闭')
        }
    } catch (err) {
        console.error('获取窗口信息失败:', err)
        ElMessage.error('无效的窗口句柄或无法获取信息')
    }
}

// 在树中递归查找节点（辅助函数）
function findNodeInTree(nodes, hwnd) {
    for (const node of nodes) {
        if (node.hwnd === hwnd) return true
        // 注意：由于懒加载，子节点可能未加载，无法在未展开时查找到
    }
    return false
}

// 生命周期
onMounted(() => {
    loadWindows()
    updateMouseWindow()   // 启动递归
})

onUnmounted(() => {
    if (mouseUpdateTimeout) clearTimeout(mouseUpdateTimeout)  // 清理
})
</script>

<style scoped>
* {
    box-sizing: border-box;
}

.app {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
    height: 100vh;
    display: flex;
    flex-direction: column;
    background-color: #f5f5f5;
    color: #333;
}

header {
    background: #2c3e50;
    color: white;
    padding: 1rem;
    text-align: center;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h1 {
    margin: 0;
    font-size: 1.5rem;
}

.main-container {
    display: flex;
    flex: 1;
    overflow: hidden;
}

.window-tree {
    width: 360px;
    background: white;
    border-right: 1px solid #ddd;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.tree-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1rem 0.5rem 1rem;
    border-bottom: 1px solid #eee;
}

.tree-header h2 {
    margin: 0;
    font-size: 1.2rem;
}

.window-tree :deep(.el-tree) {
    flex: 1;
    overflow-y: auto;
    padding: 0.5rem;
    background: transparent;
}

.tree-node {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.9rem;
    cursor: pointer;
    width: 100%;
}

.tree-node .el-icon {
    font-size: 1rem;
    color: #409eff;
    flex-shrink: 0;
}

.window-detail {
    flex: 1;
    padding: 1rem;
    overflow-y: auto;
    background: #fafafa;
}

.detail-card {
    background: white;
    border-radius: 8px;
    padding: 1rem;
    margin-bottom: 1rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.info-row {
    margin-bottom: 0.5rem;
    display: flex;
}

.label {
    width: 80px;
    font-weight: 600;
    color: #555;
}

.value {
    flex: 1;
    word-break: break-all;
}

.actions {
    margin-top: 1rem;
    display: flex;
    gap: 0.5rem;
}

.btn {
    padding: 0.4rem 0.8rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background 0.2s;
}

.highlight {
    background: #ff5252;
    color: white;
}

.highlight:hover {
    background: #ff1744;
}

.info {
    background: #2196f3;
    color: white;
}

.info:hover {
    background: #1976d2;
}

.child-list {
    margin-top: 1rem;
}

.child-list h3 {
    margin: 0 0 0.5rem;
    font-size: 1rem;
}

.child-list ul {
    max-height: 200px;
    overflow-y: auto;
    background: #f5f5f5;
    border-radius: 4px;
    padding: 0.5rem;
    list-style: none;
    margin: 0;
}

.child-list li {
    font-family: monospace;
    font-size: 0.8rem;
    padding: 0.2rem 0;
    border-bottom: 1px solid #e0e0e0;
}

.mouse-window,
.manual-input {
    background: white;
    border-radius: 8px;
    padding: 1rem;
    margin-top: 1rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.mouse-window h2,
.manual-input h2 {
    margin: 0 0 0.5rem;
    font-size: 1.2rem;
}

.mouse-info {
    min-height: 80px;
}

.manual-input {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    flex-wrap: wrap;
}

.manual-input input {
    flex: 1;
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 4px;
}

.loading,
.no-selection {
    text-align: center;
    color: #999;
    padding: 2rem;
}

/* 树形框节点悬停效果 */
.window-tree :deep(.el-tree-node__content) {
    height: auto;
    padding: 0.3rem 0;
}

.window-tree :deep(.el-tree-node__content:hover) {
    background-color: #f0f0f0;
}

.window-tree :deep(.el-tree-node.is-current > .el-tree-node__content) {
    background-color: #e3f2fd;
}
</style>