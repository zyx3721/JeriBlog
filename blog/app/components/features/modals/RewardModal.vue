<!--
项目名称：JeriBlog
文件名称：RewardModal.vue
创建时间：2026-04-16 11:55:19

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：打赏弹窗组件，展示微信和支付宝收款码
-->

<script setup lang="ts">
interface Props {
  modelValue: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const { blogConfig } = useSysConfig()

// 关闭弹窗
const closeModal = () => {
  emit('update:modelValue', false)
}

// 点击遮罩层关闭
const handleMaskClick = (e: MouseEvent) => {
  if (e.target === e.currentTarget) {
    closeModal()
  }
}

// 检查是否配置了打赏功能
const isRewardEnabled = computed(() => {
  return blogConfig.value.wechat_qrcode || blogConfig.value.alipay_qrcode
})
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="reward-modal-mask" @click="handleMaskClick">
        <div class="reward-modal-container">
          <div class="reward-modal-header">
            <h3>请作者喝杯咖啡 ☕</h3>
            <button class="close-btn" @click="closeModal" aria-label="关闭">
              <i class="ri-close-line"></i>
            </button>
          </div>

          <div v-if="!isRewardEnabled" class="reward-modal-body empty">
            <div class="empty-icon">
              <i class="ri-image-line"></i>
            </div>
            <p class="empty-text">暂未开放打赏功能</p>
          </div>

          <div v-else class="reward-modal-body">
            <div class="qrcode-container">
              <div v-if="blogConfig.wechat_qrcode" class="qrcode-item">
                <img :src="blogConfig.wechat_qrcode" alt="微信收款码" class="qrcode-image" />
                <p class="qrcode-label">微信</p>
              </div>
              <div v-if="blogConfig.alipay_qrcode" class="qrcode-item">
                <img :src="blogConfig.alipay_qrcode" alt="支付宝收款码" class="qrcode-image" />
                <p class="qrcode-label">支付宝</p>
              </div>
            </div>
            <p class="reward-tips">感谢您的支持与鼓励！</p>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
.reward-modal-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 1rem;
}

.reward-modal-container {
  background-color: var(--jeri-card-bg);
  border-radius: 1rem;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  max-width: 500px;
  width: 100%;
  overflow: hidden;
}

.reward-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem;
  border-bottom: 1px solid var(--jeri-border-color);

  h3 {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--font-color);
    margin: 0;
  }

  .close-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    color: var(--theme-meta-color);
    font-size: 1.5rem;
    line-height: 1;
    padding: 0.25rem;
    transition: color 0.2s ease;

    &:hover {
      color: var(--font-color);
    }
  }
}

.reward-modal-body {
  padding: 2rem 1.5rem;

  &.empty {
    text-align: center;
    padding: 3rem 1.5rem;

    .empty-icon {
      font-size: 4rem;
      color: var(--theme-meta-color);
      margin-bottom: 1rem;
      opacity: 0.5;
    }

    .empty-text {
      font-size: 1rem;
      color: var(--theme-meta-color);
      margin: 0;
    }
  }

  .qrcode-container {
    display: flex;
    justify-content: center;
    gap: 2rem;
    margin-bottom: 1.5rem;

    .qrcode-item {
      text-align: center;

      .qrcode-image {
        width: 160px;
        height: 160px;
        border-radius: 0.5rem;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        object-fit: cover;
      }

      .qrcode-label {
        margin-top: 0.75rem;
        font-size: 0.95rem;
        font-weight: 500;
        color: var(--font-color);
      }
    }
  }

  .reward-tips {
    text-align: center;
    font-size: 0.9rem;
    color: var(--theme-meta-color);
    margin: 0;
  }
}

// 弹窗动画
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;

  .reward-modal-container {
    transition: transform 0.3s ease;
  }
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;

  .reward-modal-container {
    transform: scale(0.9) translateY(-20px);
  }
}

.modal-enter-to,
.modal-leave-from {
  opacity: 1;

  .reward-modal-container {
    transform: scale(1) translateY(0);
  }
}

@media screen and (max-width: 768px) {
  .reward-modal-body {
    .qrcode-container {
      flex-direction: column;
      align-items: center;
      gap: 1.5rem;
    }
  }
}
</style>
