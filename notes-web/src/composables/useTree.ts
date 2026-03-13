/**
 * 树形数据处理 composable
 */

import { ref, computed } from 'vue'

export interface TreeNode {
  id: string
  label: string
  isCategory?: boolean
  isTopic?: boolean
  children?: TreeNode[]
  [key: string]: unknown
}

export type TreeData = Record<string, Record<string, string[]>>

export function useTree() {
  const rawData = ref<TreeData>({})
  const selectedNode = ref<TreeNode | null>(null)

  /**
   * 将原始数据转换为 Element Plus Tree 组件需要的格式
   */
  const treeData: ComputedRef<TreeNode[]> = computed(() => {
    const result: TreeNode[] = []
    
    for (const [category, topics] of Object.entries(rawData.value)) {
      const categoryNode: TreeNode = {
        id: category,
        label: category,
        isCategory: true,
        children: [],
      }
      
      for (const [topic, items] of Object.entries(topics)) {
        const filteredItems = items.filter(item => item !== 'README')
        if (filteredItems.length === 0) continue
        
        const topicNode: TreeNode = {
          id: `${category}/${topic}`,
          label: topic,
          isTopic: true,
          children: filteredItems.map(item => ({
            id: `${category}/${topic}/${item}`,
            label: item,
            category,
            topic,
          })),
        }
        categoryNode.children!.push(topicNode)
      }
      
      if (categoryNode.children!.length > 0) {
        result.push(categoryNode)
      }
    }
    
    return result
  })

  /**
   * 设置原始数据
   */
  function setRawData(data: TreeData) {
    rawData.value = data
  }

  /**
   * 清空选择
   */
  function clearSelection() {
    selectedNode.value = null
  }

  return {
    rawData,
    treeData,
    selectedNode,
    setRawData,
    clearSelection,
  }
}
