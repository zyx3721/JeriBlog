/*
项目名称：JeriBlog
文件名称：useExpandable.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

export function useExpandable(initialState = false) {
  const isExpanded = ref(initialState)
  const toggleExpand = () => isExpanded.value = !isExpanded.value
  return { isExpanded, toggleExpand }
}