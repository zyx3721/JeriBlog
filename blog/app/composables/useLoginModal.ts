/*
项目名称：JeriBlog
文件名称：useLoginModal.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

// 全局共享的登录弹窗状态
const showLoginModal = ref(false)

export const useLoginModal = () => {
  const open = () => {
    showLoginModal.value = true
  }
  
  const close = () => {
    showLoginModal.value = false
  }
  
  return {
    showLoginModal,
    open,
    close
  }
}
