<!--
项目名称：JeriBlog
文件名称：MobileToc.vue
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：移动端目录组件
-->

<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div v-if="visible" class="mobile-toc-overlay" @click="handleClose">
        <div class="mobile-toc-drawer" @click.stop>
          <div class="mobile-toc-header">
            <h3>目录</h3>
            <button class="close-btn" @click="handleClose">
              <i class="ri-close-line"></i>
            </button>
          </div>
          <div class="mobile-toc-content">
            <div v-if="tocItems.length === 0" class="empty-toc">
              <i class="ri-file-list-line"></i>
              <p>暂无目录</p>
            </div>
            <ul v-else class="toc-list">
              <li
                v-for="item in tocItems"
                :key="item.id"
                :class="['toc-item', `level-${item.level}`, { active: item.id === activeId }]"
                @click="scrollToHeading(item.id)"
              >
                <span class="toc-text">{{ item.text }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
interface TocItem {
  id: string
  text: string
  level: number
}

interface Props {
  visible: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
}>()

const tocItems = ref<TocItem[]>([])
const activeId = ref<string>('')

const handleClose = () => {
  emit('close')
}

const scrollToHeading = (id: string) => {
  const element = document.getElementById(id)
  if (element) {
    const offset = 80 // 顶部偏移量
    const elementPosition = element.getBoundingClientRect().top
    const offsetPosition = elementPosition + window.pageYOffset - offset

    window.scrollTo({
      top: offsetPosition,
      behavior: 'smooth'
    })

    // 关闭抽屉
    handleClose()
  }
}

const extractTocItems = () => {
  const article = document.querySelector('.markdown-content')
  if (!article) return

  const headings = article.querySelectorAll('h1, h2, h3, h4, h5, h6')
  const items: TocItem[] = []

  headings.forEach((heading) => {
    const id = heading.id
    const text = heading.textContent || ''
    const level = parseInt(heading.tagName.substring(1))

    if (id && text) {
      items.push({ id, text, level })
    }
  })

  tocItems.value = items
}

const updateActiveId = () => {
  const article = document.querySelector('.markdown-content')
  if (!article) return

  const headings = article.querySelectorAll('h1, h2, h3, h4, h5, h6')
  const scrollTop = window.pageYOffset || document.documentElement.scrollTop
  const offset = 100

  let currentId = ''

  headings.forEach((heading) => {
    const rect = heading.getBoundingClientRect()
    const top = rect.top + scrollTop

    if (top <= scrollTop + offset) {
      currentId = heading.id
    }
  })

  activeId.value = currentId
}

const handleScroll = () => {
  updateActiveId()
}

onMounted(() => {
  extractTocItems()
  updateActiveId()
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

// 监听 visible 变化，重新提取目录
watch(() => props.visible, (newVal) => {
  if (newVal) {
    extractTocItems()
    updateActiveId()
  }
})
</script>

<style scoped lang="scss">
.mobile-toc-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 2000;
  display: flex;
  justify-content: flex-end;
}

.mobile-toc-drawer {
  width: 80%;
  max-width: 320px;
  height: 100%;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  box-shadow: -2px 0 16px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

[data-theme='dark'] .mobile-toc-drawer {
  background: rgba(30, 30, 30, 0.98);
  box-shadow: -2px 0 16px rgba(0, 0, 0, 0.5);
}

.mobile-toc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  border-bottom: 1px solid var(--flec-border);

  h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: var(--flec-text);
  }

  .close-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: var(--flec-text);
    font-size: 20px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    transition: background 0.2s;

    &:hover {
      background: var(--flec-hover);
    }

    &:active {
      background: var(--flec-active);
    }
  }
}

.mobile-toc-content {
  flex: 1;
  overflow-y: auto;
  padding: 10px 0;
}

.empty-toc {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--flec-text-secondary);

  i {
    font-size: 48px;
    margin-bottom: 12px;
    opacity: 0.5;
  }

  p {
    margin: 0;
    font-size: 14px;
  }
}

.toc-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.toc-item {
  padding: 10px 20px;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--flec-text-secondary);
  font-size: 14px;
  line-height: 1.6;
  border-left: 3px solid transparent;

  &.level-1 {
    padding-left: 20px;
    font-weight: 600;
  }

  &.level-2 {
    padding-left: 30px;
  }

  &.level-3 {
    padding-left: 40px;
  }

  &.level-4 {
    padding-left: 50px;
  }

  &.level-5 {
    padding-left: 60px;
  }

  &.level-6 {
    padding-left: 70px;
  }

  &:hover {
    background: var(--flec-hover);
    color: var(--flec-text);
  }

  &:active {
    background: var(--flec-active);
  }

  &.active {
    color: var(--flec-primary);
    border-left-color: var(--flec-primary);
    background: var(--flec-primary-bg);

    .toc-text {
      font-weight: 500;
    }
  }

  .toc-text {
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

// 抽屉动画
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.3s ease;

  .mobile-toc-drawer {
    transition: transform 0.3s ease;
  }
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;

  .mobile-toc-drawer {
    transform: translateX(100%);
  }
}

// 滚动条样式
.mobile-toc-content {
  &::-webkit-scrollbar {
    width: 4px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background: var(--flec-border);
    border-radius: 2px;

    &:hover {
      background: var(--flec-text-secondary);
    }
  }
}
</style>
